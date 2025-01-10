// bubble is a package that handles the bubble popup system
package bubble

import (
	"fmt"

	"github.com/xackery/critsprinkler/common"
	"github.com/xackery/critsprinkler/dps"
)

var (
	damageEventChan = make(chan *common.DamageEvent, 10000)
	spawns          = []*common.DamageEvent{}
)

func New() error {
	err := dps.SubscribeToDamageEvent(DamageEvent)
	if err != nil {
		return fmt.Errorf("dps subscribe to damage event: %w", err)
	}
	return nil
}

// DamageEvent is called to add a damage event to the popup queue
func DamageEvent(event *common.DamageEvent) {
	damageEventChan <- event
}

// Update is called by the engine
func Update() {
	for i := 0; i < 60; i++ {
		isSuccess := false
		select {
		case event := <-damageEventChan:
			err := spawn(event)
			if err != nil {
				fmt.Println("spawn error:", err)
			}
			isSuccess = true
		default:
		}
		if !isSuccess {
			break
		}

	}
}

func spawn(event *common.DamageEvent) error {
	spawns = append(spawns, event)
	return nil
}

// Spawns is fetched by popup to get off queue
func Spawns() []*common.DamageEvent {
	out := spawns
	spawns = []*common.DamageEvent{}
	return out
}
