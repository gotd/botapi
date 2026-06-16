package main

import (
	"fmt"

	glog "github.com/gotd/log"

	"github.com/gotd/botapi"
)

// registerAdmin builds a Group whose handlers only run in group/supergroup chats
// and share an extra middleware. Groups keep related, similarly-gated handlers
// together instead of repeating the same predicate on each registration.
func registerAdmin(bot *botapi.Bot) {
	// The group predicate narrows every handler below to group chats; Use adds a
	// middleware that runs only for this group.
	groups := bot.Group(botapi.Or(
		botapi.ChatTypeIs(botapi.ChatTypeGroup),
		botapi.ChatTypeIs(botapi.ChatTypeSupergroup),
	)).Use(announce())

	// Chat info: title, members count and the caller's membership status.
	groups.OnCommand("info", "Show info about this group", func(c *botapi.Context) error {
		chat, _ := c.Chat()

		info, err := c.Bot.GetChat(c, chat)
		if err != nil {
			return err
		}

		count, err := c.Bot.GetChatMemberCount(c, chat)
		if err != nil {
			return err
		}

		status := "unknown"

		if u := c.Sender(); u != nil {
			if member, err := c.Bot.GetChatMember(c, chat, u.ID); err == nil {
				status = chatMemberStatus(member)
			}
		}

		_, err = c.Reply(fmt.Sprintf("Group %q\nMembers: %d\nYour status: %s",
			info.Title, count, status))

		return err
	})

	// React to the command message with an emoji, then pin it.
	groups.OnCommand("pin", "React to and pin your command", func(c *botapi.Context) error {
		chat, _ := c.Chat()
		id := c.Message().MessageID

		if err := c.Bot.SetMessageReaction(c, chat, id, []botapi.ReactionType{botapi.Emoji("👍")}); err != nil {
			return err
		}

		if err := c.Bot.PinChatMessage(c, chat, id, botapi.Silent()); err != nil {
			return err
		}

		_, err := c.Reply("📌 Pinned (and reacted).")

		return err
	})

	// Raw() is the escape hatch to the underlying *tg.Client for anything the
	// typed surface doesn't cover. Here we just confirm it's reachable.
	groups.OnCommand("raw", "Demonstrate the raw MTProto client escape hatch", func(c *botapi.Context) error {
		_ = c.Bot.Raw() // *tg.Client — full MTProto API available here.

		_, err := c.Reply("The raw *tg.Client is reachable via Bot.Raw() for anything not yet typed.")

		return err
	})
}

// announce is a group-scoped middleware: it logs every group-chat update it
// handles before delegating. Unlike the global metrics middleware it is attached
// only to this Group, so it runs for group commands only.
func announce() botapi.Middleware {
	return func(next botapi.Handler) botapi.Handler {
		return func(c *botapi.Context) error {
			glog.For(c.Bot.Logger()).Debug(c, "group command", glog.String("text", c.Message().Text))

			return next(c)
		}
	}
}

// chatMemberStatus extracts the status string from a ChatMember union value.
func chatMemberStatus(m botapi.ChatMember) string {
	switch v := m.(type) {
	case *botapi.ChatMemberOwner:
		return string(v.Status)
	case *botapi.ChatMemberAdministrator:
		return string(v.Status)
	case *botapi.ChatMemberMember:
		return string(v.Status)
	case *botapi.ChatMemberRestricted:
		return string(v.Status)
	case *botapi.ChatMemberLeft:
		return string(v.Status)
	case *botapi.ChatMemberBanned:
		return string(v.Status)
	default:
		return "unknown"
	}
}
