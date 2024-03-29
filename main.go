package main

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

	fmt.Println("Let's play TicTacToe without the board! Good luck trying to win!")

	// let the human choose, which player to be
	for {
		fmt.Printf("Which player do you want to be? 1 = Player1; 2 = Player2: ")
		player, _ := reader.ReadString('\n')
		// convert string player to a playerIndex
		playerIndex, err := strconv.Atoi(strings.Split(player, "\r\n")[0])
		if err != nil || playerIndex < 1 || playerIndex > 2 {
			fmt.Println("Please enter a number from 1 to 2")
			continue
		}

		ai_player = uint8(playerIndex)%2 + 1
		break
	}

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

		// if no player has won and all 9 moves have been made, the game ends in a draw
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

// Solve TicTacToe using the minimax algorithm
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
			continue // skip non empty cells
		}

		var newturn uint8 = turn%2 + 1

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
	if playerWon(current_field, ai_player) {
		// Computer wins
		return 10
	}
	if playerWon(current_field, ai_player%2+1) {
		// Human wins
		return -10
	}
	return 0
}
