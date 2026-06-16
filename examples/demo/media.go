package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"time"

	"github.com/gotd/botapi"
)

const samplePhotoURL = "https://upload.wikimedia.org/wikipedia/commons/3/3a/Cat03.jpg"

// registerMedia covers media in both directions: commands that send media a few
// different ways, and handlers that react to incoming media.
func registerMedia(bot *botapi.Bot) {
	registerOutgoingMedia(bot)
	registerIncomingMedia(bot)
}

// registerOutgoingMedia sends media built from a URL, from in-memory bytes, and
// as an album, and shows a chat action while "working".
func registerOutgoingMedia(bot *botapi.Bot) {
	// FileURL: Telegram fetches the file itself from the given URL.
	bot.OnCommand("photo", "Send a photo by URL", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		_, err := c.Bot.SendPhoto(c, chat, botapi.FileURL(samplePhotoURL), "A cat 🐈 sent by URL")

		return err
	})

	// FileFromBytes: upload a document built entirely in memory.
	bot.OnCommand("document", "Send a generated text document", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		body := []byte("Hello from gotd/botapi!\nThis file was generated in memory.\n")
		_, err := c.Bot.SendDocument(c, chat, botapi.FileFromBytes("hello.txt", body), "A generated file")

		return err
	})

	// A media group (album) of two photos sent as one grouped message. Media-group
	// items must be uploaded files — file_id and URL are not accepted — so the
	// photos are generated in memory and uploaded via FileFromBytes.
	bot.OnCommand("album", "Send a media-group album", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		media := []botapi.InputMedia{
			&botapi.InputMediaPhoto{
				Type:    botapi.InputMediaPhotoType,
				Media:   botapi.FileFromBytes("red.png", solidPNG(color.RGBA{R: 0xD0, G: 0x21, B: 0x21, A: 0xFF})),
				Caption: "First",
			},
			&botapi.InputMediaPhoto{
				Type:    botapi.InputMediaPhotoType,
				Media:   botapi.FileFromBytes("blue.png", solidPNG(color.RGBA{R: 0x15, G: 0x65, B: 0xC0, A: 0xFF})),
				Caption: "Second",
			},
		}
		_, err := c.Bot.SendMediaGroup(c, chat, media)

		return err
	})

	// A chat action ("typing…") followed by a reply.
	bot.OnCommand("typing", "Show a chat action, then reply", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		if err := c.Bot.SendChatAction(c, chat, botapi.ChatActionTyping); err != nil {
			return err
		}

		time.Sleep(2 * time.Second)

		_, err := c.Reply("Done typing!")

		return err
	})
}

// registerIncomingMedia reacts to media the user sends, each gated by a
// content-type predicate. EffectiveMessage handles new, edited and channel posts
// uniformly.
func registerIncomingMedia(bot *botapi.Bot) {
	// Echo a received photo back by its file_id (no re-upload) and report its size.
	bot.OnMessage(func(c *botapi.Context) error {
		photos := c.Message().Photo
		largest := photos[len(photos)-1]
		chat, _ := c.Chat()
		_, err := c.Bot.SendPhoto(c, chat, botapi.FileID(largest.FileID),
			fmt.Sprintf("Got your photo: %d×%d", largest.Width, largest.Height))

		return err
	}, hasPhoto)

	// Resolve a document's downloadable file via GetFile.
	bot.OnMessage(func(c *botapi.Context) error {
		doc := c.Message().Document

		file, err := c.Bot.GetFile(c, doc.FileID)
		if err != nil {
			return err
		}

		_, err = c.Reply(fmt.Sprintf("Got document %q (%d bytes)\nfile_unique_id=%s",
			doc.FileName, doc.FileSize, file.FileUniqueID))

		return err
	}, hasDocument)

	bot.OnMessage(func(c *botapi.Context) error {
		_, err := c.Reply("Nice sticker " + c.Message().Sticker.Emoji)
		return err
	}, hasSticker)

	bot.OnMessage(func(c *botapi.Context) error {
		loc := c.Message().Location
		_, err := c.Reply(fmt.Sprintf("📍 You are at %.4f, %.4f", loc.Latitude, loc.Longitude))

		return err
	}, hasLocation)

	bot.OnMessage(func(c *botapi.Context) error {
		ct := c.Message().Contact
		_, err := c.Reply("📇 Contact: " + ct.FirstName + " " + ct.PhoneNumber)

		return err
	}, hasContact)
}

// solidPNG renders a 256×256 solid-color PNG in memory, so the album demo has
// real uploadable image bytes without shipping asset files.
func solidPNG(c color.Color) []byte {
	const size = 256

	img := image.NewRGBA(image.Rect(0, 0, size, size))
	for y := range size {
		for x := range size {
			img.Set(x, y, c)
		}
	}

	var buf bytes.Buffer

	if err := png.Encode(&buf, img); err != nil {
		panic(err) // encoding an in-memory RGBA image cannot fail
	}

	return buf.Bytes()
}
