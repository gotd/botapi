package botapi

import "github.com/gotd/botapi/internal/oas"

func optString(getter func() (string, bool)) oas.OptString {
	v, ok := getter()
	if !ok {
		return oas.OptString{}
	}
	return oas.NewOptString(v)
}

func optInt(getter func() (int, bool)) oas.OptInt {
	v, ok := getter()
	if !ok {
		return oas.OptInt{}
	}
	return oas.NewOptInt(v)
}

func optFloat64(getter func() (float64, bool)) oas.OptFloat64 {
	v, ok := getter()
	if !ok {
		return oas.OptFloat64{}
	}
	return oas.NewOptFloat64(v)
}

func trueType(v bool) oas.OptBool {
	return oas.OptBool{
		Value: v,
		Set:   v,
	}
}
