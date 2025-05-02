package main

import (
	"bufio"
	"fmt"
	"os"
	// "strings"
	"time"
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

// StartGame starts the typing game
func StartGame(game *Game) {
	fmt.Printf("ゲーム開始! 難易度: %s\n", game.Difficulty.Name)
	fmt.Printf("制限時間: %v秒\n", game.Difficulty.TimeLimit.Seconds())
	fmt.Println("問題が表示されます。入力してください。")

	timer := time.NewTimer(game.Difficulty.TimeLimit)
	input := bufio.NewScanner(os.Stdin)

	for _, word := range game.Words {
		fmt.Printf("問題: %s\n", word)
		fmt.Print("入力: ")

		done := make(chan bool)
		go func() {
			if input.Scan() {
				if input.Text() == word {
					game.Score++
				} else {
					game.Mistakes++
				}
			}
			done <- true
		}()

		select {
		case <-timer.C:
			fmt.Println("\n時間切れ!")
			return
		case <-done:
			continue
		}
	}

	fmt.Println("ゲーム終了!")
	fmt.Printf("スコア: %d\n", game.Score)
	fmt.Printf("ミスタイプ数: %d\n", game.Mistakes)
}

func main() {
	// Display difficulty options
	fmt.Println("難易度選択:")
	for i, d := range difficulties {
		fmt.Printf("%d: %s - %s\n", i+1, d.Name, d.Description)
	}

	// Select difficulty
	var choice int
	fmt.Print("難易度を選択してください (1-3): ")
	fmt.Scan(&choice)
	if choice < 1 || choice > len(difficulties) {
		fmt.Println("無効な選択です。終了します。")
		return
	}
	selectedDifficulty := difficulties[choice-1]

	// Load words
	words, err := LoadWords("words.txt")
	if err != nil {
		fmt.Printf("単語ファイルの読み込みに失敗しました: %v\n", err)
		return
	}
	
	// Filter words based on difficulty
	var filteredWords []string
	for _, word := range words {
		if len(word) >= selectedDifficulty.WordLengthMin && len(word) <= selectedDifficulty.WordLengthMax {
			filteredWords = append(filteredWords, word)
		}
	}

	// Initialize game
	game := Game{
		Difficulty: selectedDifficulty,
		Words:      filteredWords,
		TimeLimit:  selectedDifficulty.TimeLimit,
	}

	// Start the game
	StartGame(&game)
}
