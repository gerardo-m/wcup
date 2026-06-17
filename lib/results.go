package lib

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	wcupDirName     = ".wcup"
	resultsFileName = "results"
)

var (
	MatchResults []MatchResult
	RoundOf32    []Team
	RoundOf16    []Team
	RoundOf8     []Team
	RoundOf4     []Team
	RoundOf2     []Team
	Podium       []Team
)

func resultsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, wcupDirName, resultsFileName), nil
}

func emptyResultsContent() string {
	return strings.Repeat(".\n", 6)
}

// EnsureResultsFile creates ~/.wcup/ and an empty results file if needed.
func EnsureResultsFile() error {
	path, err := resultsPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(emptyResultsContent()), 0o644)
}

// LoadResults reads match results and classified teams from ~/.wcup/results.
func LoadResults() error {
	path, err := resultsPath()
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	sections := splitSections(file)
	if len(sections) < 1 {
		return fmt.Errorf("results file: expected at least 1 section")
	}
	if len(sections) > 7 {
		return fmt.Errorf("results file: expected at most 7 sections, got %d", len(sections))
	}
	for len(sections) < 7 {
		sections = append(sections, nil)
	}

	teamsByAbbr := make(map[string]Team, len(Teams))
	for _, team := range Teams {
		teamsByAbbr[team.Abbr] = team
	}

	matchResults, err := parseMatchResults(sections[0])
	if err != nil {
		return err
	}
	MatchResults = matchResults

	expectedCounts := []int{32, 16, 8, 4, 2, 3}
	teamSections := []*[]Team{&RoundOf32, &RoundOf16, &RoundOf8, &RoundOf4, &RoundOf2, &Podium}

	for i, maxTeams := range expectedCounts {
		teams, err := parseTeams(sections[i+1], teamsByAbbr)
		if err != nil {
			return fmt.Errorf("section %d: %w", i+2, err)
		}
		if len(teams) > maxTeams {
			return fmt.Errorf("section %d: expected at most %d teams, got %d", i+2, maxTeams, len(teams))
		}
		*teamSections[i] = teams
	}

	return nil
}

func splitSections(file *os.File) [][]string {
	var sections [][]string
	var current []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "." {
			sections = append(sections, current)
			current = nil
			continue
		}
		if line == "" {
			continue
		}
		current = append(current, line)
	}
	sections = append(sections, current)

	return sections
}

func parseMatchResults(lines []string) ([]MatchResult, error) {
	results := make([]MatchResult, 0, len(lines))

	for _, line := range lines {
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return nil, fmt.Errorf("invalid match result %q: expected 3 numbers", line)
		}

		matchID, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, fmt.Errorf("invalid match id in %q: %w", line, err)
		}
		team1Score, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid team1 score in %q: %w", line, err)
		}
		team2Score, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, fmt.Errorf("invalid team2 score in %q: %w", line, err)
		}

		match := Matches[matchID-1]

		results = append(results, MatchResult{
			Match:      match,
			Team1Score: team1Score,
			Team2Score: team2Score,
		})
	}

	return results, nil
}

func parseTeams(lines []string, teamsByAbbr map[string]Team) ([]Team, error) {
	teams := make([]Team, 0, len(lines))

	for _, line := range lines {
		abbr := strings.TrimSpace(line)
		team, ok := teamsByAbbr[abbr]
		if !ok {
			return nil, fmt.Errorf("unknown team abbreviation %q", abbr)
		}
		teams = append(teams, team)
	}

	return teams, nil
}
