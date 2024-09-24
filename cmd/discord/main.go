package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kiasuo/bot/internal/commands"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
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
	token := helpers.GetEnv("DISCORD")
	bot, err := discordgo.New(token)

	if err != nil {
		panic(err)
	}

	bot.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		var (
			command string
			data    []string
		)

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

		err = responder.RespondWithDefer()

		if err != nil {
			log.Println(err)
			return
		}

		if user == nil {
			_ = responder.Write("Ты кто?").Respond()
			return
		}

		if !user.IsReady() {
			_ = responder.Write("Токен обнови.").Respond()
			return
		}

		context := commands.Context{
			Command: command,
			User:    *user,
		}

		formatter := helpers.DiscordFormatter{}

		if len(data) > 0 {
			commands.HandleCallback(context, &responder, &formatter, data)
			return
		}

		commands.HandleCommand(context, &responder, &formatter)
	})

	err = bot.Open()

	if err != nil {
		panic(err)
	}

	_, err = bot.ApplicationCommandBulkOverwrite(bot.State.User.ID, "", commands.ParseDiscordCommands())

	if err != nil {
		panic(err)
	}

	log.Println("Authorized on account", bot.State.User.Username)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	err = bot.Close()

	if err != nil {
		panic(err)
	}
}
