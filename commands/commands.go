package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/foxtrot/scuzzy/models"
	"github.com/foxtrot/scuzzy/permissions"
)

type ScuzzyHandler func(session *discordgo.Session, m *discordgo.MessageCreate) error

type ScuzzyCommand struct {
	Index       int
	Name        string
	Description string
	AdminOnly   bool
	Handler     ScuzzyHandler
}

type Commands struct {
	Token                 string
	Permissions           *permissions.Permissions
	Config                *models.Configuration
	ScuzzyCommands        map[string]ScuzzyCommand
	ScuzzyCommandsByIndex map[int]ScuzzyCommand
}
