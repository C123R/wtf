package football

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/rivo/tview"
	"github.com/wtfutil/wtf/view"
)

var offset = 0

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
	view.KeyboardWidget
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
			KeyboardWidget: view.NewKeyboardWidget(app, pages, settings.common),
			TextWidget:     view.NewTextWidget(app, settings.common),
			err:            fmt.Errorf("Unable to get the league id for provided league '%s'", settings.league),
			Client:         NewClient(settings.apiKey),
			settings:       settings,
		}
		return &widget
	}
	widget = Widget{
		KeyboardWidget: view.NewKeyboardWidget(app, pages, settings.common),
		TextWidget:     view.NewTextWidget(app, settings.common),
		Client:         NewClient(settings.apiKey),
		League:         leagueId,
		settings:       settings,
	}
	widget.initializeKeyboardControls()
	widget.View.SetInputCapture(widget.InputCapture)
	widget.View.SetScrollable(true)
	widget.KeyboardWidget.SetView(widget.View)
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
	table, err := widget.GetStandings(widget.League.id)
	if err != nil {
		return title, err.Error(), true
	}

	if len(table) != 0 {
		content = "Standings:\n\n"
		buf := new(bytes.Buffer)
		tStandings := createTable([]string{"No.", "Team", "MP", "Won", "Draw", "Lost", "GD", "Points"}, buf)
		for _, val := range table {
			row := []string{strconv.Itoa(val.Position), val.Team.Name, strconv.Itoa(val.PlayedGames), strconv.Itoa(val.Won), strconv.Itoa(val.Draw), strconv.Itoa(val.Lost), strconv.Itoa(val.GoalDifference), strconv.Itoa(val.Points)}
			tStandings.Append(row)
		}
		tStandings.Render()
		content += buf.String()
		content += fmt.Sprintf("Offset: %d", offset)
	}

	matches, err := widget.GetMatches(widget.League.id)
	if err != nil {
		return title, err.Error(), true
	}
	if len(matches) != 0 {

		scheduledBuf := new(bytes.Buffer)
		playedBuf := new(bytes.Buffer)
		tScheduled := createTable([]string{}, scheduledBuf)
		tPlayed := createTable([]string{}, playedBuf)
		for _, val := range matches {
			if strings.Contains(val.AwayTeam.Name, widget.settings.team) || strings.Contains(val.HomeTeam.Name, widget.settings.team) || widget.settings.team == "" {
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
	}

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
func (widget *Widget) GetStandings(leagueId int) ([]Table, error) {

	var l LeagueStandings
	var table []Table
	fmt.Println("GetStandings")

	resp, err := widget.Client.footballRequest("standings", leagueId)
	if err != nil {
		return nil, fmt.Errorf("Error fetching standings: %s", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &l)
	if err != nil {
		return nil, err
	}

	if len(l.Standings) > 0 {
		for _, i := range l.Standings[0].Table {
			if i.Position <= 5 {
				table = append(table, Table{
					Position:       i.Position,
					PlayedGames:    i.PlayedGames,
					Draw:           i.Draw,
					Won:            i.Won,
					Points:         i.Points,
					GoalDifference: i.GoalDifference,
					Lost:           i.Lost,
					Team: Team{
						Name: i.Team.Name,
					},
				})
			}
		}
	} else {
		return table, fmt.Errorf("Error fetching standings")
	}

	return table, nil
}

// GetMatches of particular league
func (widget *Widget) GetMatches(leagueId int) ([]Matches, error) {

	var l LeagueFixtuers
	fmt.Println("GetMatches")
	dateFrom, dateTo := getDateFrame(offset)
	requestPath := fmt.Sprintf("matches?dateFrom=%s&dateTo=%s", dateFrom, dateTo)
	resp, err := widget.Client.footballRequest(requestPath, leagueId)
	if err != nil {
		return nil, fmt.Errorf("Error fetching matches: %s", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &l)
	if err != nil {
		return nil, err
	}

	return l.Matches, nil
}
