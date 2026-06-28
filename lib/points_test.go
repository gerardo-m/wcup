package lib

import "testing"

func TestMatchPredictionPoints(t *testing.T) {
	match := Match{Id: 1, Team1: MEX, Team2: RSA}

	tests := []struct {
		name     string
		predict  MatchResult
		actual   MatchResult
		expected int
	}{
		{
			name:     "exact score and winner",
			predict:  MatchResult{Match: match, Team1Score: 2, Team2Score: 1},
			actual:   MatchResult{Match: match, Team1Score: 2, Team2Score: 1},
			expected: 3,
		},
		{
			name:     "winner only",
			predict:  MatchResult{Match: match, Team1Score: 3, Team2Score: 0},
			actual:   MatchResult{Match: match, Team1Score: 1, Team2Score: 0},
			expected: 2,
		},
		{
			name:     "exact draw",
			predict:  MatchResult{Match: match, Team1Score: 1, Team2Score: 1},
			actual:   MatchResult{Match: match, Team1Score: 1, Team2Score: 1},
			expected: 5,
		},
		{
			name:     "draw result",
			predict:  MatchResult{Match: match, Team1Score: 0, Team2Score: 0},
			actual:   MatchResult{Match: match, Team1Score: 1, Team2Score: 1},
			expected: 4,
		},
		{
			name:     "wrong winner",
			predict:  MatchResult{Match: match, Team1Score: 2, Team2Score: 0},
			actual:   MatchResult{Match: match, Team1Score: 0, Team2Score: 1},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchPredictionPoints(tt.predict, tt.actual)
			if got != tt.expected {
				t.Fatalf("expected %d, got %d", tt.expected, got)
			}
		})
	}
}

func TestCalculateParticipantPoints(t *testing.T) {
	MatchResults = []MatchResult{
		{Match: Match{Id: 1, Team1: MEX, Team2: RSA}, Team1Score: 2, Team2Score: 1},
	}
	RoundOf32 = []Team{MEX, USA}
	RoundOf16 = []Team{MEX}
	RoundOf8 = []Team{MEX}
	RoundOf4 = []Team{MEX}
	RoundOf2 = []Team{MEX, BRA}
	Podium = []Team{MEX, BRA, ARG}
	TopScorer = "Lionel Messi"

	prediction := &Prediction{
		Matches: []MatchResult{
			{Match: Match{Id: 1, Team1: MEX, Team2: RSA}, Team1Score: 2, Team2Score: 1},
		},
		RoundOf32: []Team{MEX, CAN},
		RoundOf16: []Team{MEX},
		RoundOf8:  []Team{MEX},
		RoundOf4:  []Team{MEX},
		RoundOf2:  []Team{MEX, ARG},
		Podium:    []Team{MEX, ARG, BRA},
		TopScorer: "Lionel Messi",
	}

	score := CalculateParticipantPoints(prediction)

	if score.MatchPoints != 3 {
		t.Fatalf("match points: got %d", score.MatchPoints)
	}
	if score.Round32 != 2 {
		t.Fatalf("round32 points: got %d", score.Round32)
	}
	if score.Final != 8 {
		t.Fatalf("final points: got %d", score.Final)
	}
	if score.Champion != 10 {
		t.Fatalf("champion points: got %d", score.Champion)
	}
	if score.TopScorer != 5 {
		t.Fatalf("top scorer points: got %d", score.TopScorer)
	}
	if score.ThirdPlace != 0 {
		t.Fatalf("third place points: got %d", score.ThirdPlace)
	}
	if score.Total != 41 {
		t.Fatalf("total points: got %d", score.Total)
	}
}
