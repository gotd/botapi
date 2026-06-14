package botapi

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"io"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

// errInvalidFileID is returned when a file_id cannot be decoded.
func errInvalidFileID() *Error {
	return &Error{Code: 400, Description: "Bad Request: invalid file_id"}
}

// inputPhotoFromFileID builds an MTProto input photo reference from a file_id.
func inputPhotoFromFileID(fileID string) (tg.InputPhotoClass, error) {
	f, err := fileid.DecodeFileID(fileID)
	if err != nil {
		return nil, errInvalidFileID()
	}

	return &tg.InputPhoto{
		ID:            f.ID,
		AccessHash:    f.AccessHash,
		FileReference: f.FileReference,
	}, nil
}

// inputDocumentFromFileID builds an MTProto input document reference from a
// file_id (documents, video, audio, voice, stickers, animations, ...).
func inputDocumentFromFileID(fileID string) (tg.InputDocumentClass, error) {
	f, err := fileid.DecodeFileID(fileID)
	if err != nil {
		return nil, errInvalidFileID()
	}

	return &tg.InputDocument{
		ID:            f.ID,
		AccessHash:    f.AccessHash,
		FileReference: f.FileReference,
	}, nil
}

// GetFile decodes a file_id and returns a File describing it.
//
// Unlike the HTTP Bot API, this library is MTProto-native: there is no file
// server, so GetFile performs no network I/O and FilePath is left empty. Use
// DownloadFile or DownloadFileToPath to fetch the contents. FileUniqueID is
// derived locally from the file_id.
func (b *Bot) GetFile(_ context.Context, fileID string) (*File, error) {
	f, err := fileid.DecodeFileID(fileID)
	if err != nil {
		return nil, errInvalidFileID()
	}

	return &File{
		FileID:       fileID,
		FileUniqueID: fileUniqueID(f),
	}, nil
}

// DownloadFile streams the file referenced by file_id into w. It follows DC
// migration transparently. The number of bytes written is returned.
func (b *Bot) DownloadFile(ctx context.Context, fileID string, w io.Writer) (int64, error) {
	f, err := fileid.DecodeFileID(fileID)
	if err != nil {
		return 0, errInvalidFileID()
	}

	counter := &countWriter{w: w}
	if loc, ok := f.AsInputFileLocation(); ok {
		if _, err := b.client.Download(loc).Stream(ctx, counter); err != nil {
			return counter.n, asAPIError(err)
		}

		return counter.n, nil
	}

	if loc, ok := f.AsInputWebFileLocation(); ok {
		if _, err := b.client.DownloadWeb(loc).Stream(ctx, counter); err != nil {
			return counter.n, asAPIError(err)
		}

		return counter.n, nil
	}

	return 0, &Error{Code: 400, Description: "Bad Request: file_id is not downloadable"}
}

// DownloadFileToPath downloads the file referenced by file_id to a local path.
func (b *Bot) DownloadFileToPath(ctx context.Context, fileID, path string) error {
	f, err := fileid.DecodeFileID(fileID)
	if err != nil {
		return errInvalidFileID()
	}

	loc, ok := f.AsInputFileLocation()
	if !ok {
		return &Error{Code: 400, Description: "Bad Request: file_id is not downloadable"}
	}

	if _, err := b.client.Download(loc).ToPath(ctx, path); err != nil {
		return asAPIError(err)
	}

	return nil
}

// countWriter counts bytes forwarded to the wrapped writer.
type countWriter struct {
	w io.Writer
	n int64
}

func (c *countWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)

	c.n += int64(n)

	return n, err
}

// file_unique_id "unique types", matching the TDLib/Bot API scheme.
const (
	uniqueTypeWeb       = 0
	uniqueTypePhoto     = 1
	uniqueTypeDocument  = 2
	uniqueTypeSecure    = 3
	uniqueTypeEncrypted = 4
	uniqueTypeTemp      = 5
)

// fileUniqueID derives the Bot API file_unique_id from a decoded file_id.
//
// It mirrors the TDLib algorithm: a little-endian record of the unique type and
// the file's identifying fields, RLE-encoded over zero bytes and base64url
// encoded without padding. Web and document-family files (documents, video,
// audio, voice, stickers, animations, video notes) are exact; legacy photos use
// their volume/local id. Newer photo-size sources that carry no volume id fall
// back to the media id, which stays stable per file but may differ from the
// value Telegram's own server would return.
func fileUniqueID(f fileid.FileID) string {
	var buf []byte

	switch {
	case f.URL != "":
		buf = make([]byte, 4, 4+len(f.URL))
		binary.LittleEndian.PutUint32(buf, uniqueTypeWeb)

		buf = append(buf, f.URL...)
	case isPhotoType(f.Type) && f.PhotoSizeSource.VolumeID != 0:
		buf = make([]byte, 16)
		binary.LittleEndian.PutUint32(buf[0:], uniqueTypePhoto)
		binary.LittleEndian.PutUint64(buf[4:], uint64(f.PhotoSizeSource.VolumeID))
		binary.LittleEndian.PutUint32(buf[12:], uint32(int32(f.PhotoSizeSource.LocalID)))
	default:
		buf = make([]byte, 12)
		binary.LittleEndian.PutUint32(buf[0:], uint32(uniqueTypeFor(f.Type)))
		binary.LittleEndian.PutUint64(buf[4:], uint64(f.ID))
	}

	return base64.RawURLEncoding.EncodeToString(rleEncode(buf))
}

func isPhotoType(t fileid.Type) bool {
	switch t {
	case fileid.Thumbnail, fileid.ProfilePhoto, fileid.Photo:
		return true
	default:
		return false
	}
}

func uniqueTypeFor(t fileid.Type) int {
	switch t {
	case fileid.Thumbnail, fileid.ProfilePhoto, fileid.Photo:
		return uniqueTypePhoto
	case fileid.Secure, fileid.SecureRaw:
		return uniqueTypeSecure
	case fileid.Encrypted, fileid.EncryptedThumbnail:
		return uniqueTypeEncrypted
	case fileid.Temp:
		return uniqueTypeTemp
	default:
		return uniqueTypeDocument
	}
}

// rleEncode run-length-encodes runs of zero bytes, matching the encoding TDLib
// applies before base64url-encoding a file_unique_id. A copy of the unexported
// helper in github.com/gotd/td/fileid.
func rleEncode(s []byte) (r []byte) {
	var count byte

	for _, cur := range s {
		if cur == 0 {
			count++
			continue
		}

		if count > 0 {
			r = append(r, 0, count)
			count = 0
		}

		r = append(r, cur)
	}

	if count > 0 {
		r = append(r, 0, count)
	}

	return r
}
