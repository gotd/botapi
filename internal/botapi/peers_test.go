package botapi

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"
)

func Test_toTDLibID(t *testing.T) {
	tests := []struct {
		name    string
		tdlibID int64
		p       tg.InputPeerClass
	}{
		{"User", 309570373, &tg.InputPeerUser{UserID: 309570373}},
		{"Bot", 140267078, &tg.InputPeerUser{UserID: 140267078}},
		{"Chat", -365219918, &tg.InputPeerChat{ChatID: 365219918}},
		{"Channel", -1001228418968, &tg.InputPeerChannel{ChannelID: 1228418968}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := require.New(t)

			tdlibID := toTDLibID(tt.p)
			a.Equal(tt.tdlibID, tdlibID)
			var mtprotoID int64
			switch t := tt.p.(type) {
			case *tg.InputPeerUser:
				mtprotoID = t.UserID
				a.True(IsUserTDLibID(tdlibID))
			case *tg.InputPeerChat:
				mtprotoID = t.ChatID
				a.True(IsChatTDLibID(tdlibID))
			case *tg.InputPeerChannel:
				mtprotoID = t.ChannelID
				a.True(IsChannelTDLibID(tdlibID))
			}
			cleanID := fromTDLibID(tdlibID)
			a.Equal(mtprotoID, cleanID)
		})
	}
}
