package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

const participantsDirName = "participants"

type Participant struct {
	Name       string
	Prediction *Prediction
}

type Prediction struct {
	Matches   []MatchResult
	RoundOf32 []Team
	RoundOf16 []Team
	RoundOf8  []Team
	RoundOf4  []Team
	RoundOf2  []Team
	Podium    []Team
	TopScorer string
}

var Participants []Participant

func participantsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, wcupDirName, participantsDirName), nil
}

// EnsureParticipantsDir creates ~/.wcup/participants/ if needed.
func EnsureParticipantsDir() error {
	dir, err := participantsDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0o755)
}

// EnsureParticipantFile creates an empty prediction file for a participant if needed.
func EnsureParticipantFile(name string) error {
	dir, err := participantsDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	if err := EnsureParticipantsDir(); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(emptyResultsContent()), 0o644)
}

// LoadParticipants reads all prediction files from ~/.wcup/participants/.
func LoadParticipants() error {
	dir, err := participantsDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			Participants = nil
			return nil
		}
		return err
	}

	participants := make([]Participant, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		path := filepath.Join(dir, name)
		prediction, err := loadPredictionFromPath(path, name)
		if err != nil {
			return fmt.Errorf("participant %q: %w", name, err)
		}

		participants = append(participants, Participant{
			Name:       name,
			Prediction: prediction,
		})
	}

	Participants = participants
	return nil
}

func loadPredictionFromPath(path, label string) (*Prediction, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sections, err := splitSections(file)
	if err != nil {
		return nil, err
	}

	return parsePredictionSections(sections, label)
}

func parsePredictionSections(sections [][]string, label string) (*Prediction, error) {
	if len(sections) < 1 {
		return nil, fmt.Errorf("%s: expected at least 1 section", label)
	}
	if len(sections) > resultsSectionCount {
		return nil, fmt.Errorf("%s: expected at most %d sections, got %d", label, resultsSectionCount, len(sections))
	}
	for len(sections) < resultsSectionCount {
		sections = append(sections, nil)
	}

	teamsByAbbr := make(map[string]Team, len(Teams))
	for _, team := range Teams {
		teamsByAbbr[team.Abbr] = team
	}

	matchResults, err := parseMatchResults(sections[0])
	if err != nil {
		return nil, err
	}

	prediction := &Prediction{Matches: matchResults}
	expectedCounts := []int{32, 16, 8, 4, 2, 3}
	teamSections := []*[]Team{
		&prediction.RoundOf32,
		&prediction.RoundOf16,
		&prediction.RoundOf8,
		&prediction.RoundOf4,
		&prediction.RoundOf2,
		&prediction.Podium,
	}

	for i, maxTeams := range expectedCounts {
		teams, err := parseTeams(sections[i+1], teamsByAbbr)
		if err != nil {
			return nil, fmt.Errorf("%s section %d: %w", label, i+2, err)
		}
		if len(teams) > maxTeams {
			return nil, fmt.Errorf("%s section %d: expected at most %d teams, got %d", label, i+2, maxTeams, len(teams))
		}
		*teamSections[i] = teams
	}

	topScorer, err := parseTopScorer(sections[7])
	if err != nil {
		return nil, fmt.Errorf("%s section 8: %w", label, err)
	}
	prediction.TopScorer = topScorer

	return prediction, nil
}

// LoadParticipantPrediction reads a single participant prediction file.
func LoadParticipantPrediction(name string) (*Prediction, error) {
	if err := EnsureParticipantFile(name); err != nil {
		return nil, err
	}

	dir, err := participantsDir()
	if err != nil {
		return nil, err
	}

	return loadPredictionFromPath(filepath.Join(dir, name), name)
}

// SaveParticipantPrediction writes a participant prediction file.
func SaveParticipantPrediction(name string, prediction *Prediction) error {
	if err := EnsureParticipantsDir(); err != nil {
		return err
	}

	dir, err := participantsDir()
	if err != nil {
		return err
	}

	content := formatResultsFile(
		prediction.Matches,
		prediction.RoundOf32,
		prediction.RoundOf16,
		prediction.RoundOf8,
		prediction.RoundOf4,
		prediction.RoundOf2,
		prediction.Podium,
		prediction.TopScorer,
	)

	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		return err
	}

	updateParticipantInMemory(name, prediction)
	return nil
}

func updateParticipantInMemory(name string, prediction *Prediction) {
	for i, participant := range Participants {
		if participant.Name == name {
			Participants[i].Prediction = prediction
			return
		}
	}

	Participants = append(Participants, Participant{
		Name:       name,
		Prediction: prediction,
	})
}

// ResetParticipant clears a single participant prediction file.
func ResetParticipant(name string) error {
	dir, err := participantsDir()
	if err != nil {
		return err
	}

	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("participant %q not found", name)
		}
		return err
	}

	if err := os.WriteFile(path, []byte(emptyResultsContent()), 0o644); err != nil {
		return err
	}

	clearParticipantInMemory(name)
	return nil
}

// ResetAllParticipants clears every participant prediction file.
func ResetAllParticipants() error {
	dir, err := participantsDir()
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			Participants = nil
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if err := os.WriteFile(path, []byte(emptyResultsContent()), 0o644); err != nil {
			return err
		}
	}

	Participants = nil
	return nil
}

func clearParticipantInMemory(name string) {
	for i, participant := range Participants {
		if participant.Name == name {
			Participants[i].Prediction = &Prediction{}
			return
		}
	}
}
