package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/kiasuo/bot/internal/commands"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
	"github.com/kiasuo/bot/internal/version"
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

		responder := commands.DiscordResponder{
			Interaction: *interaction.Interaction,
			Session:     session,
		}

		switch interaction.Type {
		case discordgo.InteractionApplicationCommand:
			command = interaction.ApplicationCommandData().Name
		case discordgo.InteractionMessageComponent:
			data = strings.Split(interaction.MessageComponentData().CustomID, ":")

			if len(data) < 2 {
				return
			}

			if data[0] != version.Version {
				_ = responder.Write("Версию меню устарела.").Respond()
				return
			}

			command = data[1]
		default:
			return
		}

		userID := GetUserID(interaction)
		id, state := users.GetPartialByDiscordID(userID)

		if err = responder.RespondWithDefer(); err != nil {
			log.Println(err)
			return
		}

		if state == users.Unknown {
			_ = responder.Write("Ты кто?").Respond()
			return
		}

		if state != users.Ready {
			_ = responder.Write("Токен обнови.").Respond()
			return
		}

		context := commands.Context{
			Command: command,
			User:    *users.GetByID(id),
		}

		formatter := helpers.DiscordFormatter{}

		if len(data) > 0 {
			commands.HandleCallback(context, &responder, &formatter, data[2:])
			return
		}

		commands.HandleCommand(context, &responder, &formatter)
	})

	if err = bot.Open(); err != nil {
		panic(err)
	}

	_, err = bot.ApplicationCommandBulkOverwrite(
		bot.State.User.ID,
		"",
		commands.ParseDiscordCommands(),
	)

	if err != nil {
		panic(err)
	}

	log.Println("Authorized on account", bot.State.User.Username)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if err = bot.Close(); err != nil {
		panic(err)
	}
}
