package common

import (
	"fmt"
	"time"
)

type DamageEvent struct {
	Category  PopupCategory
	SpellName string
	Source    string
	Target    string
	Type      string
	Damage    string
	Event     time.Time
	Origin    string
}

func (d *DamageEvent) String() string {
	return fmt.Sprintf("%+v", *d)
}
