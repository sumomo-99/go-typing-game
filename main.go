package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Difficulty represents the difficulty level of the game
type Difficulty struct {
	Name          string
	Description   string
	TimeLimit     time.Duration
	WordLengthMin int
	WordLengthMax int
}

// Predefined difficulties
var difficulties = []Difficulty{
	{
		Name:          "Easy",
		Description:   "短い単語や簡単な文章が出題され、制限時間が長めに設定される。",
		TimeLimit:     60 * time.Second,
		WordLengthMin: 1,
		WordLengthMax: 5,
	},
	{
		Name:          "Normal",
		Description:   "中程度の長さの単語や文章が出題され、制限時間が標準的に設定される。",
		TimeLimit:     45 * time.Second,
		WordLengthMin: 6,
		WordLengthMax: 10,
	},
	{
		Name:          "Hard",
		Description:   "長い単語や難しい文章が出題され、制限時間が短めに設定される。",
		TimeLimit:     30 * time.Second,
		WordLengthMin: 11,
		WordLengthMax: 20,
	},
}

// Game represents the game state
type Game struct {
	Difficulty Difficulty
	Words      []string
	Score      int
	Mistakes   int
	TimeLimit  time.Duration
	CharCount int
	StartTime time.Time
}

// LoadWords loads words from a file
func LoadWords(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	return words, scanner.Err()
}

// Model represents the TUI state
type Model struct {
	state       string
	difficulties []Difficulty
	selected    int
	game        *Game
	timer       *time.Timer
	input       string
	currentWord string
	index       int
}

func initialModel() Model {
	return Model{
		state:        "menu",
		difficulties: difficulties,
		selected:     0,
	}
}

// Init initializes the program
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case "menu":
			switch msg.String() {
			case "up":
				if m.selected > 0 {
					m.selected--
				}
			case "down":
				if m.selected < len(m.difficulties)-1 {
					m.selected++
				}
			case "enter":
				m.state = "game"
				words, err := LoadWords("words.txt")
				if err != nil {
					m.state = "error"
					return m, tea.Quit
				}
				var filteredWords []string
				for _, word := range words {
					if len(word) >= m.difficulties[m.selected].WordLengthMin &&
						len(word) <= m.difficulties[m.selected].WordLengthMax {
						filteredWords = append(filteredWords, word)
					}
				}
				m.game = &Game{
					Difficulty: m.difficulties[m.selected],
					Words:      filteredWords,
					TimeLimit:  m.difficulties[m.selected].TimeLimit,
					StartTime: time.Now(),
				}
				m.timer = time.NewTimer(m.game.TimeLimit)
				m.index = 0
				m.currentWord = m.game.Words[m.index]
			case "q":
				return m, tea.Quit
			}
		case "game":
			switch msg.String() {
			case "enter":
				if strings.TrimSpace(m.input) == m.currentWord {
					m.game.Score++
				} else {
					m.game.Mistakes++
				}
				m.input = ""
				m.index++
				if m.index >= len(m.game.Words) {
					m.state = "end"
				} else {
					m.currentWord = m.game.Words[m.index]
				}
			case "q":
				return m, tea.Quit
			default:
				m.input += msg.String()
				m.game.CharCount += len(msg.String())
			}
		case "end":
			if msg.String() == "q" {
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// View renders the UI
func (m Model) View() string {
	switch m.state {
	case "menu":
		s := "難易度選択:\n\n"
		for i, d := range m.difficulties {
			cursor := " "
			if m.selected == i {
				cursor = ">"
			}
			s += fmt.Sprintf("%s %s - %s\n", cursor, d.Name, d.Description)
		}
		s += "\n上下キーで選択、Enterで決定、qで終了"
		return s
	case "game":
		return fmt.Sprintf("問題: %s\n入力: %s\nスコア: %d ミスタイプ: %d\nqで終了",
			m.currentWord, m.input, m.game.Score, m.game.Mistakes)
	case "end":
		return fmt.Sprintf("ゲーム終了!\nスコア: %d\nミスタイプ: %d\nTPS: %.2f\nqで終了",
			m.game.Score, m.game.Mistakes, float64(m.game.CharCount)/time.Since(m.game.StartTime).Seconds())
	case "error":
		return "エラー: 単語ファイルの読み込みに失敗しました。\nqで終了"
	}
	return ""
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("エラー: %v\n", err)
		os.Exit(1)
	}
}
