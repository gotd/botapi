package botapi

import (
	"github.com/gotd/td/constant"
	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func (b *BotAPI) encodeFileID(f fileid.FileID) (fileID, fileUniqueID string) {
	fileID, _ = fileid.EncodeFileID(f)
	// TODO(tdakkota): generate unique id
	fileUniqueID = "todo"
	return fileID, fileUniqueID
}

func (b *BotAPI) profilePhotoFileID(
	id constant.TDLibPeerID,
	accessHash int64,
	photo fileid.ChatPhoto,
	big bool,
) (fileID, fileUniqueID string) {
	return b.encodeFileID(fileid.FromChatPhoto(id, accessHash, photo, big))
}

func (b *BotAPI) setChatPhoto(
	id constant.TDLibPeerID,
	accessHash int64,
	from tg.ChatPhotoClass,
) (to oas.OptChatPhoto) {
	p, ok := from.(*tg.ChatPhoto)
	if !ok {
		return
	}

	smallFileID, smallUniqueFileID := b.profilePhotoFileID(id, accessHash, p, false)
	bigFileID, bigUniqueFileID := b.profilePhotoFileID(id, accessHash, p, true)
	return oas.NewOptChatPhoto(oas.ChatPhoto{
		SmallFileID:       smallFileID,
		SmallFileUniqueID: smallUniqueFileID,
		BigFileID:         bigFileID,
		BigFileUniqueID:   bigUniqueFileID,
	})
}

func (b *BotAPI) setUserPhoto(
	id constant.TDLibPeerID,
	accessHash int64,
	from tg.UserProfilePhotoClass,
) (to oas.OptChatPhoto) {
	p, ok := from.(*tg.UserProfilePhoto)
	if !ok {
		return
	}

	smallFileID, smallUniqueFileID := b.profilePhotoFileID(id, accessHash, p, false)
	bigFileID, bigUniqueFileID := b.profilePhotoFileID(id, accessHash, p, true)
	return oas.NewOptChatPhoto(oas.ChatPhoto{
		SmallFileID:       smallFileID,
		SmallFileUniqueID: smallUniqueFileID,
		BigFileID:         bigFileID,
		BigFileUniqueID:   bigUniqueFileID,
	})
}
