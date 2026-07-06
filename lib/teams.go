package lib

import "strings"

// Team represents a national team in the 2026 FIFA World Cup.
type Team struct {
	Name string
	Abbr string
}

// Individual teams participating in the 2026 FIFA World Cup.
var (
	MEX = Team{Name: "Mexico", Abbr: "MEX"}
	KOR = Team{Name: "South Korea", Abbr: "KOR"}
	RSA = Team{Name: "South Africa", Abbr: "RSA"}
	CZE = Team{Name: "Czechia", Abbr: "CZE"}
	CAN = Team{Name: "Canada", Abbr: "CAN"}
	SUI = Team{Name: "Switzerland", Abbr: "SUI"}
	QAT = Team{Name: "Qatar", Abbr: "QAT"}
	BIH = Team{Name: "Bosnia and Herzegovina", Abbr: "BIH"}
	BRA = Team{Name: "Brazil", Abbr: "BRA"}
	MAR = Team{Name: "Morocco", Abbr: "MAR"}
	SCO = Team{Name: "Scotland", Abbr: "SCO"}
	HAI = Team{Name: "Haiti", Abbr: "HAI"}
	USA = Team{Name: "USA", Abbr: "USA"}
	PAR = Team{Name: "Paraguay", Abbr: "PAR"}
	AUS = Team{Name: "Australia", Abbr: "AUS"}
	TUR = Team{Name: "Turkiye", Abbr: "TUR"}
	GER = Team{Name: "Germany", Abbr: "GER"}
	ECU = Team{Name: "Ecuador", Abbr: "ECU"}
	CIV = Team{Name: "Ivory Coast", Abbr: "CIV"}
	CUW = Team{Name: "Curaçao", Abbr: "CUW"}
	NED = Team{Name: "Netherlands", Abbr: "NED"}
	JPN = Team{Name: "Japan", Abbr: "JPN"}
	TUN = Team{Name: "Tunisia", Abbr: "TUN"}
	SWE = Team{Name: "Sweden", Abbr: "SWE"}
	BEL = Team{Name: "Belgium", Abbr: "BEL"}
	IRN = Team{Name: "Iran", Abbr: "IRN"}
	EGY = Team{Name: "Egypt", Abbr: "EGY"}
	NZL = Team{Name: "New Zealand", Abbr: "NZL"}
	ESP = Team{Name: "Spain", Abbr: "ESP"}
	URU = Team{Name: "Uruguay", Abbr: "URU"}
	KSA = Team{Name: "Saudi Arabia", Abbr: "KSA"}
	CPV = Team{Name: "Cape Verde", Abbr: "CPV"}
	FRA = Team{Name: "France", Abbr: "FRA"}
	SEN = Team{Name: "Senegal", Abbr: "SEN"}
	NOR = Team{Name: "Norway", Abbr: "NOR"}
	IRQ = Team{Name: "Iraq", Abbr: "IRQ"}
	ARG = Team{Name: "Argentina", Abbr: "ARG"}
	AUT = Team{Name: "Austria", Abbr: "AUT"}
	ALG = Team{Name: "Algeria", Abbr: "ALG"}
	JOR = Team{Name: "Jordan", Abbr: "JOR"}
	POR = Team{Name: "Portugal", Abbr: "POR"}
	COL = Team{Name: "Colombia", Abbr: "COL"}
	UZB = Team{Name: "Uzbekistan", Abbr: "UZB"}
	COD = Team{Name: "DR Congo", Abbr: "COD"}
	ENG = Team{Name: "England", Abbr: "ENG"}
	CRO = Team{Name: "Croatia", Abbr: "CRO"}
	PAN = Team{Name: "Panama", Abbr: "PAN"}
	GHA = Team{Name: "Ghana", Abbr: "GHA"}

	Teams = []Team{
		MEX, KOR, RSA, CZE,
		CAN, SUI, QAT, BIH,
		BRA, MAR, SCO, HAI,
		USA, PAR, AUS, TUR,
		GER, ECU, CIV, CUW,
		NED, JPN, TUN, SWE,
		BEL, IRN, EGY, NZL,
		ESP, URU, KSA, CPV,
		FRA, SEN, NOR, IRQ,
		ARG, AUT, ALG, JOR,
		POR, COL, UZB, COD,
		ENG, CRO, PAN, GHA,
	}
)

// Group represents a World Cup group stage group.
type Group struct {
	Name  string
	Teams []Team
}

// Individual groups for the 2026 FIFA World Cup group stage.
var (
	GroupA = Group{
		Name:  "A",
		Teams: []Team{MEX, KOR, RSA, CZE},
	}
	GroupB = Group{
		Name:  "B",
		Teams: []Team{CAN, SUI, QAT, BIH},
	}
	GroupC = Group{
		Name:  "C",
		Teams: []Team{BRA, MAR, SCO, HAI},
	}
	GroupD = Group{
		Name:  "D",
		Teams: []Team{USA, PAR, AUS, TUR},
	}
	GroupE = Group{
		Name:  "E",
		Teams: []Team{GER, ECU, CIV, CUW},
	}
	GroupF = Group{
		Name:  "F",
		Teams: []Team{NED, JPN, TUN, SWE},
	}
	GroupG = Group{
		Name:  "G",
		Teams: []Team{BEL, IRN, EGY, NZL},
	}
	GroupH = Group{
		Name:  "H",
		Teams: []Team{ESP, URU, KSA, CPV},
	}
	GroupI = Group{
		Name:  "I",
		Teams: []Team{FRA, SEN, NOR, IRQ},
	}
	GroupJ = Group{
		Name:  "J",
		Teams: []Team{ARG, AUT, ALG, JOR},
	}
	GroupK = Group{
		Name:  "K",
		Teams: []Team{POR, COL, UZB, COD},
	}
	GroupL = Group{
		Name:  "L",
		Teams: []Team{ENG, CRO, PAN, GHA},
	}

	Groups = []Group{
		GroupA, GroupB, GroupC, GroupD,
		GroupE, GroupF, GroupG, GroupH,
		GroupI, GroupJ, GroupK, GroupL,
	}
)

// FindTeamByQuery finds the best matching team by abbreviation or name within candidates.
func FindTeamByQuery(query string, candidates []Team) (Team, bool) {
	q := strings.TrimSpace(strings.ToLower(query))
	if q == "" {
		return Team{}, false
	}

	for _, team := range candidates {
		if strings.ToLower(team.Abbr) == q {
			return team, true
		}
	}

	var abbrMatches []Team
	for _, team := range candidates {
		if strings.HasPrefix(strings.ToLower(team.Abbr), q) {
			abbrMatches = append(abbrMatches, team)
		}
	}
	if len(abbrMatches) > 0 {
		return abbrMatches[0], true
	}

	var nameMatches []Team
	for _, team := range candidates {
		name := strings.ToLower(team.Name)
		if strings.HasPrefix(name, q) || strings.Contains(name, q) {
			nameMatches = append(nameMatches, team)
		}
	}
	if len(nameMatches) > 0 {
		return nameMatches[0], true
	}

	return Team{}, false
}
