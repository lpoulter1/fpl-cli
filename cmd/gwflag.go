package cmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type gwRange struct {
	Start int
	End   int
}

type gwFlag struct {
	ranges []gwRange
}

func (g *gwFlag) String() string {
	if len(g.ranges) == 0 {
		return ""
	}
	parts := make([]string, 0, len(g.ranges))
	for _, r := range g.ranges {
		if r.Start == r.End {
			parts = append(parts, fmt.Sprintf("%d", r.Start))
			continue
		}
		parts = append(parts, fmt.Sprintf("%d-%d", r.Start, r.End))
	}
	return strings.Join(parts, ",")
}

func (g *gwFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("gameweek value cannot be empty")
	}

	for _, token := range splitTokens(value) {
		gr, err := parseGWRange(token)
		if err != nil {
			return err
		}
		g.ranges = append(g.ranges, gr)
	}
	g.normalize()
	return nil
}

func (g *gwFlag) Type() string {
	return "gwrange"
}

func (g *gwFlag) normalize() {
	if len(g.ranges) == 0 {
		return
	}

	sort.Slice(g.ranges, func(i, j int) bool {
		if g.ranges[i].Start == g.ranges[j].Start {
			return g.ranges[i].End < g.ranges[j].End
		}
		return g.ranges[i].Start < g.ranges[j].Start
	})

	merged := make([]gwRange, 0, len(g.ranges))
	current := g.ranges[0]
	for _, r := range g.ranges[1:] {
		if r.Start <= current.End+1 {
			if r.End > current.End {
				current.End = r.End
			}
			continue
		}
		merged = append(merged, current)
		current = r
	}
	merged = append(merged, current)
	g.ranges = merged
}

func (g *gwFlag) includes(round int) bool {
	if len(g.ranges) == 0 {
		return true
	}
	for _, r := range g.ranges {
		if round >= r.Start && round <= r.End {
			return true
		}
	}
	return false
}

func (g *gwFlag) Ranges() []gwRange {
	return g.ranges
}

func parseGWRange(token string) (gwRange, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return gwRange{}, fmt.Errorf("gameweek token cannot be empty")
	}

	if strings.Contains(token, "-") {
		parts := strings.SplitN(token, "-", 2)
		start, err := parseGW(parts[0])
		if err != nil {
			return gwRange{}, err
		}
		end, err := parseGW(parts[1])
		if err != nil {
			return gwRange{}, err
		}
		if end < start {
			return gwRange{}, fmt.Errorf("invalid gameweek range %s: end before start", token)
		}
		return gwRange{Start: start, End: end}, nil
	}

	week, err := parseGW(token)
	if err != nil {
		return gwRange{}, err
	}
	return gwRange{Start: week, End: week}, nil
}

func parseGW(value string) (int, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, fmt.Errorf("gameweek value cannot be empty")
	}
	num, err := strconv.Atoi(value)
	if err != nil || num <= 0 {
		return 0, fmt.Errorf("gameweek must be a positive integer: %s", value)
	}
	return num, nil
}

func splitTokens(value string) []string {
	separators := func(r rune) bool {
		return r == ',' || r == '|' || r == ' '
	}
	raw := strings.FieldsFunc(value, separators)
	tokens := make([]string, 0, len(raw))
	for _, r := range raw {
		if trimmed := strings.TrimSpace(r); trimmed != "" {
			tokens = append(tokens, trimmed)
		}
	}
	return tokens
}
