package fpl

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

// MatchSuggestion captures a fuzzy match candidate for a player name lookup.
type MatchSuggestion struct {
	Element  *Element
	Alias    string
	Distance int
}

// FindPlayerByName runs a fuzzy search across common player name variants.
func FindPlayerByName(query string, elements []Element) (*Element, []MatchSuggestion, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil, nil, errors.New("player name cannot be empty")
	}

	targets := make([]string, 0, len(elements)*2)
	targetMap := make(map[string]*Element, len(elements)*2)

	for i := range elements {
		el := &elements[i]
		for _, alias := range nameVariants(el) {
			key := fmt.Sprintf("%s|%d", alias, el.ID)
			if _, exists := targetMap[key]; exists {
				continue
			}
			targets = append(targets, key)
			targetMap[key] = el
		}
	}

	ranks := fuzzy.RankFindNormalizedFold(query, targets)
	if len(ranks) == 0 {
		return nil, nil, fmt.Errorf("no players found matching %q", query)
	}

	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].Distance < ranks[j].Distance
	})

	best := ranks[0]
	bestElement := targetMap[best.Target]
	suggestions := make([]MatchSuggestion, 0, min(5, len(ranks)))
	for i := 0; i < len(ranks) && i < 5; i++ {
		r := ranks[i]
		el := targetMap[r.Target]
		suggestions = append(suggestions, MatchSuggestion{
			Element:  el,
			Alias:    aliasFromKey(r.Target),
			Distance: r.Distance,
		})
	}

	return bestElement, suggestions, nil
}

func nameVariants(el *Element) []string {
	candidates := []string{
		el.WebName,
		strings.TrimSpace(el.FirstName + " " + el.SecondName),
		el.SecondName,
		el.KnownAs,
	}

	unique := make([]string, 0, len(candidates))
	seen := make(map[string]struct{}, len(candidates))
	for _, c := range candidates {
		c = strings.TrimSpace(c)
		if c == "" {
			continue
		}
		cLower := strings.ToLower(c)
		if _, ok := seen[cLower]; ok {
			continue
		}
		seen[cLower] = struct{}{}
		unique = append(unique, c)
	}
	return unique
}

func aliasFromKey(key string) string {
	parts := strings.SplitN(key, "|", 2)
	if len(parts) == 0 {
		return key
	}
	return parts[0]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
