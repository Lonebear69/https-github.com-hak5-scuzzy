package features

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/foxtrot/scuzzy/models"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func (f *Features) handleSetConfig(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	}

	configArgs := strings.Split(m.Content, " ")

	if len(configArgs) != 3 {
		return errors.New("Invalid arguments supplied. Usage: " + f.Config.CommandKey + "setconfig <key> <value>")
	}

	configKey := configArgs[1]
	configVal := configArgs[2]

	rt := reflect.TypeOf(f.Config)
	for i := 0; i < rt.NumField(); i++ {
		x := rt.Field(i)
		tagVal := strings.Split(x.Tag.Get("json"), ",")[0]
		tagName := x.Name

		if tagVal == configKey {
			prop := reflect.ValueOf(&f.Config).Elem().FieldByName(tagName)

			switch prop.Interface().(type) {
			case string:
				prop.SetString(configVal)
				break
			case int:
				intVal, err := strconv.ParseInt(configVal, 10, 64)
				if err != nil {
					return err
				}
				prop.SetInt(intVal)
				break
			case float64:
				floatVal, err := strconv.ParseFloat(configVal, 64)
				if err != nil {
					return err
				}
				prop.SetFloat(floatVal)
				break
			case bool:
				boolVal, err := strconv.ParseBool(configVal)
				if err != nil {
					return err
				}
				prop.SetBool(boolVal)
				break
			default:
				return errors.New("Unsupported key value type")
			}

			msgE := f.CreateDefinedEmbed("Set Configuration", "Successfully set property '"+configKey+"'!", "success", m.Author)
			_, err := s.ChannelMessageSendEmbed(m.ChannelID, msgE)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return errors.New("Unknown key specified")
}

func (f *Features) handleGetConfig(s *discordgo.Session, m *discordgo.MessageCreate) error {
	//TODO: Handle printing of slices (check the Type, loop accordingly)

	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	}

	configArgs := strings.Split(m.Content, " ")
	configKey := "all"
	if len(configArgs) == 2 {
		configKey = configArgs[1]
	}

	msg := ""

	rt := reflect.TypeOf(f.Config)
	for i := 0; i < rt.NumField(); i++ {
		x := rt.Field(i)
		tagVal := strings.Split(x.Tag.Get("json"), ",")[0]
		tagName := x.Name
		prop := reflect.ValueOf(&f.Config).Elem().FieldByName(tagName)

		if configKey == "all" {
			switch prop.Interface().(type) {
			case string:
				if len(prop.String()) > 256 {
					// Truncate large values.
					msg += "`" + tagName + "` - " + "Truncated...\n"
				} else {
					msg += "`" + tagName + "` - `" + prop.String() + "`\n"
				}
				break
			default:
				// Ignore non strings for now...
				msg += "`" + tagName + "` - Skipped Value\n"
				continue
			}
		} else {
			if tagVal == configKey {
				switch prop.Interface().(type) {
				case string:
					msg += "`" + tagName + "` - `" + prop.String() + "`\n"
				default:
					// Ignore non strings for now...
					msg += "`" + tagName + "` - Skipped Value\n"
				}

				eMsg := f.CreateDefinedEmbed("Get Configuration", msg, "success", m.Author)
				_, err := s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
				if err != nil {
					return err
				}

				return nil
			}
		}
	}

	if msg == "" {
		return errors.New("Unknown key specified")
	}

	eMsg := f.CreateDefinedEmbed("Get Configuration", msg, "success", m.Author)
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleReloadConfig(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use this command.")
	}

	fBuf, err := ioutil.ReadFile(f.Config.ConfigPath)
	if err != nil {
		return err
	}

	conf := &models.Configuration{}

	err = json.Unmarshal(fBuf, &conf)
	if err != nil {
		return err
	}

	f.Config = conf
	f.Permissions.Config = conf

	eMsg := f.CreateDefinedEmbed("Reload Configuration", "Successfully reloaded configuration from disk", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleSaveConfig(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use this command.")
	}

	j, err := json.Marshal(f.Config)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(f.Config.ConfigPath, j, os.ModePerm)
	if err != nil {
		return err
	}

	eMsg := f.CreateDefinedEmbed("Save Configuration", "Saved runtime configuration successfully", "success", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, eMsg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleCat(s *discordgo.Session, m *discordgo.MessageCreate) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://giphy.com/gifs/cat-cute-no-rCxogJBzaeZuU")
	if err != nil {
		return err
	}

	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handlePing(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var r *discordgo.Message
	var err error

	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	} else {
		msg := f.CreateDefinedEmbed("Ping", "Pong", "success", m.Author)
		r, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
		if err != nil {
			return err
		}
	}

	time.Sleep(5 * time.Second)

	err = s.ChannelMessageDelete(m.ChannelID, r.ID)
	if err != nil {
		return err
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	desc := "**Source**:   https://github.com/foxtrot/scuzzy\n"
	desc += "**Language**: Go\n"
	desc += "**Commands**: See `" + f.Config.CommandKey + "help`\n\n\n"

	gm, err := s.GuildMember(f.Config.GuildID, s.State.User.ID)
	if err != nil {
		return err
	}

	d := models.CustomEmbed{
		Title:          "Scuzzy Information",
		Desc:           desc,
		ImageURL:       "",
		ImageH:         100,
		ImageW:         100,
		Color:          0xFFA500,
		URL:            "",
		Type:           "",
		Timestamp:      "",
		FooterText:     "Made with  ❤  by Foxtrot",
		FooterImageURL: "https://cdn.discordapp.com/avatars/514163441548656641/a_ac5e022e77e62e7793711ebde8cdf4a1.gif",
		ThumbnailURL:   gm.User.AvatarURL(""),
		ThumbnailH:     150,
		ThumbnailW:     150,
		ProviderURL:    "",
		ProviderText:   "",
		AuthorText:     "",
		AuthorURL:      "",
		AuthorImageURL: "",
	}

	msg := f.CreateCustomEmbed(&d)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleHelp(s *discordgo.Session, m *discordgo.MessageCreate) error {
	desc := "**Available Commands**\n"
	desc += "__Misc__\n"
	desc += "`help` - This help dialog\n"
	desc += "`info` - Display Scuzzy info\n"
	desc += "`md` - Display Discord markdown information\n"
	desc += "`userinfo` - Display information about a user\n"
	desc += "`serverinfo` - Display information about the server\n"

	desc += "\n__User Settings__\n"
	desc += "`colors` - Available color roles\n"
	desc += "`color` - Set an available color role\n"
	desc += "`listroles` - List available roles\n"
	desc += "`joinrole` - Join an available role\n"
	desc += "`leaverole` - Leave a joined role\n"

	desc += "\n__Conversion Helpers__\n"
	desc += "`ctof` - Convert Celsius to Farenheit\n"
	desc += "`ftoc` - Convert Farenheit to Celsius\n"
	desc += "`metofe` - Convert Meters to Feet\n"
	desc += "`fetome` - Convert Feet to Meters\n"
	desc += "`cmtoin` - Convert Centimeters to Inches\n"
	desc += "`intocm` - Convert Inches to Centimeters\n"

	if f.Permissions.CheckAdminRole(m.Member) {
		desc += "\n"
		desc += "**Admin Commands**\n"
		desc += "`ping` - Ping the bot\n"
		desc += "`rules` - Display the server rules\n"
		desc += "`status` - Set the bot status\n"
		desc += "`purge` - Purge channel messages\n"
		desc += "`kick` - Kick a specified user\n"
		desc += "`ban` - Ban a specified user\n"
		desc += "`ignore` - Ignore a specified user\n"
		desc += "`unignore` - Unignore a specified user\n"
		desc += "`setconfig` - Manage the runtime configuration\n"
		desc += "`getconfig` - View the runtime configuration\n"
		desc += "`reloadconfig` - Reload configuration from disk\n"
		desc += "`saveconfig` - Save the runtime configuration to disk\n"
		desc += "`addrole` - Add a user joinable role\n"
	}

	desc += "\n\nAll commands are prefixed with `" + f.Config.CommandKey + "`\n"

	msg := f.CreateDefinedEmbed("Help", desc, "", m.Author)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleRules(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if !f.Permissions.CheckAdminRole(m.Member) {
		return errors.New("You do not have permissions to use that command.")
	}

	msg := f.Config.RulesText
	embedTitle := "Rules (" + f.Config.GuildName + ")"
	embed := f.CreateDefinedEmbed(embedTitle, msg, "success", m.Author)

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleMarkdownInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	cleanup := true
	args := strings.Split(m.Content, " ")

	if len(args) == 2 {
		if args[1] == "stay" && f.Permissions.CheckAdminRole(m.Member) {
			cleanup = false
		}
	}

	desc := "*Italic* text goes between `*single asterisks*`\n"
	desc += "**Bold** text goes between `**double asterisks**`\n"
	desc += "***Bold and Italic*** text goes between `***triple asterisks***`\n"
	desc += "__Underlined__ text goes between `__double underscore__`\n"
	desc += "~~Strikethrough~~ text goes between `~~double tilde~~`\n"
	desc += "||Spoilers|| go between `|| double pipe ||`\n\n"
	desc += "You can combine the above styles.\n\n"
	desc += "Inline Code Blocks start and end with a single ``​`​``\n"
	desc += "Multi line Code Blocks start and end with ``​```​``\n"
	desc += "Multi line Code Blocks can also specify a language with ``​```​language`` at the start\n\n"
	desc += "Single line quotes start with `>`\n"
	desc += "Multi line quotes start with `>>>`\n"

	msg := f.CreateDefinedEmbed("Discord Markdown", desc, "", m.Author)
	r, err := s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	if cleanup {
		time.Sleep(15 * time.Second)

		err = s.ChannelMessageDelete(m.ChannelID, r.ID)
		if err != nil {
			return err
		}
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Features) handleCtoF(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a temperature")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	cels := (inF * 9.0 / 5.0) + 32.0
	celsF := float64(cels)

	msg := fmt.Sprintf("`%.1f°c` is `%.1f°f`", inF, celsF)

	e := f.CreateDefinedEmbed("Celsius to Farenheit", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleFtoC(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a temperature")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	faren := (inF - 32) * 5 / 9
	farenF := float64(faren)

	msg := fmt.Sprintf("`%.1f°f` is `%.1f°c`", inF, farenF)

	e := f.CreateDefinedEmbed("Farenheit to Celsius", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleMetersToFeet(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a distance")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	meters := inF * 3.28
	metersF := float64(meters)

	msg := fmt.Sprintf("`%.1fm` is `%.1fft`", inF, metersF)

	e := f.CreateDefinedEmbed("Meters to Feet", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleFeetToMeters(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a distance")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	feet := inF / 3.28
	feetF := float64(feet)

	msg := fmt.Sprintf("`%.1fft` is `%.1fm`", inF, feetF)

	e := f.CreateDefinedEmbed("Feet to Meters", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleCentimeterToInch(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a distance")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	inch := inF / 2.54
	inchF := float64(inch)

	msg := fmt.Sprintf("`%.1fcm` is `%.1fin`", inF, inchF)

	e := f.CreateDefinedEmbed("Centimeter To Inch", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleInchToCentimeter(s *discordgo.Session, m *discordgo.MessageCreate) error {
	inS := strings.Split(m.Content, " ")

	if len(inS) < 2 {
		return errors.New("You did not specify a distance")
	}
	in := inS[1]

	inF, err := strconv.ParseFloat(in, 2)
	if err != nil {
		return errors.New("You did not specify a valid number")
	}

	cm := inF * 2.54
	cmF := float64(cm)

	msg := fmt.Sprintf("`%.1fin` is `%.1fcm`", inF, cmF)

	e := f.CreateDefinedEmbed("Inch to Centimeter", msg, "", m.Author)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, e)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleUserInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	var (
		mHandle   *discordgo.Member
		requester *discordgo.Member
		err       error
	)

	userSplit := strings.Split(m.Content, " ")

	if len(userSplit) < 2 {
		mHandle, err = s.GuildMember(f.Config.GuildID, m.Author.ID)
		requester = mHandle
		if err != nil {
			return err
		}
	} else {
		idStr := strings.ReplaceAll(userSplit[1], "<@!", "")
		idStr = strings.ReplaceAll(idStr, "<@", "")
		idStr = strings.ReplaceAll(idStr, ">", "")
		mHandle, err = s.GuildMember(f.Config.GuildID, idStr)
		if err != nil {
			return err
		}
		requester, err = s.GuildMember(f.Config.GuildID, m.Author.ID)
		if err != nil {
			return err
		}
	}

	rUserID := mHandle.User.ID
	rUserNick := mHandle.Nick
	rUsername := mHandle.User.Username
	rUserDiscrim := mHandle.User.Discriminator
	rUserAvatar := mHandle.User.AvatarURL("4096")
	rJoinTime := mHandle.JoinedAt
	rRoles := mHandle.Roles

	if len(rUserNick) == 0 {
		rUserNick = "No Nickname"
	}

	rJoinTimeP, err := rJoinTime.Parse()
	if err != nil {
		return err
	}

	rRolesTidy := ""
	if len(rRoles) == 0 {
		rRolesTidy = "No Roles"
	} else {
		for _, role := range rRoles {
			rRolesTidy += "<@&" + role + "> "
		}
	}

	msg := "**User ID**: `" + rUserID + "`\n"
	msg += "**User Name**: `" + rUsername + "`\n"
	msg += "**User Nick**: `" + rUserNick + "`\n"
	msg += "**User Discrim**: `#" + rUserDiscrim + "`\n"
	msg += "**User Join**:  `" + rJoinTimeP.String() + "`\n"
	msg += "**User Roles**: " + rRolesTidy + "\n"

	embedData := models.CustomEmbed{
		URL:            "",
		Title:          "User Info (" + rUsername + ")",
		Desc:           msg,
		Type:           "",
		Timestamp:      time.Now().Format(time.RFC3339),
		Color:          0xFFA500,
		FooterText:     "Requested by " + requester.User.Username + "#" + requester.User.Discriminator,
		FooterImageURL: "",
		ImageURL:       "",
		ImageH:         0,
		ImageW:         0,
		ThumbnailURL:   rUserAvatar,
		ThumbnailH:     512,
		ThumbnailW:     512,
		ProviderURL:    "",
		ProviderText:   "",
		AuthorText:     "",
		AuthorURL:      "",
		AuthorImageURL: "",
	}

	embed := f.CreateCustomEmbed(&embedData)
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return err
	}

	return nil
}

func (f *Features) handleServerInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	g, err := s.Guild(f.Config.GuildID)
	if err != nil {
		return err
	}

	sID := f.Config.GuildID
	sName := f.Config.GuildName

	chans, _ := s.GuildChannels(f.Config.GuildID)
	sChannels := strconv.Itoa(len(chans))
	sEmojis := strconv.Itoa(len(g.Emojis))
	sRoles := strconv.Itoa(len(g.Roles))
	sRegion := g.Region

	iID, _ := strconv.Atoi(f.Config.GuildID)
	createdMSecs := ((iID / 4194304) + 1420070400000) / 1000
	sCreatedAt := time.Unix(int64(createdMSecs), 0).Format(time.RFC1123)

	sIconURL := g.IconURL()

	user := m.Author

	desc := "**Server ID**: `" + sID + "`\n"
	desc += "**Server Name**: `" + sName + "`\n"
	desc += "**Server Channels**: `" + sChannels + "`\n"
	desc += "**Server Emojis**: `" + sEmojis + "`\n"
	desc += "**Server Roles**: `" + sRoles + "`\n"
	desc += "**Server Region**: `" + sRegion + "`\n"
	desc += "**Server Creation**: `" + sCreatedAt + "`\n"

	embedData := models.CustomEmbed{
		URL:            "",
		Title:          "Server Info (" + sName + ")",
		Desc:           desc,
		Type:           "",
		Timestamp:      time.Now().Format(time.RFC3339),
		Color:          0xFFA500,
		FooterText:     "Requested by " + user.Username + "#" + user.Discriminator,
		FooterImageURL: "",
		ImageURL:       "",
		ImageH:         0,
		ImageW:         0,
		ThumbnailURL:   sIconURL,
		ThumbnailH:     256,
		ThumbnailW:     256,
		ProviderURL:    "",
		ProviderText:   "",
		AuthorText:     "",
		AuthorURL:      "",
		AuthorImageURL: "",
	}

	msg := f.CreateCustomEmbed(&embedData)

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, msg)
	if err != nil {
		return err
	}

	return nil
}
