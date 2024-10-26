package commands

import "testing"

func TestIsSystemCommand(t *testing.T) {
	if !IsSystemCommand(StartCommandName) {
		t.Errorf("IsSystemCommand(%s) = false; want true\n", StartCommandName)
	}

	if IsSystemCommand("settings") {
		t.Errorf("IsSystemCommand(%s) = true; want false\n", "settings")
	}
}

func TestParseTelegramCommands(t *testing.T) {
	commands := ParseTelegramCommands()

	if len(commands.Commands) != len(publicCommands) {
		t.Errorf("ParseTelegramCommands() = %d; want %d\n", len(commands.Commands), len(publicCommands))
	}

	for i, command := range commands.Commands {
		if command.Command != publicCommands[i].Name {
			t.Errorf("ParseTelegramCommands()[%d].Command = %s; want %s\n", i, command.Command, publicCommands[i].Name)
		}

		if command.Description != publicCommands[i].Description {
			t.Errorf(
				"ParseTelegramCommands()[%d].Description = %s; want %s\n",
				i,
				command.Description,
				publicCommands[i].Description,
			)
		}
	}
}

func TestParseDiscordCommands(t *testing.T) {
	commands := ParseDiscordCommands()

	// because we have commands that is TelegramOnly
	// we need to filter them out
	publicDiscordCommands := make([]commandConfig, 0)

	for _, c := range publicCommands {
		if !IsSystemCommand(c.Name) {
			publicDiscordCommands = append(publicDiscordCommands, c)
		}
	}

	for i, command := range commands {
		if command.Name != publicDiscordCommands[i].Name {
			t.Errorf("ParseDiscordCommands()[%d].Name = %s; want %s\n", i, command.Name, publicDiscordCommands[i].Name)
		}

		if command.Description != publicDiscordCommands[i].Description {
			t.Errorf(
				"ParseDiscordCommands()[%d].Description = %s; want %s\n",
				i,
				command.Description,
				publicDiscordCommands[i].Description,
			)
		}
	}
}
