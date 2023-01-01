package botapi

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgmock"

	"github.com/gotd/botapi/internal/oas"
)

func TestBotAPI_SendChatAction(t *testing.T) {
	progress := 0
	tests := []struct {
		action  string
		expect  tg.SendMessageActionClass
		wantErr bool
	}{
		{"cancel", &tg.SendMessageCancelAction{}, false},
		{"typing", &tg.SendMessageTypingAction{}, false},
		{"record_video", &tg.SendMessageRecordVideoAction{}, false},
		{"upload_video", &tg.SendMessageUploadVideoAction{Progress: progress}, false},
		{"record_audio", &tg.SendMessageRecordAudioAction{}, false},
		{"record_voice", &tg.SendMessageRecordAudioAction{}, false},
		{"upload_audio", &tg.SendMessageUploadVideoAction{Progress: progress}, false},
		{"upload_voice", &tg.SendMessageUploadVideoAction{Progress: progress}, false},
		{"upload_photo", &tg.SendMessageUploadPhotoAction{Progress: progress}, false},
		{"upload_document", &tg.SendMessageUploadDocumentAction{Progress: progress}, false},
		{"choose_sticker", &tg.SendMessageChooseStickerAction{}, false},
		{"pick_up_location", &tg.SendMessageGeoLocationAction{}, false},
		{"find_location", &tg.SendMessageGeoLocationAction{}, false},
		{"record_video_note", &tg.SendMessageRecordRoundAction{}, false},
		{"upload_video_note", &tg.SendMessageUploadRoundAction{Progress: progress}, false},
		{"", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			ctx := context.Background()
			testWithCache(t, func(a *require.Assertions, mock *tgmock.Mock, api *BotAPI) {
				if !tt.wantErr {
					mock.ExpectCall(&tg.MessagesSetTypingRequest{
						Peer:   &tg.InputPeerChat{ChatID: testChat().ID},
						Action: tt.expect,
					}).ThenTrue()
				}
				_, err := api.SendChatAction(ctx, &oas.SendChatAction{
					ChatID: oas.NewInt64ID(testChatID()),
					Action: tt.action,
				})
				if tt.wantErr {
					a.Error(err)
					return
				}
				a.NoError(err)
			})
		})
	}
}
