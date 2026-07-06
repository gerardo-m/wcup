package lib

import "testing"

func TestFindTeamByQuery(t *testing.T) {
	candidates := []Team{MEX, KOR, BRA, USA}

	tests := []struct {
		query    string
		expected string
	}{
		{"MEX", "MEX"},
		{"me", "MEX"},
		{"mexico", "MEX"},
		{"brazil", "BRA"},
		{"south", "KOR"},
		{"usa", "USA"},
	}

	for _, tt := range tests {
		team, ok := FindTeamByQuery(tt.query, candidates)
		if !ok {
			t.Fatalf("query %q: expected match", tt.query)
		}
		if team.Abbr != tt.expected {
			t.Fatalf("query %q: expected %s, got %s", tt.query, tt.expected, team.Abbr)
		}
	}
}
