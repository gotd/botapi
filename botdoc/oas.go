package botdoc

import (
	"github.com/ogen-go/ogen"
)

func (a API) OAS() *ogen.Spec {
	spec := &ogen.Spec{
		OpenAPI: "3.0.0",
		Info: ogen.Info{
			Title:          "Telegram Bot API",
			TermsOfService: "https://telegram.org/tos",
			Description:    "API for Telegram bots",
			Version:        "5.3",
		},
		Servers: []ogen.Server{
			{
				Description: "production",
				URL:         "https://api.telegram.org/",
			},
		},
	}

	return spec
}
