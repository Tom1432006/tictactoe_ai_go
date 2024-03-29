package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var movesMade uint8
var ai_player uint8
var field [9]uint8

// var db, _ = gorm.Open("mysql", "root:@tcp(127.0.0.1:3306)/todolist-go")

func Healtz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": "true"}`)
}

func init() {
	// setup logger
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}

func main() {
	// defer db.Close()

	log.Info("Starting TicTacToe API server")
	router := mux.NewRouter()
	router.HandleFunc("/healtz", Healtz).Methods("GET")
	router.HandleFunc("/ai/{board}", GetBestMove).Methods("GET")

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET"},
	}).Handler(router)

	http.ListenAndServe(":8000", handler)
}

func GetBestMove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	board, _ := vars["board"]

	// convert the board string to an array
	board_split := strings.Split(board, "")

	if len(board_split) != 9 {
		// throw error the length of the board string is not 9
		log.WithFields(log.Fields{"board": board}).Info("Error code 100. Not enough data.")
		io.WriteString(w, `{"success": "false", "error": "The board configuration is not a valid configuration. See Docs for more info!", "code": "100"}`)
		return
	}

	// find out which players turn it is
	var num_ones uint8
	var num_twos uint8
	var num_zeros uint8

	for i := 0; i < 9; i++ {
		n, _ := strconv.Atoi(board_split[i])
		field[i] = uint8(n)

		if n == 1 {
			num_ones++
		} else if n == 2 {
			num_twos++
		} else if n == 0 {
			num_zeros++
		} else {
			// throw error if n != 1, 2 or 0
			log.WithFields(log.Fields{"board": board}).Info("Error code 101. Number is not a 0, 1, 2")
			io.WriteString(w, `{"success": "false", "error": "The board configuration is not a valid configuration. See Docs for more info!", "code": "101"}`)
			return
		}
	}

	if num_zeros == 0 {
		// throw error because there has to be at least one zero
		log.WithFields(log.Fields{"board": board}).Info("Error code 102. No possible moves")
		io.WriteString(w, `{"success": "false", "error": "The board configuration is not a valid configuration. See Docs for more info!", "code": "102"}`)
		return
	}

	movesMade = 9 - num_zeros

	if num_ones == num_twos {
		ai_player = 1
	} else if num_ones == num_twos+1 {
		ai_player = 2
	} else {
		// throw an error if the ratio of the ones and twos is not correct
		log.WithFields(log.Fields{"board": board}).Info("Error code 103. Wrong ratio of 1s and 2s!")
		io.WriteString(w, `{"success": "false", "error": "The board configuration is not a valid configuration. See Docs for more info!", "code": "103"}`)
		return
	}

	move := findBestMove()
	log.WithFields(log.Fields{"board": board, "move": move, "player": ai_player}).Info("Evaluated board position and send response")
	w.Header().Set("Content-Type", "application/json")
	p := map[string]interface{}{
		"success": true,
		"move":    move,
		"player":  ai_player,
	}
	json.NewEncoder(w).Encode(p)
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
