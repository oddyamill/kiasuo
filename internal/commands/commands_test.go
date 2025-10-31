package commands

import "testing"

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
