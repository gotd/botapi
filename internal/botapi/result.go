package botapi

import "github.com/gotd/botapi/internal/oas"

func resultOK(v bool) *oas.Result {
	return &oas.Result{
		Result: oas.OptBool{
			Value: v,
			Set:   v,
		},
		Ok: true,
	}
}
