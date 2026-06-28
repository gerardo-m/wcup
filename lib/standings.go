package lib

import "sort"

// TeamStanding holds group stage statistics for a team.
type TeamStanding struct {
	Team           Team
	Points         int
	Wins           int
	Draws          int
	Losses         int
	GoalsFor       int
	GoalsAgainst   int
	GoalDifference int
	Position       int
}

type miniStats struct {
	points int
	gd     int
	gf     int
}

// GroupStandings calculates ranked standings for a group from registered match results.
func GroupStandings(group Group, results map[int]MatchResult) []TeamStanding {
	matches := MatchesByGroup[group.Name]
	standings := make([]TeamStanding, len(group.Teams))

	for i, team := range group.Teams {
		standings[i] = teamStanding(team, matches, results)
	}

	sortStandings(standings, matches, results)

	for i := range standings {
		standings[i].Position = i + 1
	}

	return standings
}

// AllGroupStandings calculates standings for every group.
func AllGroupStandings(results map[int]MatchResult) map[string][]TeamStanding {
	standings := make(map[string][]TeamStanding, len(Groups))
	for _, group := range Groups {
		standings[group.Name] = GroupStandings(group, results)
	}
	return standings
}

func teamStanding(team Team, matches []Match, results map[int]MatchResult) TeamStanding {
	standing := TeamStanding{Team: team}

	for _, match := range matches {
		result, ok := results[match.Id]
		if !ok {
			continue
		}

		switch team.Abbr {
		case match.Team1.Abbr:
			standing.GoalsFor += result.Team1Score
			standing.GoalsAgainst += result.Team2Score
			switch {
			case result.Team1Score > result.Team2Score:
				standing.Wins++
				standing.Points += 3
			case result.Team1Score < result.Team2Score:
				standing.Losses++
			default:
				standing.Draws++
				standing.Points++
			}
		case match.Team2.Abbr:
			standing.GoalsFor += result.Team2Score
			standing.GoalsAgainst += result.Team1Score
			switch {
			case result.Team2Score > result.Team1Score:
				standing.Wins++
				standing.Points += 3
			case result.Team2Score < result.Team1Score:
				standing.Losses++
			default:
				standing.Draws++
				standing.Points++
			}
		}
	}

	standing.GoalDifference = standing.GoalsFor - standing.GoalsAgainst
	return standing
}

func sortStandings(standings []TeamStanding, matches []Match, results map[int]MatchResult) {
	sort.Slice(standings, func(i, j int) bool {
		return compareStandings(standings[i], standings[j], standings, matches, results) > 0
	})
}

func compareStandings(a, b TeamStanding, standings []TeamStanding, matches []Match, results map[int]MatchResult) int {
	if a.Points != b.Points {
		return a.Points - b.Points
	}

	tied := tiedTeamsOnPoints(standings, a.Points)
	if len(tied) >= 2 {
		miniA := miniGroupStats(a.Team, tied, matches, results)
		miniB := miniGroupStats(b.Team, tied, matches, results)

		if miniA.points != miniB.points {
			return miniA.points - miniB.points
		}
		if miniA.gd != miniB.gd {
			return miniA.gd - miniB.gd
		}
		if miniA.gf != miniB.gf {
			return miniA.gf - miniB.gf
		}
	}

	if a.GoalDifference != b.GoalDifference {
		return a.GoalDifference - b.GoalDifference
	}
	if a.GoalsFor != b.GoalsFor {
		return a.GoalsFor - b.GoalsFor
	}

	return 0
}

func tiedTeamsOnPoints(standings []TeamStanding, points int) []Team {
	tied := make([]Team, 0, len(standings))
	for _, standing := range standings {
		if standing.Points == points {
			tied = append(tied, standing.Team)
		}
	}
	return tied
}

func miniGroupStats(team Team, subset []Team, matches []Match, results map[int]MatchResult) miniStats {
	subsetAbbrs := make(map[string]struct{}, len(subset))
	for _, t := range subset {
		subsetAbbrs[t.Abbr] = struct{}{}
	}

	var stats miniStats

	for _, match := range matches {
		if _, ok := subsetAbbrs[match.Team1.Abbr]; !ok {
			continue
		}
		if _, ok := subsetAbbrs[match.Team2.Abbr]; !ok {
			continue
		}

		result, ok := results[match.Id]
		if !ok {
			continue
		}

		switch team.Abbr {
		case match.Team1.Abbr:
			stats.gf += result.Team1Score
			stats.gd += result.Team1Score - result.Team2Score
			switch {
			case result.Team1Score > result.Team2Score:
				stats.points += 3
			case result.Team1Score == result.Team2Score:
				stats.points++
			}
		case match.Team2.Abbr:
			stats.gf += result.Team2Score
			stats.gd += result.Team2Score - result.Team1Score
			switch {
			case result.Team2Score > result.Team1Score:
				stats.points += 3
			case result.Team2Score == result.Team1Score:
				stats.points++
			}
		}
	}

	return stats
}

func MatchResultsByID() map[int]MatchResult {
	results := make(map[int]MatchResult, len(MatchResults))
	for _, result := range MatchResults {
		results[result.Match.Id] = result
	}
	return results
}
