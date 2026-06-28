package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbletea"
	"github.com/gerardo-m/wcup/lib"
)

type editorPhase int

const (
	phaseMatches editorPhase = iota
	phaseRound32
	phaseRound16
	phaseRound8
	phaseRound4
	phaseRound2
	phasePodium
	phaseTopScorer
	phaseDone
)

type resultsEditor struct {
	participantName string
	phase           editorPhase
	matchIndex int
	scoreField int
	team1Input string
	team2Input string
	teamInput        string
	topScorer        string
	topScorerInput   string
	errMsg           string

	matchResults map[int]lib.MatchResult
	roundOf32    []lib.Team
	roundOf16    []lib.Team
	roundOf8     []lib.Team
	roundOf4     []lib.Team
	roundOf2     []lib.Team
	podium       []lib.Team

	teamsByAbbr map[string]lib.Team
}

const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiDim     = "\033[2m"
	ansiReverse = "\033[7m"
	ansiRed     = "\033[31m"
	ansiPink    = "\033[95m"
	ansiGray    = "\033[90m"
)

func styleBold(s string) string    { return ansiBold + s + ansiReset }
func styleDim(s string) string     { return ansiDim + s + ansiReset }
func styleActive(s string) string  { return ansiReverse + ansiBold + s + ansiReset }
func styleError(s string) string   { return ansiRed + s + ansiReset }
func styleTitle(s string) string   { return ansiBold + ansiPink + s + ansiReset }

func runResultsEditor() error {
	return runEditor(newResultsEditor())
}

func runPredictionEditor(participant string) error {
	m, err := newPredictionEditor(participant)
	if err != nil {
		return err
	}
	return runEditor(m)
}

func runEditor(m resultsEditor) error {
	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func newResultsEditor() resultsEditor {
	m := resultsEditor{
		teamsByAbbr:  lib.TeamsByAbbr(),
		matchResults: make(map[int]lib.MatchResult, len(lib.MatchResults)),
	}
	for _, result := range lib.MatchResults {
		m.matchResults[result.Match.Id] = result
	}

	m.roundOf32 = append([]lib.Team(nil), lib.RoundOf32...)
	m.roundOf16 = append([]lib.Team(nil), lib.RoundOf16...)
	m.roundOf8 = append([]lib.Team(nil), lib.RoundOf8...)
	m.roundOf4 = append([]lib.Team(nil), lib.RoundOf4...)
	m.roundOf2 = append([]lib.Team(nil), lib.RoundOf2...)
	m.podium = append([]lib.Team(nil), lib.Podium...)
	m.topScorer = lib.TopScorer

	m.reposition()
	return m
}

func newPredictionEditor(participant string) (resultsEditor, error) {
	prediction, err := lib.LoadParticipantPrediction(participant)
	if err != nil {
		return resultsEditor{}, err
	}

	m := resultsEditor{
		participantName: participant,
		teamsByAbbr:     lib.TeamsByAbbr(),
		matchResults:    make(map[int]lib.MatchResult, len(prediction.Matches)),
	}
	for _, result := range prediction.Matches {
		m.matchResults[result.Match.Id] = result
	}

	m.roundOf32 = append([]lib.Team(nil), prediction.RoundOf32...)
	m.roundOf16 = append([]lib.Team(nil), prediction.RoundOf16...)
	m.roundOf8 = append([]lib.Team(nil), prediction.RoundOf8...)
	m.roundOf4 = append([]lib.Team(nil), prediction.RoundOf4...)
	m.roundOf2 = append([]lib.Team(nil), prediction.RoundOf2...)
	m.podium = append([]lib.Team(nil), prediction.Podium...)
	m.topScorer = prediction.TopScorer

	m.reposition()
	return m, nil
}

func (m *resultsEditor) reposition() {
	m.phase, m.matchIndex = detectEditorProgress(*m)
	m.scoreField = 0
	m.team1Input = ""
	m.team2Input = ""
	m.teamInput = ""
	m.topScorerInput = ""
	m.errMsg = ""
}

func detectEditorProgress(m resultsEditor) (editorPhase, int) {
	for i, match := range lib.Matches {
		if _, ok := m.matchResults[match.Id]; !ok {
			return phaseMatches, i
		}
	}

	sections := []struct {
		teams []lib.Team
		max   int
		phase editorPhase
	}{
		{m.roundOf32, 32, phaseRound32},
		{m.roundOf16, 16, phaseRound16},
		{m.roundOf8, 8, phaseRound8},
		{m.roundOf4, 4, phaseRound4},
		{m.roundOf2, 2, phaseRound2},
		{m.podium, 3, phasePodium},
	}

	for _, section := range sections {
		if len(section.teams) < section.max {
			return section.phase, 0
		}
	}

	if strings.TrimSpace(m.topScorer) == "" {
		return phaseTopScorer, 0
	}

	return phaseDone, 0
}

func (m resultsEditor) Init() tea.Cmd {
	return nil
}

func (m resultsEditor) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			_ = m.syncToLib()
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		case "left":
			if m.phase == phaseMatches && m.scoreField == 1 {
				m.scoreField = 0
				m.errMsg = ""
			}
		case "backspace":
			m.handleBackspace()
		default:
			if len(msg.Runes) == 1 && msg.Runes[0] >= '0' && msg.Runes[0] <= '9' {
				m.handleDigit(msg.Runes[0])
			} else if m.phase == phaseTopScorer {
				m.handleTopScorerRune(msg.Runes)
			} else if m.phase != phaseMatches && m.phase != phaseDone {
				m.handleTeamRune(msg.Runes)
			}
		}
	}

	return m, nil
}

func (m resultsEditor) handleEnter() (tea.Model, tea.Cmd) {
	switch m.phase {
	case phaseMatches:
		if m.scoreField == 0 {
			if strings.TrimSpace(m.team1Input) == "" {
				m.team1Input = "0"
			}
			m.scoreField = 1
			m.errMsg = ""
			return m, nil
		}

		if strings.TrimSpace(m.team2Input) == "" {
			m.team2Input = "0"
		}

		team1Score, err1 := strconv.Atoi(m.team1Input)
		team2Score, err2 := strconv.Atoi(m.team2Input)
		if err1 != nil || err2 != nil {
			m.errMsg = "Marcador inválido"
			return m, nil
		}

		match := lib.Matches[m.matchIndex]
		m.matchResults[match.Id] = lib.MatchResult{
			Match:      match,
			Team1Score: team1Score,
			Team2Score: team2Score,
		}
		_ = m.syncToLib()
		m.reposition()
		if m.phase == phaseDone {
			return m, tea.Quit
		}
		return m, nil

	case phaseDone:
		return m, tea.Quit

	case phaseTopScorer:
		name := strings.TrimSpace(m.topScorerInput)
		if name == "" {
			m.errMsg = "Ingresa el nombre del goleador"
			return m, nil
		}
		m.topScorer = name
		m.topScorerInput = ""
		m.errMsg = ""
		_ = m.syncToLib()
		m.reposition()
		if m.phase == phaseDone {
			return m, tea.Quit
		}
		return m, nil

	default:
		abbr := strings.ToUpper(strings.TrimSpace(m.teamInput))
		if abbr == "" {
			m.errMsg = "Ingresa una abreviatura"
			return m, nil
		}

		team, ok := m.teamsByAbbr[abbr]
		if !ok {
			m.errMsg = fmt.Sprintf("Equipo desconocido: %s", abbr)
			return m, nil
		}

		current := m.currentTeams()
		for _, existing := range *current {
			if existing.Abbr == team.Abbr {
				m.errMsg = fmt.Sprintf("%s ya está en la lista", team.Abbr)
				return m, nil
			}
		}

		*current = append(*current, team)
		m.teamInput = ""
		m.errMsg = ""
		_ = m.syncToLib()

		if len(*current) >= m.maxTeams() {
			m.reposition()
			if m.phase == phaseDone {
				return m, tea.Quit
			}
		}
		return m, nil
	}
}

func (m *resultsEditor) handleBackspace() {
	switch m.phase {
	case phaseMatches:
		if m.scoreField == 0 && len(m.team1Input) > 0 {
			m.team1Input = m.team1Input[:len(m.team1Input)-1]
		} else if m.scoreField == 1 && len(m.team2Input) > 0 {
			m.team2Input = m.team2Input[:len(m.team2Input)-1]
		}
	case phaseDone:
		return
	case phaseTopScorer:
		if len(m.topScorerInput) > 0 {
			m.topScorerInput = m.topScorerInput[:len(m.topScorerInput)-1]
		}
	default:
		if len(m.teamInput) > 0 {
			m.teamInput = m.teamInput[:len(m.teamInput)-1]
		}
	}
	m.errMsg = ""
}

func (m *resultsEditor) handleDigit(d rune) {
	if m.phase != phaseMatches {
		return
	}

	digit := string(d)
	if m.scoreField == 0 {
		if len(m.team1Input) < 2 {
			m.team1Input += digit
		}
	} else if len(m.team2Input) < 2 {
		m.team2Input += digit
	}
	m.errMsg = ""
}

func (m *resultsEditor) handleTopScorerRune(runes []rune) {
	for _, r := range runes {
		if unicode.IsLetter(r) || r == ' ' || r == '-' || r == '\'' {
			if len(m.topScorerInput) < 40 {
				m.topScorerInput += string(r)
			}
		}
	}
	m.errMsg = ""
}

func (m *resultsEditor) handleTeamRune(runes []rune) {
	for _, r := range runes {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			if len(m.teamInput) < 3 {
				m.teamInput += strings.ToUpper(string(r))
			}
		}
	}
	m.errMsg = ""
}

func (m *resultsEditor) currentTeams() *[]lib.Team {
	switch m.phase {
	case phaseRound32:
		return &m.roundOf32
	case phaseRound16:
		return &m.roundOf16
	case phaseRound8:
		return &m.roundOf8
	case phaseRound4:
		return &m.roundOf4
	case phaseRound2:
		return &m.roundOf2
	case phasePodium:
		return &m.podium
	default:
		return &m.roundOf32
	}
}

func (m resultsEditor) maxTeams() int {
	switch m.phase {
	case phaseRound32:
		return 32
	case phaseRound16:
		return 16
	case phaseRound8:
		return 8
	case phaseRound4:
		return 4
	case phaseRound2:
		return 2
	case phasePodium:
		return 3
	default:
		return 0
	}
}

func (m resultsEditor) teamColumns() int {
	switch m.phase {
	case phaseRound32, phaseRound16:
		return 4
	default:
		return 2
	}
}

func (m resultsEditor) phaseTitle() string {
	var title string
	switch m.phase {
	case phaseMatches:
		title = fmt.Sprintf("PARTIDOS (%d/%d)", len(m.matchResults), len(lib.Matches))
	case phaseRound32:
		title = fmt.Sprintf("RONDA DE 32 (%d/32)", len(m.roundOf32))
	case phaseRound16:
		title = fmt.Sprintf("OCTAVOS (%d/16)", len(m.roundOf16))
	case phaseRound8:
		title = fmt.Sprintf("CUARTOS (%d/8)", len(m.roundOf8))
	case phaseRound4:
		title = fmt.Sprintf("SEMIFINALES (%d/4)", len(m.roundOf4))
	case phaseRound2:
		title = fmt.Sprintf("FINAL (%d/2)", len(m.roundOf2))
	case phasePodium:
		title = fmt.Sprintf("PODIO (%d/3)", len(m.podium))
	case phaseTopScorer:
		title = "GOLEADOR"
	default:
		if m.participantName != "" {
			return "PREDICCIÓN COMPLETA"
		}
		return "RESULTADOS COMPLETOS"
	}

	if m.participantName != "" {
		return fmt.Sprintf("%s · %s", m.participantName, title)
	}
	return title
}

func (m resultsEditor) syncToLib() error {
	prediction := &lib.Prediction{
		Matches:   m.matchResultsSlice(),
		RoundOf32: append([]lib.Team(nil), m.roundOf32...),
		RoundOf16: append([]lib.Team(nil), m.roundOf16...),
		RoundOf8:  append([]lib.Team(nil), m.roundOf8...),
		RoundOf4:  append([]lib.Team(nil), m.roundOf4...),
		RoundOf2:  append([]lib.Team(nil), m.roundOf2...),
		Podium:    append([]lib.Team(nil), m.podium...),
		TopScorer: m.topScorer,
	}

	if m.participantName != "" {
		return lib.SaveParticipantPrediction(m.participantName, prediction)
	}

	lib.MatchResults = prediction.Matches
	lib.RoundOf32 = prediction.RoundOf32
	lib.RoundOf16 = prediction.RoundOf16
	lib.RoundOf8 = prediction.RoundOf8
	lib.RoundOf4 = prediction.RoundOf4
	lib.RoundOf2 = prediction.RoundOf2
	lib.Podium = prediction.Podium
	lib.TopScorer = prediction.TopScorer
	return lib.SaveResults()
}

func (m resultsEditor) matchResultsSlice() []lib.MatchResult {
	results := make([]lib.MatchResult, 0, len(m.matchResults))
	for _, result := range m.matchResults {
		results = append(results, result)
	}
	for i := 1; i < len(results); i++ {
		for j := i; j > 0 && results[j].Match.Id < results[j-1].Match.Id; j-- {
			results[j], results[j-1] = results[j-1], results[j]
		}
	}
	return results
}

func (m resultsEditor) View() string {
	if m.phase == phaseDone {
		doneTitle := "RESULTADOS COMPLETOS"
		doneMsg := "Todos los resultados han sido registrados."
		if m.participantName != "" {
			doneTitle = "PREDICCIÓN COMPLETA"
			doneMsg = fmt.Sprintf("La predicción de %s ha sido registrada.", m.participantName)
		}
		return styleTitle(doneTitle) + "\n\n" +
			styleDim(doneMsg) + "\n\n" +
			styleDim("Presiona Enter o Esc para salir.")
	}

	var b strings.Builder
	b.WriteString(styleTitle(m.phaseTitle()))
	b.WriteString("\n\n")

	if m.phase == phaseMatches {
		b.WriteString(m.renderMatchView())
		b.WriteString("\n\n")
		b.WriteString(styleDim("Escribe el marcador · Enter avanza · ← retrocede · Esc guarda y sale"))
	} else if m.phase == phaseTopScorer {
		b.WriteString(m.renderTopScorerView())
		b.WriteString("\n\n")
		b.WriteString(styleDim("Escribe el nombre del goleador · Enter guarda · Esc guarda y sale"))
	} else {
		b.WriteString(m.renderTeamView())
		b.WriteString("\n\n")
		b.WriteString(styleDim("Escribe la abreviatura del equipo · Enter agrega · Esc guarda y sale"))
	}

	if m.errMsg != "" {
		b.WriteString("\n\n")
		b.WriteString(styleError(m.errMsg))
	}

	return b.String()
}

func (m resultsEditor) renderMatchView() string {
	match := lib.Matches[m.matchIndex]
	team1Score := m.displayScore(m.team1Input)
	team2Score := m.displayScore(m.team2Input)

	if m.scoreField == 0 {
		team1Score = styleActive(team1Score)
	} else {
		team2Score = styleActive(team2Score)
	}

	meta := styleDim(fmt.Sprintf("#%d · Grupo %s", match.Id, groupForMatch(match)))
	line := fmt.Sprintf("  %s  %s - %s  %s", match.Team1.Abbr, team1Score, team2Score, match.Team2.Abbr)

	return meta + "\n" + line
}

func (m resultsEditor) displayScore(input string) string {
	if input == "" {
		return "_"
	}
	return input
}

func (m resultsEditor) renderTeamView() string {
	teams := *m.currentTeams()
	cols := m.teamColumns()
	max := m.maxTeams()

	var b strings.Builder
	for i := 0; i < max; i++ {
		cell := "____"
		if i < len(teams) {
			cell = teams[i].Abbr
		}

		b.WriteString(fmt.Sprintf("%-6s", cell))
		if (i+1)%cols == 0 {
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Equipo: %s", m.teamInput))
	if m.phase == phasePodium && len(teams) < 3 {
		labels := []string{"1°", "2°", "3°"}
		b.WriteString(styleDim(fmt.Sprintf("  (%s lugar)", labels[len(teams)])))
	}

	return b.String()
}

func (m resultsEditor) renderTopScorerView() string {
	name := m.topScorerInput
	if name == "" {
		name = "_"
	}
	return fmt.Sprintf("Nombre: %s", styleActive(name))
}

func groupForMatch(match lib.Match) string {
	for _, group := range lib.Groups {
		for _, team := range group.Teams {
			if team.Abbr == match.Team1.Abbr {
				return group.Name
			}
		}
	}
	return "?"
}
