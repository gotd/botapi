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
