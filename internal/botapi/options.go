package botapi

import "go.uber.org/zap"

// Options is options of BotAPI.
type Options struct {
	Debug  bool
	Logger *zap.Logger
}

func (o *Options) setDefaults() {
	if o.Logger == nil {
		o.Logger = zap.NewNop()
	}
}
