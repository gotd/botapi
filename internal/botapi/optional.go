package botapi

import "github.com/gotd/botapi/internal/oas"

func optString(getter func() (string, bool)) oas.OptString {
	v, ok := getter()
	if !ok {
		return oas.OptString{}
	}
	return oas.NewOptString(v)
}

func optBool(v bool) oas.OptBool {
	return oas.OptBool{
		Value: v,
		Set:   v,
	}
}
