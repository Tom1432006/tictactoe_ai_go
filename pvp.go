package pvp

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

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Let's play TicTacToe!")

	printBoard()

	for {
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
