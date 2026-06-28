package lib

import (
	"sort"
	"strings"
)

// ParticipantScore holds the point breakdown for a participant.
type ParticipantScore struct {
	Name         string
	MatchPoints  int
	Round32      int
	Round16      int
	QuarterFinal int
	SemiFinal    int
	ThirdPlace   int
	Final        int
	TopScorer    int
	Champion     int
	Total        int
}

// CalculateParticipantPoints compares a prediction against the current results.
func CalculateParticipantPoints(prediction *Prediction) ParticipantScore {
	score := ParticipantScore{}

	actualMatches := MatchResultsByID()
	predictedMatches := matchResultsByID(prediction.Matches)

	for matchID, predicted := range predictedMatches {
		actual, ok := actualMatches[matchID]
		if !ok {
			continue
		}
		score.MatchPoints += matchPredictionPoints(predicted, actual)
	}

	score.Round32 = teamListPoints(prediction.RoundOf32, RoundOf32, 2)
	score.Round16 = teamListPoints(prediction.RoundOf16, RoundOf16, 3)
	score.QuarterFinal = teamListPoints(prediction.RoundOf8, RoundOf8, 4)
	score.SemiFinal = teamListPoints(prediction.RoundOf4, RoundOf4, 6)
	score.Final = teamListPoints(prediction.RoundOf2, RoundOf2, 8)
	score.ThirdPlace = podiumPlacePoints(prediction.Podium, Podium, 2, 3)
	score.Champion = podiumPlacePoints(prediction.Podium, Podium, 0, 10)

	if strings.TrimSpace(prediction.TopScorer) != "" &&
		strings.EqualFold(strings.TrimSpace(prediction.TopScorer), strings.TrimSpace(TopScorer)) {
		score.TopScorer = 5
	}

	score.Total = score.MatchPoints + score.Round32 + score.Round16 + score.QuarterFinal +
		score.SemiFinal + score.ThirdPlace + score.Final + score.TopScorer + score.Champion

	return score
}

// CalculateAllParticipantPoints scores every loaded participant.
func CalculateAllParticipantPoints() []ParticipantScore {
	scores := make([]ParticipantScore, 0, len(Participants))
	for _, participant := range Participants {
		if participant.Prediction == nil {
			continue
		}
		score := CalculateParticipantPoints(participant.Prediction)
		score.Name = participant.Name
		scores = append(scores, score)
	}

	sort.Slice(scores, func(i, j int) bool {
		if scores[i].Total != scores[j].Total {
			return scores[i].Total > scores[j].Total
		}
		return scores[i].Name < scores[j].Name
	})

	return scores
}

func matchResultsByID(results []MatchResult) map[int]MatchResult {
	byID := make(map[int]MatchResult, len(results))
	for _, result := range results {
		byID[result.Match.Id] = result
	}
	return byID
}

func matchPredictionPoints(predicted, actual MatchResult) int {
	actualDraw := actual.Team1Score == actual.Team2Score
	predictedDraw := predicted.Team1Score == predicted.Team2Score
	exactScore := predicted.Team1Score == actual.Team1Score && predicted.Team2Score == actual.Team2Score

	if actualDraw {
		if !predictedDraw {
			return 0
		}
		if exactScore {
			return 5
		}
		return 4
	}

	if predictedDraw {
		return 0
	}

	if predictedWinner(predicted) != predictedWinner(actual) {
		return 0
	}

	if exactScore {
		return 3
	}
	return 2
}

func predictedWinner(result MatchResult) string {
	if result.Team1Score > result.Team2Score {
		return result.Match.Team1.Abbr
	}
	if result.Team2Score > result.Team1Score {
		return result.Match.Team2.Abbr
	}
	return ""
}

func teamListPoints(predicted, actual []Team, pointsPerTeam int) int {
	if len(actual) == 0 {
		return 0
	}

	actualTeams := teamAbbrSet(actual)
	total := 0
	for _, team := range predicted {
		if _, ok := actualTeams[team.Abbr]; ok {
			total += pointsPerTeam
		}
	}
	return total
}

func podiumPlacePoints(predicted, actual []Team, index, points int) int {
	if len(actual) <= index || len(predicted) <= index {
		return 0
	}
	if predicted[index].Abbr == actual[index].Abbr {
		return points
	}
	return 0
}

func teamAbbrSet(teams []Team) map[string]struct{} {
	set := make(map[string]struct{}, len(teams))
	for _, team := range teams {
		set[team.Abbr] = struct{}{}
	}
	return set
}
