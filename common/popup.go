package common

import "reflect"

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

	PopupCategoryTotalDamageOut
	PopupCategoryTotalDamageIn
	PopupCategoryTotalHealOut
	PopupCategoryTotalHealIn

	PopupCategoryMax
)

func (e PopupCategory) String() string {
	switch e {
	case PopupCategoryMeleeCritOut:
		return "Melee Crit Out"
	case PopupCategoryMeleeHitOut:
		return "Melee Hit Out"
	case PopupCategoryMeleeMissOut:
		return "Melee Miss Out"
	case PopupCategoryMeleeCritIn:
		return "Melee Crit In"
	case PopupCategoryMeleeHitIn:
		return "Melee Hit In"
	case PopupCategoryMeleeMissIn:
		return "Melee Miss In"
	case PopupCategorySpellCritOut:
		return "Spell Crit Out"
	case PopupCategorySpellHitOut:
		return "Spell Hit Out"
	case PopupCategorySpellMissOut:
		return "Spell Miss Out"
	case PopupCategorySpellCritIn:
		return "Spell Crit In"
	case PopupCategorySpellHitIn:
		return "Spell Hit In"
	case PopupCategorySpellMissIn:
		return "Spell Miss In"
	case PopupCategoryHealCritOut:
		return "Heal Crit Out"
	case PopupCategoryHealHitOut:
		return "Heal Hit Out"
	case PopupCategoryHealCritIn:
		return "Heal Crit In"
	case PopupCategoryHealHitIn:
		return "Heal Hit In"
	case PopupCategoryRuneHitOut:
		return "Rune Hit Out"
	case PopupCategoryRuneHitIn:
		return "Rune Hit In"
	case PopupCategoryTotalDamageOut:
		return "Total Damage Out"
	case PopupCategoryTotalDamageIn:
		return "Total Damage In"
	case PopupCategoryTotalHealOut:
		return "Total Heal Out"
	case PopupCategoryTotalHealIn:
		return "Total Heal In"
	}
	return "unknown"

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

func (e Direction) String() string {
	typ := reflect.TypeOf(e)
	for i := 0; i < typ.NumMethod(); i++ {
		if typ.Method(i).Type.NumIn() == 1 && typ.Method(i).Type.In(0) == typ {
			return typ.Method(i).Name
		}
	}
	return "unknown"
}

func IsTotalDamageIn(category PopupCategory) bool {
	return category == PopupCategoryMeleeCritIn ||
		category == PopupCategoryMeleeHitIn ||
		category == PopupCategorySpellCritIn ||
		category == PopupCategorySpellHitIn ||
		category == PopupCategorySpellMissIn
}

func IsTotalDamageOut(category PopupCategory) bool {
	return category == PopupCategoryMeleeCritOut ||
		category == PopupCategoryMeleeHitOut ||
		category == PopupCategorySpellCritOut ||
		category == PopupCategorySpellHitOut ||
		category == PopupCategorySpellMissOut
}

func IsTotalHealIn(category PopupCategory) bool {
	return category == PopupCategoryHealCritIn ||
		category == PopupCategoryHealHitIn ||
		category == PopupCategoryRuneHitIn
}

func IsTotalHealOut(category PopupCategory) bool {
	return category == PopupCategoryHealCritOut ||
		category == PopupCategoryHealHitOut ||
		category == PopupCategoryRuneHitOut
}
