package reporter

import (
	"fmt"
	"sync"
	"time"
)

type CastResult int

const (
	CastSuccess CastResult = iota
	CastResist
	CastImmune
	CastInterrupted
	CastFizzle
)

type AttackResult int

const (
	AttackSuccess AttackResult = iota
	AttackMiss
	AttackParry
	AttackDodge
	AttackBlock
	AttackImmune
	AttackRiposte
	AttackAbsorb
	AttackShieldBlock
)

var (
	instance *Reporter
	mux      = sync.RWMutex{}
)

type Reporter struct {
	OngoingBattles  []*Battle
	FinishedBattles []*Battle
}

type AttackSummary struct {
	SourceName string
	TotalHits  int
	TotalCrits int
	TotalDmg   int
}

// A battle is a group of events
type Battle struct {
	Start          time.Time // Start time of the battle
	KillerName     string    // Name of the killer
	KillerID       int       // ID of the killer
	LastEvent      time.Time // Last event time
	End            time.Time // End time of the battle
	Target         *Mob      // Target of the battle
	TargetCorpseID int       // Corpse ID of the target
	Mobs           []*Mob    // Mobs in the battle
}

// Mob represents a player or NPC
type Mob struct {
	Name    string
	ID      int
	Casts   []*Cast
	Attacks []*Attack
}

// Cast is an attempted cast
type Cast struct {
	Event     time.Time  // Time of the cast
	SpellName string     // Name of the spell
	Result    CastResult // Result of the cast
	Value     int        // Value of the cast
	IsCrit    bool       // Is the cast a critical hit
}

// Attack is an attempted attack
type Attack struct {
	Event   time.Time    // Time of the attack
	HitName string       // Name of the hit
	Result  AttackResult // Result of the attack
	Value   int          // Value of the attack
	IsCrit  bool         // Is the attack a critical hit
}

// New creates a new reporter
func New() (*Reporter, error) {
	mux.Lock()
	defer mux.Unlock()
	instance = &Reporter{}
	return instance, nil
}

// CastEvent is called when a cast event occurs
func CastEvent(sourceName string, sourceID int, cast *Cast) error {
	mux.Lock()
	defer mux.Unlock()

	if instance == nil {
		return fmt.Errorf("reporter not initialized")
	}

	if sourceName == "" {
		return fmt.Errorf("sourceName cannot be empty")
	}

	battle := instance.battleFetchOrStart(sourceName, sourceID, cast.Event)
	if battle == nil {
		return fmt.Errorf("fetch battle returned nil?")
	}

	battle.LastEvent = cast.Event

	var source *Mob
	for _, m := range battle.Mobs {
		if m.ID != 0 && m.ID != sourceID {
			continue
		}

		if m.Name != sourceName {
			continue
		}

		source = m
		break
	}

	if source == nil {
		source = &Mob{
			Name: sourceName,
			ID:   sourceID,
		}
		battle.Mobs = append(battle.Mobs, source)
	}

	source.Casts = append(source.Casts, cast)

	return nil
}

// AttackEvent is called when an attack event occurs
func AttackEvent(sourceName string, sourceID int, targetName string, targetID int, attack *Attack) error {
	mux.Lock()
	defer mux.Unlock()

	if instance == nil {
		return fmt.Errorf("reporter not initialized")
	}

	if targetName == "" && sourceName == "" && sourceID == 0 && targetID == 0 {
		return fmt.Errorf("sourceName and targetName cannot be empty")
	}

	battle := instance.battleFetchOrStart(targetName, targetID, attack.Event)
	if battle == nil {
		return fmt.Errorf("fetch battle returned nil?")
	}

	battle.LastEvent = attack.Event

	var source *Mob
	for _, m := range battle.Mobs {
		if m.ID != 0 && m.ID != sourceID {
			continue
		}

		if m.Name != sourceName {
			continue
		}

		source = m
		break
	}

	if source == nil {
		source = &Mob{
			Name: sourceName,
			ID:   sourceID,
		}
		battle.Mobs = append(battle.Mobs, source)
	}

	source.Attacks = append(source.Attacks, attack)

	return nil
}

// DeathEvent is called when a death event occurs
func DeathEvent(targetName string, targetID int, killerName string, killerID int, event time.Time) error {
	mux.Lock()
	defer mux.Unlock()

	if instance == nil {
		return fmt.Errorf("reporter not initialized")
	}

	if targetName == "" {
		return fmt.Errorf("targetName cannot be empty")
	}

	battle := instance.battleFetchOrStart(targetName, targetID, event)
	if battle == nil {
		return fmt.Errorf("fetch battle returned nil?")
	}

	battle.KillerName = killerName
	battle.KillerID = killerID
	battle.LastEvent = event
	battle.End = event

	instance.FinishedBattles = append(instance.FinishedBattles, battle)
	instance.OngoingBattles = append(instance.OngoingBattles[:0], instance.OngoingBattles[1:]...)
	return nil
}

func (e *Reporter) battleFetchOrStart(targetName string, targetID int, event time.Time) *Battle {
	isSummaryDirty := false
	var battle *Battle
	for _, b := range e.OngoingBattles {
		if b.Target.ID != 0 && b.Target.ID != targetID {
			continue
		}

		if b.Target.Name != targetName {
			continue
		}

		if event.Sub(b.LastEvent) > 1*time.Minute {
			e.FinishedBattles = append(e.FinishedBattles, b)
			isSummaryDirty = true
			e.OngoingBattles = append(e.OngoingBattles[:0], e.OngoingBattles[1:]...)

			continue
		}

		battle = b

		break
	}

	if battle == nil {
		battle = &Battle{
			Start:     event,
			LastEvent: event,
			Target: &Mob{
				Name: targetName,
				ID:   targetID,
			},
		}
		e.OngoingBattles = append(e.OngoingBattles, battle)
	}

	if isSummaryDirty {
		e.updateSummaries()
	}

	return battle
}

func (e *Reporter) updateSummaries() {
}
