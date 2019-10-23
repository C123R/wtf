package football

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
)

var leagueID = map[string]leagueInfo{
	"BSA": {2013, "Brazil Série A"},
	"PL":  {2021, "English Premier League"},
	"EC":  {2016, "English Championship"},
	"EUC": {2018, "European Championship"},
	"EL2": {444, "Campeonato Brasileiro da Série A"},
	"CL":  {2001, "UEFA Champions League"},
	"FL1": {2015, "French Ligue 1"},
	"GB":  {2002, "German Bundesliga"},
	"ISA": {2019, "Italy Serie A"},
	"NE":  {2003, "Netherlands Eredivisie"},
	"PPL": {2017, "Portugal Primeira Liga"},
	"SPD": {2014, "Spain Primera Division"},
	"WC":  {2000, "FIFA World Cup"},
}

type Widget struct {
	view.TextWidget
	*Client
	settings *Settings
	League   leagueInfo
	err      error
}

func NewWidget(app *tview.Application, pages *tview.Pages, settings *Settings) *Widget {
	var widget Widget
	leagueId, err := getLeague(settings.league)
	if err != nil {
		widget = Widget{
			TextWidget: view.NewTextWidget(app, settings.common),
			err:        fmt.Errorf("Unable to get the league id for provided league '%s'", settings.league),
			Client:     NewClient(settings.apiKey),
			settings:   settings,
		}
		return &widget
	}
	widget = Widget{
		TextWidget: view.NewTextWidget(app, settings.common),
		Client:     NewClient(settings.apiKey),
		League:     leagueId,
		settings:   settings,
	}

	return &widget
}

func (widget *Widget) Refresh() {
	widget.Redraw(widget.content)
}

func (widget *Widget) content() (string, string, bool) {

	var content string
	title := fmt.Sprintf("%s %s", widget.CommonSettings().Title, widget.League.caption)
	wrap := false
	if widget.err != nil {
		return title, widget.err.Error(), true
	}

	content += widget.GetStandings(widget.League.id)
	content += widget.GetMatches(widget.League.id)

	return title, content, wrap
}

func getLeague(league string) (leagueInfo, error) {

	var l leagueInfo
	if val, ok := leagueID[league]; ok {
		return val, nil
	}
	return l, fmt.Errorf("No such league")
}

// GetStandings of particular league
func (widget *Widget) GetStandings(leagueId int) string {

	var l LeagueStandings
	var content string
	buf := new(bytes.Buffer)
	tStandings := createTable([]string{"No.", "Team", "MP", "Won", "Draw", "Lost", "GD", "Points"}, buf)
	resp, err := widget.Client.footballRequest("standings", leagueId)
	if err != nil {
		return fmt.Sprintf("Error fetching standings: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error fetching standings: %s", err.Error())
	}
	err = json.Unmarshal(data, &l)
	if err != nil {
		return fmt.Sprintf("Error fetching standings: %s", err.Error())
	}

	if err != nil {
		return fmt.Sprintf("Error reading standingCount: %s", err.Error())
	}
	if len(l.Standings) > 0 {
		for _, i := range l.Standings[0].Table {
			if i.Position <= widget.settings.standingCount {
				row := []string{strconv.Itoa(i.Position), i.Team.Name, strconv.Itoa(i.PlayedGames), strconv.Itoa(i.Won), strconv.Itoa(i.Draw), strconv.Itoa(i.Lost), strconv.Itoa(i.GoalDifference), strconv.Itoa(i.Points)}
				tStandings.Append(row)
			}

		}
	} else {
		return fmt.Sprintf("Error fetching standings!")
	}

	tStandings.Render()
	content += buf.String()

	return content
}

// GetMatches of particular league
func (widget *Widget) GetMatches(leagueId int) string {

	var l LeagueFixtuers
	var content string
	scheduledBuf := new(bytes.Buffer)
	playedBuf := new(bytes.Buffer)
	tScheduled := createTable([]string{}, scheduledBuf)
	tPlayed := createTable([]string{}, playedBuf)
	today := time.Now()

	requestPath := fmt.Sprintf("matches?dateFrom=%s&dateTo=%s", today.AddDate(0, 0, -widget.settings.matchesFrom).Format("2006-01-02"), today.AddDate(0, 0, widget.settings.matchesTo).Format("2006-01-02"))
	resp, err := widget.Client.footballRequest(requestPath, leagueId)
	if err != nil {
		return fmt.Sprintf("Error fetching matches: %s", err.Error())
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Error fetching matches: %s", err.Error())
	}
	err = json.Unmarshal(data, &l)
	if err != nil {
		return fmt.Sprintf("Error fetching matches: %s", err.Error())
	}

	for _, val := range l.Matches {
		if strings.Contains(val.AwayTeam.Name, widget.settings.favTeam) || strings.Contains(val.HomeTeam.Name, widget.settings.favTeam) || widget.settings.favTeam == "" {
			if val.Status == "SCHEDULED" {
				row := []string{"⚽", val.HomeTeam.Name, "🆚", val.AwayTeam.Name, parseDateString(val.Date)}
				tScheduled.Append(row)
			} else if val.Status == "FINISHED" {
				row := []string{"⚽", val.HomeTeam.Name, strconv.Itoa(val.Score.FullTime.HomeTeam), "🆚", val.AwayTeam.Name, strconv.Itoa(val.Score.FullTime.AwayTeam)}
				tPlayed.Append(row)
			}
		}
	}
	tScheduled.Render()
	tPlayed.Render()
	if playedBuf.String() != "" {
		content += "\nMatches Played:\n\n"
		content += playedBuf.String()

	}
	if scheduledBuf.String() != "" {
		content += "\nUpcoming Matches:\n\n"
		content += scheduledBuf.String()
	}

	return content
}
