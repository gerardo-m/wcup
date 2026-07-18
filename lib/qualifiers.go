package lib

import "sort"

// GroupQualification holds the group stage finishers and whether an extra team qualifies.
type GroupQualification struct {
	Group          string
	First          Team
	Second         Team
	Third          Team
	Fourth         Team
	ThirdQualified bool // group sends an extra team beyond 1st and 2nd
	ExtraIsFourth  bool // when ThirdQualified, use Fourth instead of Third
}

// ExtraQualifier returns the third-place qualifier team for the group, if any.
func (q GroupQualification) ExtraQualifier() (Team, bool) {
	if !q.ThirdQualified {
		return Team{}, false
	}
	if q.ExtraIsFourth {
		return q.Fourth, true
	}
	return q.Third, true
}

// ComputeGroupQualifications derives group finishers and the 8 best third-placed teams.
func ComputeGroupQualifications(results map[int]MatchResult) []GroupQualification {
	allStandings := AllGroupStandings(results)
	qualifications := make([]GroupQualification, 0, len(Groups))
	thirdPlaceCandidates := make([]thirdPlaceCandidate, 0, len(Groups))

	for i, group := range Groups {
		standings := allStandings[group.Name]
		qualification := GroupQualification{
			Group:  group.Name,
			First:  standings[0].Team,
			Second: standings[1].Team,
			Third:  standings[2].Team,
			Fourth: standings[3].Team,
		}
		qualifications = append(qualifications, qualification)
		thirdPlaceCandidates = append(thirdPlaceCandidates, thirdPlaceCandidate{
			groupIndex: i,
			standing:   standings[2],
		})
	}

	sort.Slice(thirdPlaceCandidates, func(i, j int) bool {
		return compareThirdPlaceStandings(thirdPlaceCandidates[i].standing, thirdPlaceCandidates[j].standing) > 0
	})

	for k := 0; k < 8 && k < len(thirdPlaceCandidates); k++ {
		qualifications[thirdPlaceCandidates[k].groupIndex].ThirdQualified = true
	}

	return qualifications
}

// RoundOf32FromQualifications builds the flat Round of 32 list from group qualifications.
func RoundOf32FromQualifications(qualifications []GroupQualification) []Team {
	teams := make([]Team, 0, 32)
	for _, qualification := range qualifications {
		teams = append(teams, qualification.First, qualification.Second)
		if extra, ok := qualification.ExtraQualifier(); ok {
			teams = append(teams, extra)
		}
	}
	return teams
}

// ApplyRoundOf32ToQualifications restores extra qualifiers from a saved Round of 32 list.
func ApplyRoundOf32ToQualifications(qualifications []GroupQualification, roundOf32 []Team) []GroupQualification {
	qualified := make(map[string]struct{}, len(roundOf32))
	for _, team := range roundOf32 {
		qualified[team.Abbr] = struct{}{}
	}

	for i := range qualifications {
		q := &qualifications[i]
		_, thirdIn := qualified[q.Third.Abbr]
		_, fourthIn := qualified[q.Fourth.Abbr]
		switch {
		case fourthIn && !thirdIn:
			q.ThirdQualified = true
			q.ExtraIsFourth = true
		case thirdIn:
			q.ThirdQualified = true
			q.ExtraIsFourth = false
		default:
			q.ThirdQualified = false
			q.ExtraIsFourth = false
		}
	}

	return qualifications
}

// TeamInGroup reports whether a team belongs to the given group.
func TeamInGroup(team Team, group Group) bool {
	for _, groupTeam := range group.Teams {
		if groupTeam.Abbr == team.Abbr {
			return true
		}
	}
	return false
}

type thirdPlaceCandidate struct {
	groupIndex int
	standing   TeamStanding
}

func compareThirdPlaceStandings(a, b TeamStanding) int {
	if a.Points != b.Points {
		return a.Points - b.Points
	}
	if a.GoalDifference != b.GoalDifference {
		return a.GoalDifference - b.GoalDifference
	}
	if a.GoalsFor != b.GoalsFor {
		return a.GoalsFor - b.GoalsFor
	}
	return 0
}
