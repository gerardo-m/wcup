package lib

import "testing"

func TestComputeGroupQualifications(t *testing.T) {
	group := Group{
		Name:  "T",
		Teams: []Team{{Name: "A", Abbr: "AAA"}, {Name: "B", Abbr: "BBB"}, {Name: "C", Abbr: "CCC"}, {Name: "D", Abbr: "DDD"}},
	}

	originalGroups := Groups
	originalMatchesByGroup := MatchesByGroup
	t.Cleanup(func() {
		Groups = originalGroups
		MatchesByGroup = originalMatchesByGroup
	})

	Groups = []Group{group}
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

	qualifications := ComputeGroupQualifications(results)
	if len(qualifications) != 1 {
		t.Fatalf("expected 1 group qualification, got %d", len(qualifications))
	}

	qualification := qualifications[0]
	if qualification.First.Abbr != "CCC" || qualification.Second.Abbr != "AAA" {
		t.Fatalf("unexpected top two: %s, %s", qualification.First.Abbr, qualification.Second.Abbr)
	}
	if qualification.Third.Abbr != "BBB" {
		t.Fatalf("expected BBB third, got %s", qualification.Third.Abbr)
	}
	if !qualification.ThirdQualified {
		t.Fatal("expected sole third place team to qualify")
	}

	roundOf32 := RoundOf32FromQualifications(qualifications)
	if len(roundOf32) != 3 {
		t.Fatalf("expected 3 teams for single group test, got %d", len(roundOf32))
	}
}

func TestRoundOf32FromQualificationsCountsThirds(t *testing.T) {
	qualifications := make([]GroupQualification, 12)
	for i := range qualifications {
		qualifications[i] = GroupQualification{
			Group:          string(rune('A' + i)),
			First:          Team{Abbr: "A1"},
			Second:         Team{Abbr: "A2"},
			Third:          Team{Abbr: "A3"},
			Fourth:         Team{Abbr: "A4"},
			ThirdQualified: i < 8,
			ExtraIsFourth:  i == 0,
		}
	}

	roundOf32 := RoundOf32FromQualifications(qualifications)
	if len(roundOf32) != 32 {
		t.Fatalf("expected 32 teams, got %d", len(roundOf32))
	}
	if roundOf32[2].Abbr != "A4" {
		t.Fatalf("expected first group's extra qualifier to be A4, got %s", roundOf32[2].Abbr)
	}
}

func TestApplyRoundOf32ToQualificationsFourth(t *testing.T) {
	qualifications := []GroupQualification{{
		Group:  "A",
		First:  Team{Abbr: "A1"},
		Second: Team{Abbr: "A2"},
		Third:  Team{Abbr: "A3"},
		Fourth: Team{Abbr: "A4"},
	}}

	restored := ApplyRoundOf32ToQualifications(qualifications, []Team{
		{Abbr: "A1"}, {Abbr: "A2"}, {Abbr: "A4"},
	})
	if !restored[0].ThirdQualified || !restored[0].ExtraIsFourth {
		t.Fatal("expected fourth place to be restored as extra qualifier")
	}
}
