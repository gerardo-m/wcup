package lib

import "testing"

func TestGroupStandingsTiebreakers(t *testing.T) {
	group := Group{
		Name:  "T",
		Teams: []Team{{Name: "A", Abbr: "AAA"}, {Name: "B", Abbr: "BBB"}, {Name: "C", Abbr: "CCC"}, {Name: "D", Abbr: "DDD"}},
	}

	MatchesByGroup = map[string][]Match{
		"T": {
			{1, group.Teams[0], group.Teams[1], scheduleDate(2026, 1, 1, 0, 0)},
			{2, group.Teams[2], group.Teams[3], scheduleDate(2026, 1, 1, 0, 0)},
			{3, group.Teams[0], group.Teams[2], scheduleDate(2026, 1, 2, 0, 0)},
			{4, group.Teams[1], group.Teams[3], scheduleDate(2026, 1, 2, 0, 0)},
			{5, group.Teams[0], group.Teams[3], scheduleDate(2026, 1, 3, 0, 0)},
			{6, group.Teams[1], group.Teams[2], scheduleDate(2026, 1, 3, 0, 0)},
		},
	}

	results := map[int]MatchResult{
		1: {Match: MatchesByGroup["T"][0], Team1Score: 1, Team2Score: 0},
		2: {Match: MatchesByGroup["T"][1], Team1Score: 1, Team2Score: 0},
		3: {Match: MatchesByGroup["T"][2], Team1Score: 0, Team2Score: 0},
		4: {Match: MatchesByGroup["T"][3], Team1Score: 1, Team2Score: 0},
		5: {Match: MatchesByGroup["T"][4], Team1Score: 0, Team2Score: 1},
		6: {Match: MatchesByGroup["T"][5], Team1Score: 0, Team2Score: 1},
	}

	standings := GroupStandings(group, results)

	if standings[0].Team.Abbr != "CCC" {
		t.Fatalf("expected CCC first with 7 pts, got %s", standings[0].Team.Abbr)
	}
	if standings[1].Team.Abbr != "AAA" {
		t.Fatalf("expected AAA second with 4 pts, got %s", standings[1].Team.Abbr)
	}
	if standings[2].Team.Abbr != "BBB" {
		t.Fatalf("expected BBB third on head-to-head over DDD, got %s", standings[2].Team.Abbr)
	}
	if standings[3].Team.Abbr != "DDD" {
		t.Fatalf("expected DDD last, got %s", standings[3].Team.Abbr)
	}
}

func TestGroupStandingsStepTwoGoalDifference(t *testing.T) {
	group := Group{
		Name:  "T",
		Teams: []Team{{Name: "A", Abbr: "AAA"}, {Name: "B", Abbr: "BBB"}, {Name: "C", Abbr: "CCC"}, {Name: "D", Abbr: "DDD"}},
	}

	MatchesByGroup = map[string][]Match{
		"T": {
			{1, group.Teams[0], group.Teams[1], scheduleDate(2026, 1, 1, 0, 0)},
			{2, group.Teams[2], group.Teams[3], scheduleDate(2026, 1, 1, 0, 0)},
			{3, group.Teams[0], group.Teams[2], scheduleDate(2026, 1, 2, 0, 0)},
			{4, group.Teams[1], group.Teams[3], scheduleDate(2026, 1, 2, 0, 0)},
			{5, group.Teams[0], group.Teams[3], scheduleDate(2026, 1, 3, 0, 0)},
			{6, group.Teams[1], group.Teams[2], scheduleDate(2026, 1, 3, 0, 0)},
		},
	}

	results := map[int]MatchResult{
		1: {Match: MatchesByGroup["T"][0], Team1Score: 1, Team2Score: 1},
		2: {Match: MatchesByGroup["T"][1], Team1Score: 0, Team2Score: 0},
		3: {Match: MatchesByGroup["T"][2], Team1Score: 2, Team2Score: 0},
		4: {Match: MatchesByGroup["T"][3], Team1Score: 1, Team2Score: 0},
		5: {Match: MatchesByGroup["T"][4], Team1Score: 1, Team2Score: 0},
		6: {Match: MatchesByGroup["T"][5], Team1Score: 0, Team2Score: 0},
	}

	standings := GroupStandings(group, results)

	if standings[0].Team.Abbr != "AAA" {
		t.Fatalf("expected AAA first on overall GD, got %s", standings[0].Team.Abbr)
	}
	if standings[1].Team.Abbr != "BBB" {
		t.Fatalf("expected BBB second, got %s", standings[1].Team.Abbr)
	}
}
