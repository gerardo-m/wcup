package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	tea "github.com/charmbracelet/bubbletea"
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
	matchIndex      int
	scoreField      int
	team1Input      string
	team2Input      string
	teamInput       string
	topScorer       string
	topScorerInput  string
	errMsg          string

	matchResults map[int]lib.MatchResult
	roundOf32    []lib.Team
	roundOf16    []lib.Team
	roundOf8     []lib.Team
	roundOf4     []lib.Team
	roundOf2     []lib.Team
	podium       []lib.Team

	groupQualifiers   []lib.GroupQualification
	qualifierGroupIdx int
	qualifierSlot     int
	pendingEditGroup  int
	pendingEditSlot   int
	teamSearchInput   string
	round32Reviewed   bool

	knockoutCursor int

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

func styleBold(s string) string   { return ansiBold + s + ansiReset }
func styleDim(s string) string    { return ansiDim + s + ansiReset }
func styleActive(s string) string { return ansiReverse + ansiBold + s + ansiReset }
func styleActivePadded(s string, width int) string {
	return styleActive(fmt.Sprintf("%-*s", width, s))
}
func styleError(s string) string  { return ansiRed + s + ansiReset }
func styleTitle(s string) string  { return ansiBold + ansiPink + s + ansiReset }

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
	m.round32Reviewed = len(m.roundOf32) == 32

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
	m.round32Reviewed = len(m.roundOf32) == 32

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
	m.teamSearchInput = ""
	m.knockoutCursor = 0
	m.errMsg = ""

	if m.phase == phaseRound32 && !m.round32Reviewed {
		m.ensureGroupQualifiers()
		m.pendingEditGroup = m.qualifierGroupIdx
		m.pendingEditSlot = m.qualifierSlot
	}
}

func (m *resultsEditor) ensureGroupQualifiers() {
	computed := lib.ComputeGroupQualifications(m.matchResults)
	if len(m.roundOf32) == 32 {
		m.groupQualifiers = lib.ApplyRoundOf32ToQualifications(computed, m.roundOf32)
	} else {
		m.groupQualifiers = computed
		m.roundOf32 = lib.RoundOf32FromQualifications(computed)
	}
}

func detectEditorProgress(m resultsEditor) (editorPhase, int) {
	for i, match := range lib.Matches {
		if _, ok := m.matchResults[match.Id]; !ok {
			return phaseMatches, i
		}
	}

	if !m.round32Reviewed || len(m.roundOf32) < 32 {
		return phaseRound32, 0
	}

	sections := []struct {
		teams []lib.Team
		max   int
		phase editorPhase
	}{
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
			} else if m.phase == phaseRound32 && m.qualifierSlot > 0 {
				m.qualifierSlot--
				m.pendingEditSlot = m.qualifierSlot
				m.teamSearchInput = ""
				m.errMsg = ""
			} else if m.isKnockoutSelectPhase() {
				m.moveKnockoutCursor(-1)
				m.teamSearchInput = ""
			}
		case "right":
			if m.phase == phaseRound32 && m.qualifierSlot < 2 {
				m.qualifierSlot++
				m.pendingEditSlot = m.qualifierSlot
				m.teamSearchInput = ""
				m.errMsg = ""
			} else if m.isKnockoutSelectPhase() {
				m.moveKnockoutCursor(1)
				m.teamSearchInput = ""
			}
		case "up":
			if m.phase == phaseRound32 && m.qualifierGroupIdx > 0 {
				m.qualifierGroupIdx--
				m.pendingEditGroup = m.qualifierGroupIdx
				m.teamSearchInput = ""
				m.errMsg = ""
			} else if m.isKnockoutSelectPhase() {
				m.moveKnockoutCursor(-m.knockoutColumns())
				m.teamSearchInput = ""
			}
		case "down":
			if m.phase == phaseRound32 && m.qualifierGroupIdx < len(m.groupQualifiers)-1 {
				m.qualifierGroupIdx++
				m.pendingEditGroup = m.qualifierGroupIdx
				m.teamSearchInput = ""
				m.errMsg = ""
			} else if m.isKnockoutSelectPhase() {
				m.moveKnockoutCursor(m.knockoutColumns())
				m.teamSearchInput = ""
			}
		case " ":
			if m.phase == phaseRound32 && m.qualifierSlot == 2 {
				m.toggleThirdQualified()
			} else if m.isKnockoutSelectPhase() {
				m.toggleKnockoutSelection()
			}
		case "backspace":
			m.handleBackspace()
		default:
			if len(msg.Runes) == 1 && msg.Runes[0] >= '0' && msg.Runes[0] <= '9' {
				m.handleDigit(msg.Runes[0])
			} else if m.phase == phaseTopScorer {
				m.handleTopScorerRune(msg.Runes)
			} else if m.phase == phaseRound32 || m.isKnockoutSelectPhase() {
				m.handleTeamSearchRune(msg.Runes)
			} else if m.phase == phasePodium {
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
		m.round32Reviewed = false
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

	case phaseRound32:
		return m.handleRound32Enter()

	case phaseRound16, phaseRound8, phaseRound4, phaseRound2:
		if strings.TrimSpace(m.teamSearchInput) != "" {
			m.teamSearchInput = ""
			m.errMsg = ""
			if len(*m.currentTeams()) != m.maxTeams() {
				return m, nil
			}
		}
		return m.handleKnockoutEnter()

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
	case phaseRound32:
		if len(m.teamSearchInput) > 0 {
			m.teamSearchInput = m.teamSearchInput[:len(m.teamSearchInput)-1]
			m.updateCursorFromSearch()
		}
	case phaseRound16, phaseRound8, phaseRound4, phaseRound2:
		if len(m.teamSearchInput) > 0 {
			m.teamSearchInput = m.teamSearchInput[:len(m.teamSearchInput)-1]
			m.updateCursorFromSearch()
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

func (m *resultsEditor) handleTeamSearchRune(runes []rune) {
	for _, r := range runes {
		if unicode.IsLetter(r) || r == ' ' {
			if len(m.teamSearchInput) < 40 {
				m.teamSearchInput += string(r)
			}
		}
	}
	m.updateCursorFromSearch()
	m.errMsg = ""
}

func (m *resultsEditor) updateCursorFromSearch() {
	query := strings.TrimSpace(m.teamSearchInput)
	if query == "" {
		return
	}

	switch {
	case m.phase == phaseRound32:
		if team, ok := lib.FindTeamByQuery(query, m.round32SearchPool()); ok {
			m.moveCursorToRound32Team(team)
		}
	case m.isKnockoutSelectPhase():
		if team, ok := lib.FindTeamByQuery(query, m.previousRoundTeams()); ok {
			m.moveKnockoutCursorToTeam(team)
		}
	}
}

func (m resultsEditor) round32SearchPool() []lib.Team {
	teams := make([]lib.Team, 0, len(m.groupQualifiers)*3)
	for _, qualification := range m.groupQualifiers {
		teams = append(teams, qualification.First, qualification.Second, qualification.Third)
	}
	return teams
}

func (m *resultsEditor) moveCursorToRound32Team(team lib.Team) {
	for groupIdx, qualification := range m.groupQualifiers {
		switch team.Abbr {
		case qualification.First.Abbr:
			m.qualifierGroupIdx = groupIdx
			m.qualifierSlot = 0
			return
		case qualification.Second.Abbr:
			m.qualifierGroupIdx = groupIdx
			m.qualifierSlot = 1
			return
		case qualification.Third.Abbr:
			m.qualifierGroupIdx = groupIdx
			m.qualifierSlot = 2
			return
		}
	}
}

func (m *resultsEditor) moveKnockoutCursorToTeam(team lib.Team) {
	for i, candidate := range m.previousRoundTeams() {
		if candidate.Abbr == team.Abbr {
			m.knockoutCursor = i
			return
		}
	}
}

func (m resultsEditor) handleRound32Enter() (tea.Model, tea.Cmd) {
	if strings.TrimSpace(m.teamSearchInput) != "" {
		group := lib.Groups[m.pendingEditGroup]
		team, ok := lib.FindTeamByQuery(m.teamSearchInput, group.Teams)
		if !ok {
			m.errMsg = fmt.Sprintf("Equipo desconocido: %s", strings.TrimSpace(m.teamSearchInput))
			return m, nil
		}

		qualification := &m.groupQualifiers[m.pendingEditGroup]
		switch m.pendingEditSlot {
		case 0:
			qualification.First = team
		case 1:
			qualification.Second = team
		case 2:
			qualification.Third = team
		}

		m.teamSearchInput = ""
		m.errMsg = ""
		m.moveCursorToRound32Team(team)
		m.roundOf32 = lib.RoundOf32FromQualifications(m.groupQualifiers)
		_ = m.syncToLib()
		return m, nil
	}

	if m.countQualifiedThirds() != 8 {
		m.errMsg = "Deben clasificar exactamente 8 terceros"
		return m, nil
	}

	m.roundOf32 = lib.RoundOf32FromQualifications(m.groupQualifiers)
	m.round32Reviewed = true
	m.errMsg = ""
	_ = m.syncToLib()
	m.reposition()
	if m.phase == phaseDone {
		return m, tea.Quit
	}
	return m, nil
}

func (m *resultsEditor) toggleThirdQualified() {
	if m.qualifierSlot != 2 {
		return
	}

	qualification := &m.groupQualifiers[m.qualifierGroupIdx]
	if qualification.ThirdQualified {
		qualification.ThirdQualified = false
	} else if m.countQualifiedThirds() >= 8 {
		m.errMsg = "Ya hay 8 terceros clasificados"
		return
	} else {
		qualification.ThirdQualified = true
	}

	m.roundOf32 = lib.RoundOf32FromQualifications(m.groupQualifiers)
	m.teamSearchInput = ""
	m.errMsg = ""
	_ = m.syncToLib()
}

func (m resultsEditor) countQualifiedThirds() int {
	count := 0
	for _, qualification := range m.groupQualifiers {
		if qualification.ThirdQualified {
			count++
		}
	}
	return count
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
	} else if m.phase == phaseRound32 {
		b.WriteString(m.renderRound32View())
		b.WriteString("\n\n")
		b.WriteString(styleDim("↑↓←→ navegar · Escribe abreviatura o país · Espacio alterna 3° · Enter confirma · Esc guarda y sale"))
	} else if m.isKnockoutSelectPhase() {
		b.WriteString(m.renderKnockoutView())
		b.WriteString("\n\n")
		b.WriteString(styleDim("↑↓←→ navegar · Escribe abreviatura o país · Espacio alternar · Enter confirma · Esc guarda y sale"))
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

func (m resultsEditor) renderRound32View() string {
	var b strings.Builder

	for i, qualification := range m.groupQualifiers {
		b.WriteString(fmt.Sprintf("  %s ", qualification.Group))
		b.WriteString(m.renderQualifierCell(i, 0, "1°", qualification.First.Abbr))
		b.WriteString("  ")
		b.WriteString(m.renderQualifierCell(i, 1, "2°", qualification.Second.Abbr))

		thirdLabel := "3°"
		if qualification.ThirdQualified {
			thirdLabel = "3°✓"
		}
		b.WriteString("  ")
		b.WriteString(m.renderQualifierCell(i, 2, thirdLabel, qualification.Third.Abbr))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	search := m.teamSearchInput
	if search == "" {
		search = "_"
	}
	slotLabels := []string{"1°", "2°", "3°"}
	b.WriteString(fmt.Sprintf("Grupo %s %s · Buscar: %s",
		m.groupQualifiers[m.pendingEditGroup].Group,
		slotLabels[m.pendingEditSlot],
		search,
	))
	b.WriteString(styleDim(fmt.Sprintf("  (%d/8 mejores terceros)", m.countQualifiedThirds())))

	return b.String()
}

func (m resultsEditor) renderQualifierCell(groupIdx, slot int, label, abbr string) string {
	cell := fmt.Sprintf("%s %s", label, abbr)
	if groupIdx == m.qualifierGroupIdx && slot == m.qualifierSlot {
		return styleActive(cell)
	}
	return cell
}

func (m resultsEditor) isKnockoutSelectPhase() bool {
	return m.phase >= phaseRound16 && m.phase <= phaseRound2
}

func (m resultsEditor) previousRoundTeams() []lib.Team {
	switch m.phase {
	case phaseRound16:
		return m.roundOf32
	case phaseRound8:
		return m.roundOf16
	case phaseRound4:
		return m.roundOf8
	case phaseRound2:
		return m.roundOf4
	default:
		return nil
	}
}

func (m resultsEditor) knockoutColumns() int {
	switch m.phase {
	case phaseRound16, phaseRound8:
		return 4
	default:
		return 2
	}
}

func (m *resultsEditor) moveKnockoutCursor(delta int) {
	pool := m.previousRoundTeams()
	next := m.knockoutCursor + delta
	if next >= 0 && next < len(pool) {
		m.knockoutCursor = next
		m.errMsg = ""
	}
}

func (m *resultsEditor) toggleKnockoutSelection() {
	pool := m.previousRoundTeams()
	if m.knockoutCursor >= len(pool) {
		return
	}

	team := pool[m.knockoutCursor]
	current := m.currentTeams()

	for i, existing := range *current {
		if existing.Abbr == team.Abbr {
			*current = append((*current)[:i], (*current)[i+1:]...)
			m.teamSearchInput = ""
			m.errMsg = ""
			_ = m.syncToLib()
			return
		}
	}

	if len(*current) >= m.maxTeams() {
		m.errMsg = fmt.Sprintf("Ya hay %d equipos seleccionados", m.maxTeams())
		return
	}

	*current = append(*current, team)
	m.teamSearchInput = ""
	m.errMsg = ""
	_ = m.syncToLib()
}

func (m resultsEditor) handleKnockoutEnter() (tea.Model, tea.Cmd) {
	current := m.currentTeams()
	if len(*current) != m.maxTeams() {
		m.errMsg = fmt.Sprintf("Selecciona exactamente %d equipos", m.maxTeams())
		return m, nil
	}

	m.errMsg = ""
	_ = m.syncToLib()
	m.reposition()
	if m.phase == phaseDone {
		return m, tea.Quit
	}
	return m, nil
}

func (m resultsEditor) renderKnockoutView() string {
	pool := m.previousRoundTeams()
	selected := make(map[string]struct{}, len(pool))
	for _, team := range *m.currentTeams() {
		selected[team.Abbr] = struct{}{}
	}

	cols := m.knockoutColumns()
	var b strings.Builder

	b.WriteString(styleDim("Equipos de la ronda anterior:"))
	b.WriteString("\n\n")

	for i, team := range pool {
		cell := team.Abbr
		if _, ok := selected[team.Abbr]; ok {
			cell += "✓"
		}
		if i == m.knockoutCursor {
			fmt.Fprintf(&b, "%s", styleActivePadded(cell, 6))
		} else {
			fmt.Fprintf(&b, "%-6s", cell)
		}
		if (i+1)%cols == 0 {
			b.WriteString("\n")
		}
	}
	if len(pool)%cols != 0 {
		b.WriteString("\n")
	}

	b.WriteString("\n")
	search := m.teamSearchInput
	if search == "" {
		search = "_"
	}
	b.WriteString(fmt.Sprintf("Buscar: %s", search))
	b.WriteString("\n")
	b.WriteString(styleDim(fmt.Sprintf("Seleccionados: %d/%d", len(*m.currentTeams()), m.maxTeams())))

	return b.String()
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
