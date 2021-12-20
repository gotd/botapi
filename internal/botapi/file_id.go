package botapi

import "github.com/gotd/td/fileid"

func (b *BotAPI) encodeFileID(f fileid.FileID) (fileID, fileUniqueID string) {
	fileID, _ = fileid.EncodeFileID(f)
	// TODO(tdakkota): generate unique id
	fileUniqueID = "todo"
	return fileID, fileUniqueID
}
