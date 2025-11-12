package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/lpt10/fpl-cli/internal/fpl"
	"github.com/spf13/cobra"
)

type playerOptions struct {
	id   int
	name string
	gws  gwFlag
}

func newPlayerCmd() *cobra.Command {
	opts := &playerOptions{}
	cmd := &cobra.Command{
		Use:   "player",
		Short: "Show stats for an FPL player",
		Long: `Display Fantasy Premier League stats for a single player.

You can identify the target by ID (exact) or by name (fuzzy match). Gameweeks
can be filtered using --gw flags with single values or inclusive ranges.`,
		Example: `  fpl player --id 123
  fpl player --name "Haaland"
  fpl player --name "Haaland" --gw 1-3
  fpl player --name "Salah" --gw 1|4|6-8 --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPlayer(cmd.Context(), cmd, opts)
		},
	}

	cmd.Flags().IntVar(&opts.id, "id", 0, "FPL player ID to query")
	cmd.Flags().StringVar(&opts.name, "name", "", "player name to fuzzy match (web name, full name, or known-as)")
	cmd.Flags().Var(&opts.gws, "gw", "filter to a specific gameweek or inclusive range (e.g. --gw 5 --gw 1-3 --gw 6|8)")

	return cmd
}

func init() {
	rootCmd.AddCommand(newPlayerCmd())
}

func runPlayer(ctx context.Context, cmd *cobra.Command, opts *playerOptions) error {
	if opts.id <= 0 && strings.TrimSpace(opts.name) == "" {
		return errors.New("either --id or --name must be provided")
	}

	client := fpl.NewClient(nil, rootOpts.cacheTTL)
	bootstrap, err := client.Bootstrap(ctx)
	if err != nil {
		return err
	}

	var (
		target      *fpl.Element
		suggestions []fpl.MatchSuggestion
	)
	if opts.id > 0 {
		target = findElementByID(bootstrap.Elements, opts.id)
		if target == nil {
			return fmt.Errorf("player with ID %d not found in bootstrap data", opts.id)
		}
	} else {
		target, suggestions, err = fpl.FindPlayerByName(opts.name, bootstrap.Elements)
		if err != nil {
			return err
		}
	}

	summary, err := client.PlayerSummary(ctx, target.ID)
	if err != nil {
		return err
	}

	report := buildPlayerReport(target, bootstrap, summary, &opts.gws)
	if rootOpts.outputJSON {
		return printPlayerJSON(cmd, report)
	}

	return printPlayerTable(cmd, report, suggestions, opts.name)
}

func buildPlayerReport(player *fpl.Element, bootstrap *fpl.BootstrapStatic, summary *fpl.PlayerSummary, gw *gwFlag) playerReport {
	team := findTeam(bootstrap.Teams, player.Team)
	position := findElementType(bootstrap.ElementTypes, player.ElementType)

	filtered := make([]fpl.HistoryEntry, 0, len(summary.History))
	for _, entry := range summary.History {
		if gw == nil || gw.includes(entry.Round) {
			filtered = append(filtered, entry)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Round < filtered[j].Round
	})

	rows := make([]historyRow, 0, len(filtered))
	totals := historyTotals{}

	for _, entry := range filtered {
		rows = append(rows, historyRow{
			Round:       entry.Round,
			Opponent:    opponentLabel(entry, bootstrap.Teams),
			Home:        entry.WasHome,
			Minutes:     entry.Minutes,
			Goals:       entry.GoalsScored,
			Assists:     entry.Assists,
			CleanSheets: entry.CleanSheets,
			Points:      entry.TotalPoints,
		})
		totals.Gameweeks = append(totals.Gameweeks, entry.Round)
		totals.Matches++
		totals.Minutes += entry.Minutes
		totals.Goals += entry.GoalsScored
		totals.Assists += entry.Assists
		totals.CleanSheets += entry.CleanSheets
		totals.Points += entry.TotalPoints
	}

	return playerReport{
		Player: playerSummaryInfo{
			ID:          player.ID,
			Name:        playerDisplayName(player),
			Team:        teamLabel(team),
			Position:    positionLabel(position),
			Cost:        float64(player.NowCost) / 10.0,
			Form:        player.Form,
			ICTIndex:    player.ICTIndex,
			SelectedBy:  player.SelectedBy,
			TotalPoints: player.TotalPoints,
		},
		Gameweeks: rows,
		Totals:    totals,
	}
}

func printPlayerJSON(cmd *cobra.Command, report playerReport) error {
	out := cmd.OutOrStdout()
	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func printPlayerTable(cmd *cobra.Command, report playerReport, suggestions []fpl.MatchSuggestion, requestedName string) error {
	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "%s | %s | %s | £%.1f\n",
		report.Player.Name,
		report.Player.Team,
		report.Player.Position,
		report.Player.Cost,
	)
	fmt.Fprintf(out, "Form %s | Total Points %d | Selected by %s%% | ICT %s\n\n",
		report.Player.Form,
		report.Player.TotalPoints,
		report.Player.SelectedBy,
		report.Player.ICTIndex,
	)

	tw := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "GW\tOpponent\tMin\tG\tA\tCS\tPts")
	for _, row := range report.Gameweeks {
		fmt.Fprintf(tw, "%d\t%s\t%d\t%d\t%d\t%d\t%d\n",
			row.Round,
			row.Opponent,
			row.Minutes,
			row.Goals,
			row.Assists,
			row.CleanSheets,
			row.Points,
		)
	}
	tw.Flush()

	if len(report.Gameweeks) > 0 {
		fmt.Fprintf(out, "\nTotals (GW %s): %d matches | %d pts | %d min | %d G | %d A | %d CS\n",
			formatGWList(report.Totals.Gameweeks),
			report.Totals.Matches,
			report.Totals.Points,
			report.Totals.Minutes,
			report.Totals.Goals,
			report.Totals.Assists,
			report.Totals.CleanSheets,
		)
	} else {
		fmt.Fprintln(out, "No fixtures recorded for the selected gameweeks.")
	}

	if shouldSuggestAlternatives(requestedName, suggestions) {
		fmt.Fprintln(out, "\nOther close matches:")
		for _, s := range suggestions {
			if s.Element == nil {
				continue
			}
			fmt.Fprintf(out, "- %s (alias: %s, distance: %d)\n",
				playerDisplayName(s.Element),
				s.Alias,
				s.Distance,
			)
		}
	}

	return nil
}

func shouldSuggestAlternatives(requestedName string, suggestions []fpl.MatchSuggestion) bool {
	if strings.TrimSpace(requestedName) == "" || len(suggestions) == 0 {
		return false
	}
	// Inform users when the best match still has a relatively large edit distance.
	return suggestions[0].Distance > 2 && len(suggestions) > 1
}

func opponentLabel(entry fpl.HistoryEntry, teams []fpl.Team) string {
	team := findTeam(teams, entry.OpponentTeam)
	label := "Unknown"
	if team != nil {
		label = team.ShortName
	}
	if entry.WasHome {
		return fmt.Sprintf("%s (H)", label)
	}
	return fmt.Sprintf("%s (A)", label)
}

func formatGWList(weeks []int) string {
	if len(weeks) == 0 {
		return "-"
	}
	copied := append([]int(nil), weeks...)
	sort.Ints(copied)

	var builder strings.Builder
	start, prev := copied[0], copied[0]
	for i := 1; i <= len(copied); i++ {
		if i == len(copied) || copied[i] != prev+1 {
			if builder.Len() > 0 {
				builder.WriteString(",")
			}
			if start == prev {
				builder.WriteString(fmt.Sprintf("%d", start))
			} else {
				builder.WriteString(fmt.Sprintf("%d-%d", start, prev))
			}
			if i < len(copied) {
				start = copied[i]
				prev = copied[i]
			}
			continue
		}
		prev = copied[i]
	}
	return builder.String()
}

func findElementByID(elements []fpl.Element, id int) *fpl.Element {
	for i := range elements {
		if elements[i].ID == id {
			return &elements[i]
		}
	}
	return nil
}

func findTeam(teams []fpl.Team, id int) *fpl.Team {
	for i := range teams {
		if teams[i].ID == id {
			return &teams[i]
		}
	}
	return nil
}

func findElementType(types []fpl.ElementType, id int) *fpl.ElementType {
	for i := range types {
		if types[i].ID == id {
			return &types[i]
		}
	}
	return nil
}

func teamLabel(team *fpl.Team) string {
	if team == nil {
		return "Unknown"
	}
	return fmt.Sprintf("%s (%s)", team.Name, team.ShortName)
}

func positionLabel(pos *fpl.ElementType) string {
	if pos == nil {
		return "Unknown"
	}
	return pos.SingularName
}

func playerDisplayName(player *fpl.Element) string {
	name := strings.TrimSpace(player.FirstName + " " + player.SecondName)
	if name == "" {
		return player.WebName
	}
	return name
}

type playerReport struct {
	Player    playerSummaryInfo `json:"player"`
	Gameweeks []historyRow      `json:"gameweeks"`
	Totals    historyTotals     `json:"totals"`
}

type playerSummaryInfo struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Team        string  `json:"team"`
	Position    string  `json:"position"`
	Cost        float64 `json:"cost"`
	Form        string  `json:"form"`
	ICTIndex    string  `json:"ict_index"`
	SelectedBy  string  `json:"selected_by_percent"`
	TotalPoints int     `json:"total_points"`
}

type historyRow struct {
	Round       int    `json:"round"`
	Opponent    string `json:"opponent"`
	Home        bool   `json:"home"`
	Minutes     int    `json:"minutes"`
	Goals       int    `json:"goals"`
	Assists     int    `json:"assists"`
	CleanSheets int    `json:"clean_sheets"`
	Points      int    `json:"points"`
}

type historyTotals struct {
	Gameweeks   []int `json:"gameweeks"`
	Matches     int   `json:"matches"`
	Minutes     int   `json:"minutes"`
	Goals       int   `json:"goals"`
	Assists     int   `json:"assists"`
	CleanSheets int   `json:"clean_sheets"`
	Points      int   `json:"points"`
}
