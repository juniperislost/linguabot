package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type CommandFunc func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)

var CommandRegistry map[string]CommandFunc

const embedColor = 0xC8A2C8

func init() {
	CommandRegistry = map[string]CommandFunc{
		"botinfo":    botInfoCommand,
		"userinfo":   userInfoCommand,
		"serverinfo": serverInfoCommand,
		"ping":       pingCommand,
		"help":       helpCommand,
	}
}

func getCommandCount() int {
	return len(CommandRegistry)
}

func botInfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	botUser, err := s.User("@me")
	var avatarURL string
	if err == nil {
		avatarURL = botUser.AvatarURL("")
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Bot Information",
		Description: "Official bot for the Lingua Commons Discord server.",
		Color:       embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Version",
				Value: "v1.0",
			},
			{
				Name:  "Author",
				Value: "<@1335732044143001712>",
			},
			{
				Name:  "Library",
				Value: "discordgo (Go)",
			},
			{
				Name:   "Commands",
				Value:  fmt.Sprintf("%d commands loaded", getCommandCount()),
				Inline: true,
			},
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func userInfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	user := m.Author
	if len(args) > 0 {
		u, err := s.User(args[0])
		if err == nil {
			user = u
		}
	}

	avatarURL := user.AvatarURL("")

	embed := &discordgo.MessageEmbed{
		Title: "User Information",
		Color: embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Username",
				Value:  fmt.Sprintf("%s#%s", user.Username, user.Discriminator),
				Inline: true,
			},
			{
				Name:   "User ID",
				Value:  user.ID,
				Inline: true,
			},
			{
				Name:   "Bot?",
				Value:  fmt.Sprintf("%t", user.Bot),
				Inline: true,
			},
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func serverInfoCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		guild, err = s.Guild(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error fetching server info.")
			return
		}
	}

	iconURL := guild.IconURL("")

	embed := &discordgo.MessageEmbed{
		Title: "Server Information",
		Color: embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: iconURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Server Name",
				Value:  guild.Name,
				Inline: true,
			},
			{
				Name:   "Server ID",
				Value:  guild.ID,
				Inline: true,
			},
			{
				Name:   "Member Count",
				Value:  fmt.Sprintf("%d", guild.MemberCount),
				Inline: true,
			},
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func pingCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	start := time.Now()
	restLatency := time.Since(start).Milliseconds()
	gatewayLatency := s.HeartbeatLatency().Milliseconds()

	botUser, err := s.User("@me")
	var avatarURL string
	if err == nil {
		avatarURL = botUser.AvatarURL("")
	}

	embed := &discordgo.MessageEmbed{
		Title: "Pong!",
		Color: embedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: avatarURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Gateway Latency",
				Value:  fmt.Sprintf("`%dms`", gatewayLatency),
				Inline: true,
			},
			{
				Name:   "REST Latency",
				Value:  fmt.Sprintf("`%dms`", restLatency),
				Inline: true,
			},
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func helpCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	prefix := "!"
	embed := &discordgo.MessageEmbed{
		Title:       "Help — LinguaBot Commands",
		Description: "Here's a list of available commands.",
		Color:       embedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Information",
				Value: fmt.Sprintf("`%sbotinfo`, `%suserinfo`, `%sserverinfo`, `%sping`, `%shelp`", prefix, prefix, prefix, prefix, prefix),
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "LinguaBot",
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
