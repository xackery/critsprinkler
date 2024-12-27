package common

type PopupCategory int

const (
	PopupCategoryMeleeCritOut PopupCategory = iota
	PopupCategoryMeleeHitOut
	PopupCategoryMeleeMissOut
	PopupCategoryMeleeCritIn
	PopupCategoryMeleeHitIn
	PopupCategoryMeleeMissIn

	PopupCategorySpellCritOut
	PopupCategorySpellHitOut
	PopupCategorySpellMissOut
	PopupCategorySpellCritIn
	PopupCategorySpellHitIn
	PopupCategorySpellMissIn

	PopupCategoryHealCritOut
	PopupCategoryHealHitOut
	PopupCategoryHealCritIn
	PopupCategoryHealHitIn

	//PopupCategoryRuneCritOut
	PopupCategoryRuneHitOut
	//PopupCategoryRuneCritIn
	PopupCategoryRuneHitIn

	PopupCategoryGlobalCritOut
	PopupCategoryGlobalHitOut
	PopupCategoryGlobalMissOut
	PopupCategoryGlobalCritIn
	PopupCategoryGlobalHitIn
	PopupCategoryGlobalMissIn

	PopupCategoryMax
)

// String returns the string representation of the PopupCategory.
func (c PopupCategory) String() string {
	switch c {
	case PopupCategoryMeleeCritOut:
		return "MeleeCritOut"
	case PopupCategoryMeleeHitOut:
		return "MeleeHitOut"
	case PopupCategoryMeleeMissOut:
		return "MeleeMissOut"
	case PopupCategoryMeleeCritIn:
		return "MeleeCritIn"
	case PopupCategoryMeleeHitIn:
		return "MeleeHitIn"
	case PopupCategoryMeleeMissIn:
		return "MeleeMissIn"
	case PopupCategorySpellCritOut:
		return "SpellCritOut"
	case PopupCategorySpellHitOut:
		return "SpellHitOut"
	case PopupCategorySpellMissOut:
		return "SpellMissOut"
	case PopupCategorySpellCritIn:
		return "SpellCritIn"
	case PopupCategorySpellHitIn:
		return "SpellHitIn"
	case PopupCategorySpellMissIn:
		return "SpellMissIn"
	case PopupCategoryHealCritOut:
		return "HealCritOut"
	case PopupCategoryHealHitOut:
		return "HealHitOut"
	case PopupCategoryHealCritIn:
		return "HealCritIn"
	case PopupCategoryHealHitIn:
		return "HealHitIn"
	case PopupCategoryRuneHitOut:
		return "RuneHitOut"
	case PopupCategoryRuneHitIn:
		return "RuneHitIn"
	}
	return "Unknown"
}

type Direction int

const (
	DirectionUp Direction = iota
	DirectionUpRight
	DirectionRight
	DirectionDownRight
	DirectionDown
	DirectionDownLeft
	DirectionLeft
	DirectionUpLeft
)

func (d Direction) String() string {
	switch d {
	case DirectionUp:
		return "Up"
	case DirectionUpRight:
		return "UpRight"
	case DirectionRight:
		return "Right"
	case DirectionDownRight:
		return "DownRight"
	case DirectionDown:
		return "Down"
	case DirectionDownLeft:
		return "DownLeft"
	case DirectionLeft:
		return "Left"
	case DirectionUpLeft:
		return "UpLeft"
	}
	return "Unknown"
}
