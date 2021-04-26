package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/blyndusk/cofy/app/commands"
	"github.com/blyndusk/cofy/app/core"
	"github.com/blyndusk/cofy/app/helpers"
	"github.com/blyndusk/cofy/app/middlewares"
	"github.com/blyndusk/cofy/app/services"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func init() {
	// get Discord token
	Token := helpers.EnvVar("BOT_TOKEN")

	// create a variable named "t", as the bot token
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	pingApi()
	setupSession()
}

func messageCreationHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// if the message author is the bot itself, return to avoid infinite loops
	if m.Author.ID == s.State.User.ID {
		return
	}

	// if the message is not a command
	if !strings.HasPrefix(m.Content, core.Prefix) {
		// systematically get user each time someone create a message
		user := middlewares.GetUser(s, m.Author.ID)
		// if the user is not found, create a new user in db
		services.UserNotFoundHandler(s, m, user)
		// then update gains
		services.GainsHandler(s, m, user)

	} else {
		commands.ProfileCommandHandler(s, m)
		// commands.Info(s, m)
		// commands.Dev(s, m)
	}

}

func MessageReactionHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	matched, err := regexp.MatchString(`coffee|coffe|cofee|cofe|cafe|café`, m.Content)
	fmt.Println(matched, err)

	if matched == true {
		if m.Author.ID == s.State.User.ID {
			return
		}
		s.MessageReactionAdd(m.ChannelID, m.ID, "☕")

	}

}

func setupSession() {
	Token := os.Getenv("BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreationHandler)
	dg.AddHandler(MessageReactionHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func pingApi() {
	_, err := http.Get(helpers.EnvVar("API_URL"))
	if err != nil {
		log.Fatal(err)
	}
	log.Info("API available !")
}
