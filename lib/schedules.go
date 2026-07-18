package lib

import (
	"sort"
	"time"
)

type Match struct {
	Id    int
	Team1 Team
	Team2 Team
	Date  time.Time
}

type MatchResult struct {
	Match      Match
	Team1Score int
	Team2Score int
}

func (mr MatchResult) Winner() *Team {
	if mr.Team1Score > mr.Team2Score {
		return &mr.Match.Team1
	}
	if mr.Team2Score > mr.Team1Score {
		return &mr.Match.Team2
	}
	return nil
}

func (mr MatchResult) IsDraw() bool {
	return mr.Team1Score == mr.Team2Score
}

var Matches []Match

// MatchesByGroup indexes group stage matches by group name (e.g. "A", "B").
var MatchesByGroup map[string][]Match

// MatchesByTeam indexes group stage matches by team abbreviation (e.g. "MEX", "USA").
var MatchesByTeam map[string][]Match

// MatchesGroupOrder lists group stage matches ordered by group (A–L), then by match ID.
var MatchesGroupOrder []Match

// scheduleDate builds a match date using the official local kick-off time.
func scheduleDate(year int, month time.Month, day, hour, min int) time.Time {
	return time.Date(year, month, day, hour, min, 0, 0, time.UTC)
}

// BuildSchedule loads the official FIFA World Cup 2026 group stage fixtures.
func BuildSchedule() {
	Matches = []Match{
		{1, MEX, RSA, scheduleDate(2026, time.June, 11, 13, 0)},
		{2, KOR, CZE, scheduleDate(2026, time.June, 11, 20, 0)},
		{3, CAN, BIH, scheduleDate(2026, time.June, 12, 15, 0)},
		{4, USA, PAR, scheduleDate(2026, time.June, 12, 18, 0)},
		{5, HAI, SCO, scheduleDate(2026, time.June, 13, 21, 0)},
		{6, AUS, TUR, scheduleDate(2026, time.June, 13, 21, 0)},
		{7, BRA, MAR, scheduleDate(2026, time.June, 13, 18, 0)},
		{8, QAT, SUI, scheduleDate(2026, time.June, 13, 12, 0)},
		{9, CIV, ECU, scheduleDate(2026, time.June, 14, 19, 0)},
		{10, GER, CUW, scheduleDate(2026, time.June, 14, 12, 0)},
		{11, NED, JPN, scheduleDate(2026, time.June, 14, 15, 0)},
		{12, SWE, TUN, scheduleDate(2026, time.June, 14, 20, 0)},
		{13, KSA, URU, scheduleDate(2026, time.June, 15, 18, 0)},
		{14, ESP, CPV, scheduleDate(2026, time.June, 15, 12, 0)},
		{15, IRN, NZL, scheduleDate(2026, time.June, 15, 18, 0)},
		{16, BEL, EGY, scheduleDate(2026, time.June, 15, 12, 0)},
		{17, FRA, SEN, scheduleDate(2026, time.June, 16, 15, 0)},
		{18, IRQ, NOR, scheduleDate(2026, time.June, 16, 18, 0)},
		{19, ARG, ALG, scheduleDate(2026, time.June, 16, 20, 0)},
		{20, AUT, JOR, scheduleDate(2026, time.June, 16, 21, 0)},
		{21, GHA, PAN, scheduleDate(2026, time.June, 17, 19, 0)},
		{22, ENG, CRO, scheduleDate(2026, time.June, 17, 15, 0)},
		{23, POR, COD, scheduleDate(2026, time.June, 17, 12, 0)},
		{24, UZB, COL, scheduleDate(2026, time.June, 17, 20, 0)},
		{25, CZE, RSA, scheduleDate(2026, time.June, 18, 12, 0)},
		{26, SUI, BIH, scheduleDate(2026, time.June, 18, 12, 0)},
		{27, CAN, QAT, scheduleDate(2026, time.June, 18, 15, 0)},
		{28, MEX, KOR, scheduleDate(2026, time.June, 18, 19, 0)},
		{29, BRA, HAI, scheduleDate(2026, time.June, 19, 21, 0)},
		{30, SCO, MAR, scheduleDate(2026, time.June, 19, 18, 0)},
		{31, TUR, PAR, scheduleDate(2026, time.June, 19, 20, 0)},
		{32, USA, AUS, scheduleDate(2026, time.June, 19, 12, 0)},
		{33, GER, CIV, scheduleDate(2026, time.June, 20, 16, 0)},
		{34, ECU, CUW, scheduleDate(2026, time.June, 20, 19, 0)},
		{35, NED, SWE, scheduleDate(2026, time.June, 20, 12, 0)},
		{36, TUN, JPN, scheduleDate(2026, time.June, 20, 22, 0)},
		{37, URU, CPV, scheduleDate(2026, time.June, 21, 18, 0)},
		{38, ESP, KSA, scheduleDate(2026, time.June, 21, 12, 0)},
		{39, BEL, IRN, scheduleDate(2026, time.June, 21, 12, 0)},
		{40, NZL, EGY, scheduleDate(2026, time.June, 21, 18, 0)},
		{41, NOR, SEN, scheduleDate(2026, time.June, 22, 20, 0)},
		{42, FRA, IRQ, scheduleDate(2026, time.June, 22, 17, 0)},
		{43, ARG, AUT, scheduleDate(2026, time.June, 22, 12, 0)},
		{44, JOR, ALG, scheduleDate(2026, time.June, 22, 20, 0)},
		{45, ENG, GHA, scheduleDate(2026, time.June, 23, 16, 0)},
		{46, PAN, CRO, scheduleDate(2026, time.June, 23, 19, 0)},
		{47, POR, UZB, scheduleDate(2026, time.June, 23, 12, 0)},
		{48, COL, COD, scheduleDate(2026, time.June, 23, 20, 0)},
		{49, SCO, BRA, scheduleDate(2026, time.June, 24, 18, 0)},
		{50, MAR, HAI, scheduleDate(2026, time.June, 24, 18, 0)},
		{51, SUI, CAN, scheduleDate(2026, time.June, 24, 12, 0)},
		{52, BIH, QAT, scheduleDate(2026, time.June, 24, 12, 0)},
		{53, CZE, MEX, scheduleDate(2026, time.June, 24, 19, 0)},
		{54, RSA, KOR, scheduleDate(2026, time.June, 24, 19, 0)},
		{55, CUW, CIV, scheduleDate(2026, time.June, 25, 16, 0)},
		{56, ECU, GER, scheduleDate(2026, time.June, 25, 16, 0)},
		{57, JPN, SWE, scheduleDate(2026, time.June, 25, 18, 0)},
		{58, TUN, NED, scheduleDate(2026, time.June, 25, 18, 0)},
		{59, TUR, USA, scheduleDate(2026, time.June, 25, 19, 0)},
		{60, PAR, AUS, scheduleDate(2026, time.June, 25, 19, 0)},
		{61, NOR, FRA, scheduleDate(2026, time.June, 26, 15, 0)},
		{62, SEN, IRQ, scheduleDate(2026, time.June, 26, 15, 0)},
		{63, EGY, IRN, scheduleDate(2026, time.June, 26, 20, 0)},
		{64, NZL, BEL, scheduleDate(2026, time.June, 26, 20, 0)},
		{65, CPV, KSA, scheduleDate(2026, time.June, 26, 19, 0)},
		{66, URU, ESP, scheduleDate(2026, time.June, 26, 18, 0)},
		{67, PAN, ENG, scheduleDate(2026, time.June, 27, 17, 0)},
		{68, CRO, GHA, scheduleDate(2026, time.June, 27, 17, 0)},
		{69, ALG, AUT, scheduleDate(2026, time.June, 27, 21, 0)},
		{70, JOR, ARG, scheduleDate(2026, time.June, 27, 21, 0)},
		{71, COL, POR, scheduleDate(2026, time.June, 27, 19, 30)},
		{72, COD, UZB, scheduleDate(2026, time.June, 27, 19, 30)},
	}
	buildMatchIndexes()
}

func buildMatchIndexes() {
	teamToGroup := make(map[string]string, len(Teams))
	for _, group := range Groups {
		for _, team := range group.Teams {
			teamToGroup[team.Abbr] = group.Name
		}
	}

	MatchesByGroup = make(map[string][]Match, len(Groups))
	MatchesByTeam = make(map[string][]Match, len(Teams))

	for _, match := range Matches {
		group := teamToGroup[match.Team1.Abbr]
		MatchesByGroup[group] = append(MatchesByGroup[group], match)

		MatchesByTeam[match.Team1.Abbr] = append(MatchesByTeam[match.Team1.Abbr], match)
		MatchesByTeam[match.Team2.Abbr] = append(MatchesByTeam[match.Team2.Abbr], match)
	}

	MatchesGroupOrder = make([]Match, 0, len(Matches))
	for _, group := range Groups {
		groupMatches := append([]Match(nil), MatchesByGroup[group.Name]...)
		sort.Slice(groupMatches, func(i, j int) bool {
			return groupMatches[i].Id < groupMatches[j].Id
		})
		MatchesGroupOrder = append(MatchesGroupOrder, groupMatches...)
	}
}
