package dps

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xackery/critsprinkler/tracker"
)

type DPS struct {
	parseStart    time.Time
	zone          string
	damageEvents  map[string][]*DamageEvent
	onDamageEvent []func(*DamageEvent)
}

type Orientation int

const (
	OrientationNone Orientation = iota
	OrientationTopLeft
	OrientationTop
	OrientationTopRight
	OrientationRight
	OrientationBottomRight
	OrientationBottom
	OrientationBottomLeft
	OrientationLeft
)

type DamageEvent struct {
	Orientation Orientation
	SpellName   string
	Source      string
	Target      string
	Type        string
	Damage      int
	Event       time.Time
	Origin      string
	IsCritical  bool
}

func (d *DamageEvent) String() string {
	return fmt.Sprintf("%+v", *d)
}

var (
	instance              *DPS
	myLastSpellCrit       = 0
	myLastSpellName       = ""
	myLastHealCrit        = 0
	myLastMeleeCrit       = 0
	lastOtherHealCritName = ""
)

func New() (*DPS, error) {
	if instance != nil {
		return nil, fmt.Errorf("dps already exists")
	}
	a := &DPS{
		zone:         "Unknown",
		parseStart:   time.Now(),
		damageEvents: make(map[string][]*DamageEvent),
	}

	err := tracker.Subscribe(a.onLine)
	if err != nil {
		return nil, fmt.Errorf("tracker subscribe: %w", err)
	}

	err = tracker.SubscribeToZoneEvent(a.onZone)
	if err != nil {
		return nil, fmt.Errorf("tracker subscribe to zone: %w", err)
	}

	instance = a
	return a, nil
}

func (a *DPS) onLine(event time.Time, line string) {
	a.onMySpell(line)
	a.onMySpellCrit(line)
	a.onMyMelee(event, line)
	a.onMyMeleeCrit(line)
	a.onMyMeleeSlay(line)
	a.onMyMeleeCleaving(line)
	a.onSpell(event, line)
	//a.onDotDamageDPS(event, line)
	a.onMyHealCrit(line)
	a.onHealCrit(line)
	a.onHeal(event, line)
	a.onRune(event, line)
	a.dumpDPS(event)
}

func (a *DPS) onZone(event time.Time, zoneName string) {
	a.zone = zoneName

	a.dumpDPS(event)
}

func (a *DPS) dumpDPS(event time.Time) {
	//dpsPerSec := float64(a.totalDPSGained) / time.Since(a.parseStart).Seconds()
	//dpsPerHour := dpsPerSec * 3600

	if a.zone == "The Bazaar" {
		return
	}

	if len(a.damageEvents) == 0 {
		//fmt.Println("No damage events to report")
		return
	}

	//fmt.Println(len(a.damageEvents), "events to report")
	type dpsReport struct {
		total    int
		maxMelee int
		maxSpell int
	}

	damageTotals := make(map[string]dpsReport)

	tmpDamageEvents := make(map[string][]*DamageEvent)

	for name, dmgEvents := range a.damageEvents {

		for _, dmgEvent := range dmgEvents {

			// skip any events older than 60s
			if event.Sub(dmgEvent.Event).Seconds() > 60 {
				continue
			}

			if len(tmpDamageEvents[name]) == 0 {
				tmpDamageEvents[name] = make([]*DamageEvent, 0)
			}

			tmpDamageEvents[name] = append(tmpDamageEvents[name], dmgEvent)
			dps, ok := damageTotals[name]
			if !ok {
				dps = dpsReport{}
			}

			dps.total += dmgEvent.Damage
			if dmgEvent.Origin == "melee" {
				if dps.maxMelee < dmgEvent.Damage {
					dps.maxMelee = dmgEvent.Damage
				}
			} else {
				if dps.maxSpell < dmgEvent.Damage {
					dps.maxSpell = dmgEvent.Damage
				}
			}
			damageTotals[name] = dps
		}
	}

	a.damageEvents = tmpDamageEvents

	//fmt.Println(len(a.damageEvents), "events to report after filtering")

}

func (a *DPS) onMySpell(line string) {
	match := regexp.MustCompile(`\] You begin casting (.*)\.`).FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	myLastSpellName = match[1]
}

func (a *DPS) onMyHealCrit(line string) {
	match := regexp.MustCompile(`\] You perform an exceptional heal! \((.*)\)`).FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	myLastHealCrit = amount
}

func (a *DPS) onMySpellCrit(line string) {
	match := regexp.MustCompile(`\] You deliver a critical blast! \((.*)\) \((.*)\)`).FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	myLastSpellCrit = amount
}

func (a *DPS) onMyMeleeCrit(line string) {
	match := regexp.MustCompile(`\] (.*) scores a critical hit! \((.*)\)`).FindStringSubmatch(line)
	if len(match) < 3 {
		return
	}

	amount, err := strconv.Atoi(match[2])
	if err != nil {
		fmt.Println("atoi", line, match[2], err)
		return
	}

	if match[1] != tracker.PlayerName() {
		return
	}

	myLastMeleeCrit = amount
}

func (a *DPS) onMyMeleeSlay(line string) {
	match := regexp.MustCompile(`\] (.*) holy blade cleanses (.) target!\((.*)\)`).FindStringSubmatch(line)
	if len(match) < 3 {
		return
	}

	amount, err := strconv.Atoi(match[2])
	if err != nil {
		return
	}

	if match[0] != tracker.PlayerName() {
		return
	}

	myLastMeleeCrit = amount
}

func (a *DPS) onMyMeleeCleaving(line string) {
	match := regexp.MustCompile(`\] (.*) lands a Cleaving Blow! \((.*)\)`).FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	if match[0] != tracker.PlayerName() {
		return
	}

	myLastMeleeCrit = amount
}

func (a *DPS) onHealCrit(line string) {

	match := regexp.MustCompile(`\] (.*) performs an exceptional heal! \((.*)\)`).FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	lastOtherHealCritName = match[1]

}

func (a *DPS) onHeal(event time.Time, line string) {
	match := regexp.MustCompile(`\] (.*) has healed (.*) for (.*) points of damage. \((.*)\)`).FindStringSubmatch(line)
	if len(match) < 5 {
		return
	}

	amount, err := strconv.Atoi(match[3])
	if err != nil {
		fmt.Println("atoi", line, match[3], err)
		return
	}

	source := match[1]
	target := match[2]
	if source == "you" {
		target = tracker.PlayerName()
	}

	if strings.EqualFold(target, "itself") || strings.EqualFold(target, "himself") || strings.EqualFold(target, "herself") {
		target = source
	}

	damageEvent := &DamageEvent{
		Source:    source,
		Target:    target,
		SpellName: match[4],
		Type:      "heal",
		Damage:    amount,
		Event:     event,
		Origin:    "heal",
	}

	if lastOtherHealCritName == match[1] {
		damageEvent.IsCritical = true
		lastOtherHealCritName = ""
	}

	if source == tracker.PlayerName() && myLastHealCrit == amount {
		damageEvent.IsCritical = true
		myLastHealCrit = 0
	}

	for _, fn := range a.onDamageEvent {
		fn(damageEvent)
	}

}

func (a *DPS) onRune(event time.Time, line string) {
	match := regexp.MustCompile(`\] (.*) has shielded (.*) from (.*) points of damage. \((.*)\)`).FindStringSubmatch(line)
	if len(match) < 4 {
		return
	}

	amount, err := strconv.Atoi(match[3])
	if err != nil {
		return
	}

	source := match[1]
	target := match[2]
	if source == "you" {
		target = tracker.PlayerName()
	}

	if strings.EqualFold(target, "itself") || strings.EqualFold(target, "himself") || strings.EqualFold(target, "herself") {
		target = source
	}

	damageEvent := &DamageEvent{
		Source: source,
		Target: target,
		Type:   "heal",
		Damage: amount,
		Event:  event,
		Origin: "heal",
	}

	if lastOtherHealCritName == match[1] {
		damageEvent.IsCritical = true
		lastOtherHealCritName = ""
	}

	for _, fn := range a.onDamageEvent {
		fn(damageEvent)
	}

}

func (a *DPS) onMyMelee(event time.Time, line string) {
	match := regexp.MustCompile(`\] (.*) for (.*) points of damage.`).FindStringSubmatch(line)
	if len(match) < 3 {
		return
	}

	amount, err := strconv.Atoi(match[2])
	if err != nil {
		return
	}

	chunk := match[1]

	if strings.Contains(chunk, " was hit ") {
		return
	}

	pos := 0
	pickedAdj := ""
	for _, adj := range adjectives {
		pos = strings.Index(chunk, adj)
		if pos <= 0 {
			continue
		}
		pickedAdj = adj
		break
	}
	if pos <= 0 {
		return
	}

	source := chunk[:pos]
	if strings.EqualFold(source, "you") {
		source = tracker.PlayerName()
	}
	target := chunk[pos+len(pickedAdj):]
	if strings.EqualFold(target, "you") {
		target = tracker.PlayerName()
	}
	if strings.Contains(source, "`s doppleganger") {
		source = strings.ReplaceAll(source, "`s doppleganger", "")
	}
	damageEvent := &DamageEvent{
		Source: source,
		Target: target,
		Type:   strings.TrimSpace(pickedAdj),
		Damage: amount,
		Event:  event,
		Origin: "melee",
	}

	if amount >= myLastMeleeCrit && myLastMeleeCrit > 0 {
		damageEvent.IsCritical = true
		myLastMeleeCrit = 0
	}

	for _, fn := range a.onDamageEvent {
		fn(damageEvent)
	}
	_, ok := a.damageEvents[damageEvent.Source]
	if !ok {
		a.damageEvents[damageEvent.Source] = make([]*DamageEvent, 0)
	}

	a.damageEvents[damageEvent.Source] = append(a.damageEvents[damageEvent.Source], damageEvent)

}

func (a *DPS) onSpell(event time.Time, line string) {
	match := regexp.MustCompile(`\] (.*) hit (.*) for (.*) points of non-melee damage. \((.*)\)`).FindStringSubmatch(line)
	if len(match) < 3 {
		return
	}

	amount, err := strconv.Atoi(match[3])
	if err != nil {
		return
	}

	source := match[1]

	if strings.Contains(source, "`s doppleganger") {
		fmt.Println("FO?UND DOPPLE:", source)
		source = strings.ReplaceAll(source, "`s doppleganger", "")
	}

	damageEvent := &DamageEvent{
		Source:    source,
		Type:      "hit",
		Target:    match[2],
		Damage:    amount,
		SpellName: match[4],
		Event:     event,
		Origin:    "direct",
	}

	if source == tracker.PlayerName() {
		damageEvent.SpellName = myLastSpellName
	}
	if myLastSpellCrit == amount {
		damageEvent.IsCritical = true
		myLastSpellCrit = 0
	}
	for _, fn := range a.onDamageEvent {
		fn(damageEvent)
	}

	_, ok := a.damageEvents[damageEvent.Source]
	if !ok {
		a.damageEvents[damageEvent.Source] = make([]*DamageEvent, 0)
	}

	a.damageEvents[damageEvent.Source] = append(a.damageEvents[damageEvent.Source], damageEvent)
}

func (a *DPS) onDotDamageDPS(event time.Time, line string) {
	match := regexp.MustCompile(`\] (.*) has taken (.*) damage from your (.*).`).FindStringSubmatch(line)

	if len(match) < 3 {
		return
	}

	amount, err := strconv.Atoi(match[2])
	if err != nil {
		return
	}

	source := tracker.PlayerName()
	target := match[1]
	if strings.Contains(source, "`s doppleganger") {
		source = strings.ReplaceAll(source, "`s doppleganger", "")
	}
	damageEvent := &DamageEvent{
		Source: source,
		Target: target,
		Type:   match[3][0 : len(match[3])-2],
		Damage: amount,
		Event:  event,
		Origin: "dot",
	}
	for _, fn := range a.onDamageEvent {
		fn(damageEvent)
	}
	_, ok := a.damageEvents[damageEvent.Source]
	if !ok {
		a.damageEvents[damageEvent.Source] = make([]*DamageEvent, 0)
	}

	a.damageEvents[damageEvent.Source] = append(a.damageEvents[damageEvent.Source], damageEvent)
}

func SubscribeToDamageEvent(fn func(*DamageEvent)) error {
	if instance == nil {
		return fmt.Errorf("dps not initialized")
	}
	instance.onDamageEvent = append(instance.onDamageEvent, fn)
	return nil
}

var adjectives = []string{
	" mauls ",
	" maul ",
	" bites ",
	" bite ",
	" claws ",
	" claw ",
	" gores ",
	" gore ",
	" stings ",
	" slices ",
	" slice ",
	" sting ",
	" smashes ",
	" smash ",
	" rend ",
	" rends ",
	" slash ",
	" slashes ",
	" punch ",
	" punches ",
	" hit ",
	" hits ",
	" You ",
	" yourself ",
	" YOU ",
	" himself ",
	" herself ",
	" itself ",
	" crush ",
	" crushes ",
	" pierce ",
	" pierces ",
	" kick ",
	" kicks ",
	" strike ",
	" strikes ",
	" backstab ",
	" backstabs ",
	" bash ",
	" bashes ",
}
