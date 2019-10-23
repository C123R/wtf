package football

import (
	"os"

	"github.com/olebedev/config"
	"github.com/wtfutil/wtf/cfg"
)

const (
	defaultFocusable = true
	defaultTitle     = "football"
)

type Settings struct {
	common        *cfg.Common
	apiKey        string `help:"Your Football-data API token."`
	league        string `help:"Name of the competition. For example PL"`
	favTeam       string `help:"Teams to follow in mentioned league"`
	standingCount int    `help:"Number of positions to be displayed in standings widget"`
	matchesFrom   int    `help:"List matches from number days before. For example 5 = Today - 5 days"`
	matchesTo     int    `help:"List matches till number of days. For example 5 = Today + 5 days"`
}

func NewSettingsFromYAML(name string, ymlConfig *config.Config, globalConfig *config.Config) *Settings {

	settings := Settings{
		common:        cfg.NewCommonSettingsFromModule(name, defaultTitle, defaultFocusable, ymlConfig, globalConfig),
		apiKey:        ymlConfig.UString("apiKey", ymlConfig.UString("apikey", os.Getenv("WTF_FOOTBALL_API_KEY"))),
		league:        ymlConfig.UString("league", ymlConfig.UString("league", os.Getenv("WTF_FOOTBALL_LEAGUE"))),
		favTeam:       ymlConfig.UString("favTeam", ymlConfig.UString("favTeam", os.Getenv("WTF_FOOTBALL_TEAM"))),
		standingCount: ymlConfig.UInt("standingCount", 5),
		matchesFrom:   ymlConfig.UInt("matchesFrom", 2),
		matchesTo:     ymlConfig.UInt("matchesTo", 5),
	}
	return &settings
}
