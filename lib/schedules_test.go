package lib

import "testing"

func TestMatchesGroupOrder(t *testing.T) {
	BuildSchedule()

	if len(MatchesGroupOrder) != len(Matches) {
		t.Fatalf("expected %d matches, got %d", len(Matches), len(MatchesGroupOrder))
	}

	groupIndex := map[string]int{}
	for i, group := range Groups {
		groupIndex[group.Name] = i
	}

	prevGroup := -1
	prevID := 0
	for _, match := range MatchesGroupOrder {
		group := groupForTeam(match.Team1)
		idx, ok := groupIndex[group]
		if !ok {
			t.Fatalf("unknown group for match %d", match.Id)
		}
		if idx < prevGroup {
			t.Fatalf("groups out of order at match %d", match.Id)
		}
		if idx == prevGroup && match.Id < prevID {
			t.Fatalf("match IDs out of order within group %s: %d after %d", group, match.Id, prevID)
		}
		if idx != prevGroup {
			prevID = 0
		}
		prevGroup = idx
		prevID = match.Id
	}
}

func groupForTeam(team Team) string {
	for _, group := range Groups {
		for _, groupTeam := range group.Teams {
			if groupTeam.Abbr == team.Abbr {
				return group.Name
			}
		}
	}
	return ""
}
