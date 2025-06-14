package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	config struct {
		Token         string `json:"token"`
		Prefix        string `json:"prefix"`
		QOTDChannelID string `json:"qotd_channel_id"`
		QOTDRoleID    string `json:"qotd_role_id"`
	}
	qotdQuotes []string
)

func main() {
	loadConfig("config.json")
	loadQuotes("quotes.json")

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Fatalf("Failed to create Discord session: %v", err)
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalf("Failed to open Discord connection: %v", err)
	}
	defer dg.Close()

	log.Println("LinguaBot is now running!")

	go startQOTDAutoPost(dg)

	select {}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	if len(m.Content) < len(config.Prefix) || m.Content[:len(config.Prefix)] != config.Prefix {
		return
	}

	args := strings.Fields(m.Content[len(config.Prefix):])
	if len(args) == 0 {
		return
	}

	cmd := strings.ToLower(args[0])
	args = args[1:]

	if handler, ok := CommandRegistry[cmd]; ok {
		handler(s, m, args)
	}
}

func loadConfig(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
}

func loadQuotes(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read quotes: %v", err)
	}
	err = json.Unmarshal(data, &qotdQuotes)
	if err != nil {
		log.Fatalf("Failed to parse quotes: %v", err)
	}
}

func startQOTDAutoPost(s *discordgo.Session) {
	postHour := 9
	postMinute := 0

	for {
		now := time.Now().UTC()
		nextPost := time.Date(now.Year(), now.Month(), now.Day(), postHour, postMinute, 0, 0, time.UTC)
		if !nextPost.After(now) {
			nextPost = nextPost.Add(24 * time.Hour)
		}
		time.Sleep(nextPost.Sub(now))

		if len(qotdQuotes) == 0 || config.QOTDChannelID == "" {
			continue
		}

		question := qotdQuotes[rand.Intn(len(qotdQuotes))]

		content := ""
		if config.QOTDRoleID != "" {
			content = "<@&" + config.QOTDRoleID + ">"
		}

		embed := &discordgo.MessageEmbed{
			Title:       "Question of the Day",
			Description: question,
			Color:       0xC8A2C8,
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
		}

		_, err := s.ChannelMessageSendComplex(config.QOTDChannelID, &discordgo.MessageSend{
			Content: content,
			Embed:   embed,
		})
		if err != nil {
			log.Printf("Failed to send QOTD: %v", err)
		}

		time.Sleep(24 * time.Hour)
	}
}
