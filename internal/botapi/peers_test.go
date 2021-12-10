package botapi

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"
)

func Test_toTDLibID(t *testing.T) {
	tests := []struct {
		name string
		p    tg.InputPeerClass
	}{
		{"User", &tg.InputPeerUser{UserID: 309570373}},
		{"Bot", &tg.InputPeerUser{UserID: 140267078}},
		{"Chat", &tg.InputPeerChat{ChatID: 365219918}},
		{"Channel", &tg.InputPeerChat{ChatID: 1228418968}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := require.New(t)

			tdlibID := toTDLibID(tt.p)
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
