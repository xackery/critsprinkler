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

var (
	instance              *DPS
	meleeDamageRegex      = regexp.MustCompile(`\] (.*) for (.*) points of damage.`)
	directDamageRegex     = regexp.MustCompile(`\] (.*) hit (.*) for (.*) points of non-melee damage. \((.*)\)`)
	dotDamageRegex        = regexp.MustCompile(`\] (.*) has taken (.*) damage from your (.*).`)
	directDamageCritRegex = regexp.MustCompile(`\] You deliver a critical blast! \((.*)\) \((.*)\)`)
	myHealCritRegex       = regexp.MustCompile(`You perform an exceptional heal! \((.*)\)`)
	myHealRegex           = regexp.MustCompile(`You have healed (.*) for (.*) points.`)
	healCritRegex         = regexp.MustCompile(`(.*) performs an exceptional heal! \((.*)\)`)
	healRegex             = regexp.MustCompile(`(.*) has healed (.*) for (.*) points.`)
	spellNameRegex        = regexp.MustCompile(`You begin casting (.*)\.`)
	myLastCrit            = 0
	myLastCritHeal        = 0
	myLastSpellName       = ""
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
	a.onMySpellCast(line)
	a.onMySpellCrit(line)
	a.onMeleeDPS(event, line)
	a.onDirectDamageDPS(event, line)
	//a.onDotDamageDPS(event, line)
	a.onMyHealSpellCrit(line)
	a.onMyHealSpell(event, line)
	a.onHealSpellCrit(line)
	a.onHealSpell(event, line)

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

func (a *DPS) onMySpellCast(line string) {
	match := spellNameRegex.FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	myLastSpellName = match[1]
}

func (a *DPS) onMyHealSpellCrit(line string) {
	match := myHealCritRegex.FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	myLastCritHeal = amount
}

func (a *DPS) onMyHealSpell(event time.Time, line string) {
	match := myHealRegex.FindStringSubmatch(line)
	if len(match) < 3 {
		return
	}

	amount, err := strconv.Atoi(match[2])
	if err != nil {
		return
	}

	damageEvent := &DamageEvent{
		Source: tracker.PlayerName(),
		Target: match[1],
		Type:   "heal",
		Damage: amount,
		Event:  event,
		Origin: "heal",
	}

	if myLastCritHeal == amount {
		damageEvent.IsCritical = true
		myLastCritHeal = 0
	}

	for _, fn := range a.onDamageEvent {
		fn(damageEvent)
	}

}

func (a *DPS) onMySpellCrit(line string) {
	match := directDamageCritRegex.FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	myLastCrit = amount
}

func (a *DPS) onHealSpellCrit(line string) {

	match := healCritRegex.FindStringSubmatch(line)
	if len(match) < 2 {
		return
	}

	lastOtherHealCritName = match[1]

}

func (a *DPS) onHealSpell(event time.Time, line string) {
	match := healRegex.FindStringSubmatch(line)
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

func (a *DPS) onMeleeDPS(event time.Time, line string) {
	match := meleeDamageRegex.FindStringSubmatch(line)
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
	for _, fn := range a.onDamageEvent {
		fn(damageEvent)
	}
	_, ok := a.damageEvents[damageEvent.Source]
	if !ok {
		a.damageEvents[damageEvent.Source] = make([]*DamageEvent, 0)
	}

	a.damageEvents[damageEvent.Source] = append(a.damageEvents[damageEvent.Source], damageEvent)

}

func (a *DPS) onDirectDamageDPS(event time.Time, line string) {
	match := directDamageRegex.FindStringSubmatch(line)
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
	if myLastCrit == amount {
		damageEvent.IsCritical = true
		myLastCrit = 0
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
	match := dotDamageRegex.FindStringSubmatch(line)

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
