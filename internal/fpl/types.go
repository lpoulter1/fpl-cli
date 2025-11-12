package fpl

import "time"

// BootstrapStatic mirrors the subset of data returned by /bootstrap-static/.
type BootstrapStatic struct {
	Elements     []Element     `json:"elements"`
	Teams        []Team        `json:"teams"`
	ElementTypes []ElementType `json:"element_types"`
}

// Element captures the player metadata needed for CLI output.
type Element struct {
	ID          int    `json:"id"`
	WebName     string `json:"web_name"`
	FirstName   string `json:"first_name"`
	SecondName  string `json:"second_name"`
	KnownAs     string `json:"known_as"`
	Team        int    `json:"team"`
	ElementType int    `json:"element_type"`
	NowCost     int    `json:"now_cost"`
	SelectedBy  string `json:"selected_by_percent"`
	TotalPoints int    `json:"total_points"`
	Form        string `json:"form"`
	ICTIndex    string `json:"ict_index"`
}

// Team describes a Premier League club.
type Team struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

// ElementType maps the player's position (e.g. Forward, Midfielder).
type ElementType struct {
	ID                int    `json:"id"`
	SingularName      string `json:"singular_name"`
	SingularNameShort string `json:"singular_name_short"`
}

// PlayerSummary is returned by /element-summary/{id}/.
type PlayerSummary struct {
	History []HistoryEntry `json:"history"`
}

// HistoryEntry represents the stats for a single gameweek.
type HistoryEntry struct {
	Round         int        `json:"round"`
	OpponentTeam  int        `json:"opponent_team"`
	WasHome       bool       `json:"was_home"`
	TotalPoints   int        `json:"total_points"`
	Minutes       int        `json:"minutes"`
	GoalsScored   int        `json:"goals_scored"`
	Assists       int        `json:"assists"`
	CleanSheets   int        `json:"clean_sheets"`
	GoalsConceded int        `json:"goals_conceded"`
	YellowCards   int        `json:"yellow_cards"`
	RedCards      int        `json:"red_cards"`
	BPS           int        `json:"bps"`
	Influence     string     `json:"influence"`
	Creativity    string     `json:"creativity"`
	Threat        string     `json:"threat"`
	ICTIndex      string     `json:"ict_index"`
	Value         int        `json:"value"`
	KickoffTime   *time.Time `json:"kickoff_time"`
}
