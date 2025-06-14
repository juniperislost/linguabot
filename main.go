package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Config struct {
	Token  string `json:"token"`
	Prefix string `json:"prefix"`
}

var (
	config Config
	botID  string
)

func main() {
	if err := loadConfig("config.json"); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Discord session: %v\n", err)
		os.Exit(1)
	}

	u, err := dg.User("@me")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get bot user: %v\n", err)
		os.Exit(1)
	}
	botID = u.ID

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open Discord session: %v\n", err)
		os.Exit(1)
	}
	defer dg.Close()

	fmt.Println("LinguaBot is now running. Press CTRL-C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}

func loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &config)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botID {
		return
	}

	if !strings.HasPrefix(m.Content, config.Prefix) {
		return
	}

	content := strings.TrimPrefix(m.Content, config.Prefix)
	args := strings.Fields(content)
	if len(args) == 0 {
		return
	}

	cmdName := strings.ToLower(args[0])
	cmdArgs := args[1:]

	if handler, exists := CommandRegistry[cmdName]; exists {
		handler(s, m, cmdArgs)
	} else {
		s.ChannelMessageSend(m.ChannelID, "Unknown command. Try `l!help`.")
	}
}
