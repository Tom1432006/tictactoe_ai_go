package only_wins

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var field [9]uint8
var turn uint8 = 1 // either 1 or 2
var movesMade uint8 = 0
var ai_player uint8 = 1 // either 1 or 2

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Let's play TicTacToe! This is the you cannot loose version!")

	printBoard()

	for {
		if turn == ai_player {
			// let the ai player make a move
			move := findBestMove()
			takeTurn(move, turn)

			fmt.Printf("AI Player took turn %v\n", move+1)
		} else {
			// let the player make a move
			// first get the input from the player
			fmt.Printf("Player %v: ", turn)
			playerInput, _ := reader.ReadString('\n')

			// if the input is exit, exit the loop
			if strings.TrimRight(playerInput, "\r\n") == "exit" {
				os.Exit(200)
			}

			// convert string playerInput to an index
			index, err := strconv.Atoi(strings.Split(playerInput, "\r\n")[0])
			if err != nil {
				fmt.Println("Please enter a number from 1 to 9")
				continue
			}

			err = takeTurn(index-1, turn)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		printBoard()

		movesMade++

		if playerWon(field, uint8(turn)) {
			fmt.Printf("Player %v has won\n", turn)
			break
		}

		if movesMade == 9 {
			fmt.Println("The game was a draw!")
			break
		}

		turn += 1
		if turn == 3 {
			turn = 1
		}
	}

	fmt.Printf("Press enter to exit:")
	reader.ReadString('\n')
}

func takeTurn(index int, turn uint8) error {
	if index < 0 || index > 8 || field[index] != 0 {
		return errors.New("This field is not free. Try again.")
	}

	field[index] = turn
	return nil
}

func printBoard() {
	for i := 0; i < 9; i++ {
		if i%3 == 0 && i != 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("%v ", field[i])
	}
	fmt.Printf("\n")
}

/**
check if any of the two players have won
*/
func playerWon(field [9]uint8, p uint8) bool {
	if (field[0] == p && field[1] == p && field[2] == p) ||
		(field[3] == p && field[4] == p && field[5] == p) ||
		(field[6] == p && field[7] == p && field[8] == p) ||
		(field[0] == p && field[3] == p && field[6] == p) ||
		(field[1] == p && field[4] == p && field[7] == p) ||
		(field[2] == p && field[5] == p && field[8] == p) ||
		(field[0] == p && field[4] == p && field[8] == p) ||
		(field[2] == p && field[4] == p && field[6] == p) {
		return true
	}
	return false
}

func findBestMove() int {
	field_copy := field
	depth := 9 - movesMade
	bestMove := minimax(field_copy, depth, ai_player)

	return bestMove[0]
}

func minimax(field_copy [9]uint8, depth uint8, turn uint8) [2]int {
	// base state
	if depth == 0 || playerWon(field_copy, 2) || playerWon(field_copy, 1) {
		return [2]int{-1, scoreGame(field_copy)}
	}

	var best [2]int
	if turn == ai_player {
		// Computer
		best = [2]int{-1, -99999}
	} else {
		// human
		best = [2]int{-1, 999999}
	}

	var score [2]int

	// loop through the empty cells
	for i := 0; i < 9; i++ {
		if field_copy[i] != 0 {
			continue // skip
		}

		var newturn uint8 = turn + 1
		if newturn == 3 {
			newturn = 1
		}

		field_copy[i] = turn
		score = minimax(field_copy, depth-1, newturn)
		field_copy[i] = 0
		score[0] = i

		if turn == ai_player {
			// Computer
			if score[1] > best[1] {
				best = score
			}
		} else {
			// human
			if score[1] < best[1] {
				best = score
			}
		}
	}

	return best
}

func scoreGame(current_field [9]uint8) int {
	if playerWon(current_field, 1) {
		return -10
	}
	if playerWon(current_field, 2) {
		return 10
	}
	return 0
}
