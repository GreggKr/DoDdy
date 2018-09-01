package commands

import (
	"strings"

	"github.com/anmitsu/go-shlex"
	"github.com/bwmarrin/discordgo"
)

// Commands is an object containing all commands that can be called, it converts discordgo commands to a single thread each
// ResultMessages is the return channel for successful commands
type Commands struct {
	commands         map[string]Command
	ResultMessages   chan CommandResultMessage
	incomingMessages chan *discordgo.MessageCreate
}

// Init constructs the Commands object
func (c *Commands) Init() {
	c.commands = make(map[string]Command)
	c.ResultMessages = make(chan CommandResultMessage)
	c.incomingMessages = make(chan *discordgo.MessageCreate)
	go func() {
		for {
			incomingMessage := <-c.incomingMessages
			c.parse(incomingMessage)
		}
	}()
}

// Register associates a Command name to a Handler
func (c *Commands) Register(command Command) {
	commandNameSplit := strings.Split(command.Name, " ")
	if len(commandNameSplit) < 1 {
		return
	}
	name := commandNameSplit[0]
	c.commands[name] = command
}

func (c *Commands) parse(commandMessage *discordgo.MessageCreate) {
	commandParsed, err := shlex.Split(commandMessage.Content, true)
	if err != nil {
		c.ResultMessages <- &CommandError{
			CommandMessage: commandMessage,
			Message:        "Error happened " + err.Error(),
			Color:          0xb30000,
		}
	}
	commandCount := len(commandParsed)
	if commandCount < 1 {
		c.ResultMessages <- &CommandError{
			CommandMessage: commandMessage,
			Message:        "Invalid Command",
			Color:          0xb30000,
		}
	}
	commandName := commandParsed[0]
	if command, exists := c.commands[commandName]; exists {
		if commandCount < 2 {
			resultMessage := command.Handler(commandMessage, nil)
			resultMessage.setCommandMessage(commandMessage)
			c.ResultMessages <- resultMessage
		} else {
			resultMessage := command.Handler(commandMessage, commandParsed[1:])
			resultMessage.setCommandMessage(commandMessage)
			c.ResultMessages <- resultMessage
		}
	} else {
		c.ResultMessages <- &CommandError{
			CommandMessage: commandMessage,
			Message:        "Command doesn't exists",
			Color:          0xb30000,
		}
	}
}

// Parse is the input sink for commands
func (c *Commands) Parse(commandMessage *discordgo.MessageCreate) {
	c.incomingMessages <- commandMessage
}
