package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kiasuo/bot/commands"
	"github.com/kiasuo/bot/users"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func GetUserID(interaction *discordgo.InteractionCreate) string {
	if interaction.Member == nil {
		return interaction.User.ID
	}

	return interaction.Member.User.ID
}

func main() {
	token, ok := os.LookupEnv("DISCORD")

	if !ok {
		panic("DISCORD not set")
	}

	bot, err := discordgo.New(token)

	if err != nil {
		panic(err)
	}

	bot.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		var command string
		var data []string

		switch interaction.Type {
		case discordgo.InteractionApplicationCommand:
			command = interaction.ApplicationCommandData().Name
		case discordgo.InteractionMessageComponent:
			data = strings.Split(interaction.MessageComponentData().CustomID, ":")
			command = data[0]
		default:
			return
		}

		userID := GetUserID(interaction)
		user := users.GetByDiscordID(userID)

		responder := commands.DiscordResponder{
			Interaction: *interaction.Interaction,
			Session:     session,
		}

		if user == nil {
			responder.Respond("Ты кто такой? Cъебал.")
			return
		}

		if user.State != users.Ready {
			responder.Respond("Пошел нахуй.")
			return
		}

		context := commands.Context{
			Command: command,
			User:    *user,
		}

		formatter := commands.DiscordFormatter{}

		if len(data) != 0 {
			commands.HandleCallback(context, &responder, &formatter, data)
			return
		}

		commands.HandleCommand(context, &responder, &formatter)
	})

	err = bot.Open()

	if err != nil {
		panic(err)
	}

	log.Println("Authorized on account", bot.State.User.Username)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	bot.Close()
}
