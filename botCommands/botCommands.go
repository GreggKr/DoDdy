package botcommands

import (
	"github.com/Devs-On-Discord/DoDdy/commands"
	"github.com/Devs-On-Discord/DoDdy/guilds"
	"github.com/bwmarrin/discordgo"
)

// Instance is a globally accessible BotCommands object
var Instance BotCommands

// Init constructs the global BotCommands object
func Init(guilds *guilds.Guilds, session *discordgo.Session) {
	Instance = BotCommands{}
	Instance.Init(guilds, session)
}

// BotCommands is an object that encapsulates both Commands and a result handler
type BotCommands struct {
	commands                    *commands.Commands
	discordCommandResultHandler *commands.DiscordCommandResultHandler
	guilds                      *guilds.Guilds
}

// Init constructs the BotCommands object
func (b *BotCommands) Init(guilds *guilds.Guilds, session *discordgo.Session) {
	b.guilds = guilds
	b.commands = &commands.Commands{}
	b.commands.Init(session)
	b.discordCommandResultHandler = &commands.DiscordCommandResultHandler{}
	b.discordCommandResultHandler.Init(b.commands, session)
	b.RegisterCommands()
	session.AddHandler(b.messageHandler)
}

// RegisterCommands registers commands with the Commands object
func (b *BotCommands) RegisterCommands() {
	guildAdminCommands := guildAdminCommands{guilds: b.guilds}
	helpCommands := helpCommands{&b.commands.Commands}
	b.commands.Register(commands.Command{Name: "prefix", Handler: guildAdminCommands.setPrefix, Description: "Changes / Displays the prefix."})
	b.commands.Register(commands.Command{Name: "setAnnouncementsChannel", Handler: guildAdminCommands.setAnnouncementsChannel, Description: "Redefines this node's announcement channel."})
	b.commands.Register(commands.Command{Name: "announce announcement", Handler: guildAdminCommands.postAnnouncement, Description: "Post an announcement in this node."})
	b.commands.Register(commands.Command{Name: "clearAnnouncements", Handler: guildAdminCommands.clearAnnouncements, Description: "Empties this node's announcement channel."})
	b.commands.Register(commands.Command{Name: "postLastMessageAsAnnouncement", Handler: guildAdminCommands.postLastMessageAsAnnouncement, Description: "Repost the last message sent in this channel as an announcement"})
	b.commands.Register(commands.Command{Name: "setVotesChannel", Handler: guildAdminCommands.setVotesChannel, Description: "Redefines this node's voting channel."})
	b.commands.Register(commands.Command{Name: "survey vote", Handler: guildAdminCommands.postVote, Description: "Starts a DoD-wide survey."})
	b.commands.Register(commands.Command{Name: "setup", Handler: guildAdminCommands.setup, Description: "Modifies basic configuration settings"})
	b.commands.Register(commands.Command{Name: "help", Handler: helpCommands.helpCommand, Description: "Begins a vote in this node's voting channel."})
}

// Parse is the input sink for commands
func (b *BotCommands) Parse(message *discordgo.MessageCreate) {
	b.commands.Parse(message)
}

func (b *BotCommands) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	botID := s.State.User.ID
	if m.Author.ID == botID {
		return
	}
	if len(m.Content) == 0 {
		return
	}
	valid := false
	//TODO: support multiple mentions and append all mentions at the bot reply
	for _, mention := range m.Mentions {
		if mention.ID == botID {
			m.Content = m.Content[len(mention.ID)+3:] //<@{botID}>
			valid = true
			break
		}
	}
	if !valid {
		input := m.Content
		if len(input) > 1 {
			guild, err := b.guilds.Guild(m.GuildID)
			if err != nil {
				return
			}
			if guild.Prefix == input[:1] {
				m.Content = m.Content[1:]
			} else {
				return
			}
		} else {
			return
		}
	}
	b.Parse(m)
}
