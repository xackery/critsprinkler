package dps

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/reporter"
	"github.com/xackery/critsprinkler/tracker"
)

var (
	zone          = "Unknown"
	parseStart    = time.Now()
	damageEvents  = make(map[string][]*common.DamageEvent)
	onDamageEvent []func(*common.DamageEvent)
)

var (
	myLastSpellCrit       = 0
	myLastSpellCritName   = ""
	myLastSpellName       = ""
	myLastHealCrit        = 0
	myLastMeleeCrit       = 0
	lastOtherHealCritName = ""
)

func New() error {
	err := tracker.Subscribe(onLine)
	if err != nil {
		return fmt.Errorf("tracker subscribe: %w", err)
	}

	err = tracker.SubscribeToZoneEvent(onZone)
	if err != nil {
		return fmt.Errorf("tracker subscribe to zone: %w", err)
	}
	return nil
}

func onLine(event time.Time, line string) {
	onMySpellCrit(line)
	//onMyMelee(event, line)
	onMyMeleeCrit(line)
	onMyMeleeSlay(line)
	onMyMeleeFrenzy(event, line)
	onMyMeleeCleaving(line)
	onMelee(event, line)
	onMeleeMiss(event, line)
	onMyMeleeMiss(event, line)
	onSpellCast(event, line)
	onSpellInterrupt(event, line)
	onSpellFizzle(event, line)
	onSpellHit(event, line)
	onMyHealCrit(line)
	onHealCrit(line)
	onHeal(event, line)
	onRune(event, line)
	onDeath(event, line)
	dumpDPS(event)
}

func onZone(event time.Time, zoneName string) {
	zone = zoneName

	dumpDPS(event)
}

func dumpDPS(event time.Time) {
	//dpsPerSec := float64(totalDPSGained) / time.Since(parseStart).Seconds()
	//dpsPerHour := dpsPerSec * 3600

	if zone == "The Bazaar" {
		return
	}

	if len(damageEvents) == 0 {
		//fmt.Println("No damage events to report")
		return
	}

	//fmt.Println(len(damageEvents), "events to report")
	type dpsReport struct {
		total    int
		maxMelee int
		maxSpell int
	}

	damageTotals := make(map[string]dpsReport)

	tmpDamageEvents := make(map[string][]*common.DamageEvent)

	for name, dmgEvents := range damageEvents {

		for _, dmgEvent := range dmgEvents {

			// skip any events older than 60s
			if event.Sub(dmgEvent.Event).Seconds() > 60 {
				continue
			}

			if len(tmpDamageEvents[name]) == 0 {
				tmpDamageEvents[name] = make([]*common.DamageEvent, 0)
			}

			tmpDamageEvents[name] = append(tmpDamageEvents[name], dmgEvent)
			dps, ok := damageTotals[name]
			if !ok {
				dps = dpsReport{}
			}

			val, err := strconv.Atoi(dmgEvent.Damage)
			if err != nil {
				fmt.Println("atoi", dmgEvent.Damage, err)
				continue
			}

			dps.total += val
			if dmgEvent.Origin == "melee" {
				if dps.maxMelee < val {
					dps.maxMelee = val
				}
			} else {
				if dps.maxSpell < val {
					dps.maxSpell = val
				}
			}
			damageTotals[name] = dps
		}
	}

	damageEvents = tmpDamageEvents

	//fmt.Println(len(damageEvents), "events to report after filtering")

}

func onMyHealCrit(line string) {
	match, ok := easyParse(line, `\] You perform an exceptional heal! \((.*)\)`, 1)
	if !ok {
		return
	}

	amount, err := strconv.Atoi(match[0])
	if err != nil {
		return
	}

	myLastHealCrit = amount
}

func onMySpellCrit(line string) {
	match, ok := easyParse(line, `\] You deliver a critical blast! \((.*)\) \((.*)\)`, 2)
	if !ok {
		return
	}

	amount, err := strconv.Atoi(match[0])
	if err != nil {
		return
	}

	myLastSpellCrit = amount
	myLastSpellName = match[1]
}

func onMyMeleeCrit(line string) {
	match, ok := easyParse(line, `\] (.*) scores a critical hit! \((.*)\)`, 2)
	if !ok {
		return
	}

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		fmt.Println("atoi", line, match[2], err)
		return
	}

	if match[0] != tracker.PlayerName() {
		return
	}

	myLastMeleeCrit = amount
}

func onMyMeleeSlay(line string) {
	match, ok := easyParse(line, `\] (.*) holy blade cleanses (.) target!\((.*)\)`, 3)
	if !ok {
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

func onMyMeleeFrenzy(event time.Time, line string) {
	match, ok := easyParse(line, `\] You frenzy on (.*) for (.*) points of damage.`, 2)
	if !ok {
		return
	}

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	source := tracker.PlayerName()
	target := match[0]

	category := common.PopupCategoryMeleeHitOut
	damageEvent := &common.DamageEvent{
		Category: category,
		Source:   source,
		Target:   target,
		Type:     "frenzy",
		Damage:   fmt.Sprintf("%d", amount),
		Event:    event,
		Origin:   "melee",
	}

	if lastOtherHealCritName == match[1] {
		lastOtherHealCritName = ""
	}

	if source == tracker.PlayerName() && myLastMeleeCrit == amount && myLastMeleeCrit > 0 {
		damageEvent.Category = common.PopupCategoryMeleeCritOut
		myLastMeleeCrit = 0
	}

	for _, fn := range onDamageEvent {
		fn(damageEvent)
	}
}

func onMyMeleeCleaving(line string) {
	match, ok := easyParse(line, `\] (.*) lands a Cleaving Blow! \((.*)\)`, 2)
	if !ok {
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

func onHealCrit(line string) {
	match, ok := easyParse(line, `\] (.*) performs an exceptional heal! \((.*)\)`, 2)
	if !ok {
		return
	}

	if match[0] != tracker.PlayerName() {
		return
	}

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	myLastHealCrit = amount

}

func onHeal(event time.Time, line string) {
	match, ok := easyParse(line, `\] (.*) has healed (.*) for (.*) points of damage. \((.*)\)`, 4)
	if !ok {
		return
	}

	amount, err := strconv.Atoi(match[2])
	if err != nil {
		fmt.Println("atoi", line, match[2], err)
		return
	}

	category := common.PopupCategoryHealHitOut
	source := match[0]
	target := match[1]

	if strings.EqualFold(target, "itself") || strings.EqualFold(target, "himself") || strings.EqualFold(target, "herself") {
		target = source
		category = common.PopupCategoryHealHitIn
	}

	damageEvent := &common.DamageEvent{
		Category:  category,
		Source:    source,
		Target:    target,
		SpellName: match[3],
		Type:      "heal",
		Damage:    fmt.Sprintf("%d", amount),
		Event:     event,
		Origin:    "heal",
	}

	if source == tracker.PlayerName() && myLastHealCrit == amount && myLastHealCrit > 0 {
		if category == common.PopupCategoryHealHitOut {
			category = common.PopupCategoryHealCritOut
		}
		if category == common.PopupCategoryHealHitIn {
			category = common.PopupCategoryHealCritIn
		}

		myLastHealCrit = 0
	}

	for _, fn := range onDamageEvent {
		fn(damageEvent)
	}

}

func onRune(event time.Time, line string) {
	match, ok := easyParse(line, `\] (.*) has shielded (.*) from (.*) points of damage. \((.*)\)`, 4)
	if !ok {
		return
	}

	amount, err := strconv.Atoi(match[2])
	if err != nil {
		return
	}

	category := common.PopupCategoryRuneHitOut
	source := match[0]
	if source == "you" {
		source = tracker.PlayerName()
	}

	target := match[1]
	if strings.EqualFold(target, "itself") || strings.EqualFold(target, "himself") || strings.EqualFold(target, "herself") {
		target = source
		category = common.PopupCategoryRuneHitIn
	}

	damageEvent := &common.DamageEvent{
		Category:  category,
		Source:    source,
		Target:    target,
		Type:      "heal",
		Damage:    fmt.Sprintf("%d", amount),
		Event:     event,
		Origin:    "heal",
		SpellName: match[3],
	}

	for _, fn := range onDamageEvent {
		fn(damageEvent)
	}
}

func onMyMelee(event time.Time, line string) {
	match, ok := easyParse(line, `\] You (.*) for (.*) points of damage.`, 2)
	if !ok {
		return
	}

	chunk := match[0]
	firstSpacePos := strings.Index(chunk, " ")
	hitAdj := chunk[:firstSpacePos]
	target := chunk[firstSpacePos+1:]

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	pickedAdj := ""
	for _, adj := range adjectives {
		if hitAdj == adj {
			pickedAdj = adj
			break
		}
	}
	if pickedAdj == "" {
		fmt.Println("Missed parsing melee hit", hitAdj, target, chunk, line)
		return
	}

	source := tracker.PlayerName()
	if strings.EqualFold(target, "you") {
		target = tracker.PlayerName()
	}

	category := common.PopupCategoryMeleeHitOut
	damageEvent := &common.DamageEvent{
		Category: category,
		Source:   source,
		Target:   target,
		Type:     strings.TrimSpace(pickedAdj),
		Damage:   fmt.Sprintf("%d", amount),
		Event:    event,
		Origin:   "melee",
	}

	isCrit := false
	if source == tracker.PlayerName() && amount >= myLastMeleeCrit && myLastMeleeCrit > 0 {
		damageEvent.Category = common.PopupCategoryMeleeCritOut
		myLastMeleeCrit = 0
		isCrit = true
	}

	for _, fn := range onDamageEvent {
		fn(damageEvent)
	}

	_, ok = damageEvents[damageEvent.Source]
	if !ok {
		damageEvents[damageEvent.Source] = make([]*common.DamageEvent, 0)
	}

	damageEvents[damageEvent.Source] = append(damageEvents[damageEvent.Source], damageEvent)

	reporter.AttackEvent(source, 0, target, 0, &reporter.Attack{
		Event:   event,
		HitName: pickedAdj,
		Value:   amount,
		IsCrit:  isCrit,
	})
}

func onMelee(event time.Time, line string) {
	match, ok := easyParse(line, `\] (.*) for (.*) points of damage.`, 2)
	if !ok {
		return
	}

	chunks := strings.Split(match[0], " ")
	if len(chunks) < 2 {
		return
	}

	adjIndex := -1
	for i := 0; i < len(chunks); i++ {
		for _, adj := range adjectives {
			if adj == chunks[i] {
				adjIndex = i
				break
			}
		}
		if adjIndex > -1 {
			break
		}
	}

	if adjIndex == -1 {
		return
	}

	source := strings.TrimSpace(strings.Join(chunks[:adjIndex], " "))
	hitAdj := chunks[adjIndex]
	target := strings.TrimSpace(strings.Join(chunks[adjIndex+1:], " "))

	category := common.PopupCategoryMeleeHitOut
	if strings.EqualFold(target, "you") {
		target = tracker.PlayerName()
		category = common.PopupCategoryMeleeHitIn
	}
	if strings.EqualFold(source, "you") {
		source = tracker.PlayerName()
		category = common.PopupCategoryMeleeHitOut
	}

	amount, err := strconv.Atoi(match[1])
	if err != nil {
		return
	}

	// if strings.Contains(source, "`s doppleganger") {
	// 	source = strings.ReplaceAll(source, "`s doppleganger", "")
	// }
	// if strings.Contains(target, "`s doppleganger") {
	// 	target = strings.ReplaceAll(target, "`s doppleganger", "")
	// }

	damageEvent := &common.DamageEvent{
		Category: category,
		Source:   source,
		Target:   target,
		Type:     hitAdj,
		Damage:   fmt.Sprintf("%d", amount),
		Event:    event,
		Origin:   "melee",
	}

	isCrit := false
	if source == tracker.PlayerName() && amount >= myLastMeleeCrit && myLastMeleeCrit > 0 {
		damageEvent.Category = common.PopupCategoryMeleeCritOut
		myLastMeleeCrit = 0
		isCrit = true
	}

	for _, fn := range onDamageEvent {
		fn(damageEvent)
	}

	_, ok = damageEvents[damageEvent.Source]
	if !ok {
		damageEvents[damageEvent.Source] = make([]*common.DamageEvent, 0)
	}

	damageEvents[damageEvent.Source] = append(damageEvents[damageEvent.Source], damageEvent)

	reporter.AttackEvent(source, 0, target, 0, &reporter.Attack{
		Event:   event,
		HitName: hitAdj,
		Value:   amount,
		IsCrit:  isCrit,
	})
}

func onMeleeMiss(event time.Time, line string) {
	match, ok := easyParse(line, `\] (.*) tries to (.*), but (.*)!`, 3)
	if !ok {
		return
	}

	source := match[0]

	chunks := strings.Split(match[1], " ")
	if len(chunks) < 2 {
		return
	}

	adjIndex := -1
	for i := 0; i < len(chunks); i++ {
		for _, adj := range adjectives {
			if adj == chunks[i] {
				adjIndex = i
				break
			}
		}
		if adjIndex > -1 {
			break
		}
	}

	if adjIndex == -1 {
		return
	}

	hitAdj := chunks[adjIndex]
	target := strings.TrimSpace(strings.Join(chunks[adjIndex+1:], " "))

	missName := match[2]
	if missName == "misses" {
		missName = "miss"
	} else {
		missName = ""
		opts := []string{"parry", "dodge", "block", "shield block", "riposte", "absorb"}
		for _, opt := range opts {
			if strings.Contains(missName, opt) {
				if opt == "absorb" {
					missName = "rune"
				}
				missName = opt
				break
			}
		}
		if missName == "" {
			missName = strings.ReplaceAll(match[2], "YOU ", "")
		}
	}

	category := common.PopupCategoryMeleeMissIn
	if strings.EqualFold(target, "you") {
		target = tracker.PlayerName()
		category = common.PopupCategoryMeleeMissIn
	}
	if strings.EqualFold(source, "you") {
		source = tracker.PlayerName()
		category = common.PopupCategoryMeleeMissOut
	}

	// if strings.Contains(source, "`s doppleganger") {
	// 	source = strings.ReplaceAll(source, "`s doppleganger", "")
	// }
	// if strings.Contains(target, "`s doppleganger") {
	// 	target = strings.ReplaceAll(target, "`s doppleganger", "")
	// }

	damageEvent := &common.DamageEvent{
		Category: category,
		Source:   source,
		Target:   target,
		Type:     hitAdj,
		Damage:   missName,
		Event:    event,
		Origin:   "melee",
	}

	for _, fn := range onDamageEvent {
		fn(damageEvent)
	}

	// _, ok = damageEvents[damageEvent.Source]
	// if !ok {
	// 	damageEvents[damageEvent.Source] = make([]*common.DamageEvent, 0)
	// }

	// damageEvents[damageEvent.Source] = append(damageEvents[damageEvent.Source], damageEvent)

	result := reporter.AttackMiss
	switch missName {
	case "parry":
		result = reporter.AttackParry
	case "dodge":
		result = reporter.AttackDodge
	case "block":
		result = reporter.AttackBlock
	case "shield block":
		result = reporter.AttackShieldBlock
	case "riposte":
		result = reporter.AttackRiposte
	case "absorb":
		result = reporter.AttackAbsorb
	}

	reporter.AttackEvent(source, 0, target, 0, &reporter.Attack{
		Event:   event,
		HitName: hitAdj,
		Value:   0,
		Result:  result,
		IsCrit:  false,
	})
}

func onMyMeleeMiss(event time.Time, line string) {
	match, ok := easyParse(line, `\] You try to (.*), but (.*)!`, 2)
	if !ok {
		return
	}

	// You try to slash Bob Barker, but Bob Barker dodges!

	source := tracker.PlayerName()

	chunks := strings.Split(match[0], " ")
	if len(chunks) < 2 {
		return
	}

	adjIndex := -1
	for i := 0; i < len(chunks); i++ {
		for _, adj := range adjectives {
			if adj == chunks[i] {
				adjIndex = i
				break
			}
		}
		if adjIndex > -1 {
			break
		}
	}

	if adjIndex == -1 {
		return
	}

	hitAdj := chunks[adjIndex]
	target := strings.TrimSpace(strings.Join(chunks[adjIndex+1:], " "))

	missName := match[1]
	missName = strings.ReplaceAll(missName, target, "")

	category := common.PopupCategoryMeleeMissOut
	if strings.EqualFold(target, "you") {
		target = tracker.PlayerName()
		category = common.PopupCategoryMeleeMissIn
	}

	// if strings.Contains(source, "`s doppleganger") {
	// 	source = strings.ReplaceAll(source, "`s doppleganger", "")
	// }
	// if strings.Contains(target, "`s doppleganger") {
	// 	target = strings.ReplaceAll(target, "`s doppleganger", "")
	// }

	damageEvent := &common.DamageEvent{
		Category: category,
		Source:   source,
		Target:   target,
		Type:     hitAdj,
		Damage:   missName,
		Event:    event,
		Origin:   "melee",
	}

	for _, fn := range onDamageEvent {
		fn(damageEvent)
	}

	// _, ok = damageEvents[damageEvent.Source]
	// if !ok {
	// 	damageEvents[damageEvent.Source] = make([]*common.DamageEvent, 0)
	// }

	// damageEvents[damageEvent.Source] = append(damageEvents[damageEvent.Source], damageEvent)
}

func onSpellCast(event time.Time, line string) {
	match, ok := easyParse(line, `\] You begin to cast (.*).`, 1)
	if !ok {
		return
	}

	myLastSpellName = match[0]
}

func onSpellInterrupt(event time.Time, line string) {
	_, ok := easyParse(line, `\] Your spell is interrupted.`, 0)
	if !ok {
		return
	}

	reporter.CastEvent(tracker.PlayerName(), 0, &reporter.Cast{
		Event:     event,
		SpellName: myLastSpellName,
		Result:    reporter.CastInterrupted,
		Value:     0,
		IsCrit:    false,
	})
	myLastSpellName = ""
}

func onDeath(event time.Time, line string) {
	match, ok := easyParse(line, `\] (.*) has been killed by (.*)!`, 2)
	if !ok {
		return
	}

	source := match[1]
	target := match[0]

	reporter.DeathEvent(source, 0, target, 0, event)
}

func onSpellFizzle(event time.Time, line string) {
	_, ok := easyParse(line, `\] Your spell fizzles!`, 0)
	if !ok {
		return
	}

	reporter.CastEvent(tracker.PlayerName(), 0, &reporter.Cast{
		Event:     event,
		SpellName: myLastSpellName,
		Result:    reporter.CastFizzle,
		Value:     0,
		IsCrit:    false,
	})
	myLastSpellName = ""
}

func onSpellHit(event time.Time, line string) {
	match, ok := easyParse(line, `\] (.*) hit (.*) for (.*) points of non-melee damage. \((.*)\)`, 4)
	if !ok {
		return
	}

	amount, err := strconv.Atoi(match[2])
	if err != nil {
		return
	}

	source := match[0]
	target := match[1]

	damageEvent := &common.DamageEvent{
		Category:  common.PopupCategorySpellHitOut,
		Source:    source,
		Type:      "hit",
		Target:    target,
		Damage:    fmt.Sprintf("%d", amount),
		SpellName: match[3],
		Event:     event,
		Origin:    "direct",
	}

	if source == tracker.PlayerName() {
		damageEvent.Category = common.PopupCategorySpellHitOut
	}
	isCrit := false
	if source == tracker.PlayerName() && myLastSpellCrit == amount {
		damageEvent.Category = common.PopupCategorySpellCritOut
		myLastSpellCrit = 0
		isCrit = true
	}
	if target == tracker.PlayerName() {
		damageEvent.Category = common.PopupCategorySpellHitIn
	}
	for _, fn := range onDamageEvent {
		fn(damageEvent)
	}

	_, ok = damageEvents[damageEvent.Source]
	if !ok {
		damageEvents[damageEvent.Source] = make([]*common.DamageEvent, 0)
	}

	damageEvents[damageEvent.Source] = append(damageEvents[damageEvent.Source], damageEvent)
	reporter.CastEvent(source, 0, &reporter.Cast{
		Event:     event,
		SpellName: match[3],
		Result:    reporter.CastSuccess,
		Value:     amount,
		IsCrit:    isCrit,
	})

}

func SubscribeToDamageEvent(fn func(*common.DamageEvent)) error {
	onDamageEvent = append(onDamageEvent, fn)
	return nil
}

// easyParse is a helper function to parse a line with a regex and return the match
func easyParse(line string, regex string, size int) ([]string, bool) {
	match := regexp.MustCompile(regex).FindStringSubmatch(line)
	if len(match) < 1 {
		return nil, false
	}
	match = match[1:]
	if len(match) != size {
		return nil, false
	}

	return match, true
}

var adjectives = []string{
	"mauls",
	"maul",
	"bites",
	"bite",
	"claws",
	"claw",
	"gores",
	"gore",
	"stings",
	"slices",
	"slice",
	"sting",
	"smashes",
	"smash",
	"rend",
	"rends",
	"slash",
	"slashes",
	"punch",
	"punches",
	"hit",
	"hits",
	//"You",
	//"yourself",
	//"YOU",
	//"himself",
	//"herself",
	//"itself",
	"crush",
	"crushes",
	"pierce",
	"pierces",
	"kick",
	"kicks",
	"strike",
	"strikes",
	"backstab",
	"backstabs",
	"bash",
	"bashes",
}
