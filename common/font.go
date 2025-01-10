package common

import "reflect"

type Font int

const (
	FontNotoSansRegular42 Font = iota
	FontNotoSansBold42
	FontNotoSansRegular36
	FontNotoSansBold36
)

// String returns the string representation of the font
func (f Font) String() string {
	typ := reflect.TypeOf(f)
	for i := 0; i < typ.NumMethod(); i++ {
		if typ.Method(i).Type.NumIn() == 1 && typ.Method(i).Type.In(0) == typ {
			return typ.Method(i).Name
		}
	}
	return "unknown"
}
