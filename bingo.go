package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Package variable for storing the global state of the game.
var bingo struct {
	phrases []string
	players []player
	room    string
	topic   string
}

type bingoCell struct {
	Phrase string `json:"phrase"`
	Marked bool   `json:"marked"`
}

type player struct {
	Username   string `json:"username"`
	uid        string
	BingoBoard []bingoCell `json:"bingo_board"`
}

type updateBingoCellReturn struct {
	Error  string `json:"error,omitempty"`
	Marked bool   `json:"marked"`
}

func (ret updateBingoCellReturn) String() string {
	retJSON, err := json.Marshal(ret)
	if err == nil {
		return string(retJSON)
	}
	return "{}"
}

// Update the marked state of a bingo cell.
func updateBingoCell(w http.ResponseWriter, r *http.Request) {
	ret := updateBingoCellReturn{}

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		ret.Error = "Missing UID"
		fmt.Fprint(w, ret)
		return
	}

	cell := r.URL.Query().Get("cell")
	if cell == "" {
		ret.Error = "Missing cell"
		fmt.Fprint(w, ret)
		return
	}
	cellInt, err := strconv.Atoi(cell)
	if err != nil {
		ret.Error = err.Error()
		fmt.Fprint(w, ret)
		return
	}

	marked := r.URL.Query().Get("marked")
	if marked == "" {
		ret.Error = "Missing marked"
		fmt.Fprint(w, ret)
		return
	}
	markedBool, err := strconv.ParseBool(marked)
	if err != nil {
		ret.Error = err.Error()
		fmt.Fprint(w, ret)
		return
	}

	// Find the player with the given UID.
	playerIndex := -1
	for i, player := range bingo.players {
		if player.uid == uid {
			playerIndex = i
		}
	}
	if playerIndex < 0 {
		ret.Error = "Invalid UID"
		fmt.Fprint(w, ret)
		return
	}

	// Make sure the cell is valid.
	if cellInt < 0 || cellInt+1 > len(bingo.players[playerIndex].BingoBoard) {
		ret.Error = "Cell out of bounds"
		fmt.Fprint(w, ret)
		return
	}

	bingo.players[playerIndex].BingoBoard[cellInt].Marked = markedBool
	ret.Marked = markedBool
	fmt.Fprint(w, ret)
}

type getGameDataReturn struct {
	Error   string   `json:"error,omitempty"`
	Topic   string   `json:"topic"`
	Players []player `json:"players"`
}

func (ret getGameDataReturn) String() string {
	retJSON, err := json.Marshal(ret)
	if err == nil {
		return string(retJSON)
	}
	return "{}"
}

// Return all the game data needed for rendering.
func getGameData(w http.ResponseWriter, r *http.Request) {
	ret := getGameDataReturn{}

	uid := r.URL.Query().Get("uid")
	if uid == "" {
		ret.Error = "Missing UID"
		fmt.Fprint(w, ret)
		return
	}

	// Make sure thet player supplied a valid UID, and if so return the game
	// data.
	for _, player := range bingo.players {
		if player.uid == uid {
			ret.Players = bingo.players
			ret.Topic = bingo.topic
			fmt.Fprint(w, ret)
			return
		}
	}

	ret.Error = "Invalid UID"
	fmt.Fprint(w, ret)
}

type newPlayerReturn struct {
	Error string `json:"error,omitempty"`
	UID   string `json:"uid,omitempty"`
}

func (ret newPlayerReturn) String() string {
	retJSON, err := json.Marshal(ret)
	if err == nil {
		return string(retJSON)
	}
	return "{}"
}

// Add a new player to the game given their username, and return a UID to them.
func newPlayer(w http.ResponseWriter, r *http.Request) {
	ret := newPlayerReturn{}
	player := player{}

	// Make sure the room code is valid.
	room := r.URL.Query().Get("room")
	if room != bingo.room {
		ret.Error = "Invalid room code"
		fmt.Fprint(w, ret)
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		ret.Error = "Missing username"
		fmt.Fprint(w, ret)
		return
	}

	// Make sure the username is unique among current players.
	for _, player := range bingo.players {
		if player.Username == username {
			ret.Error = "Username already taken"
			fmt.Fprint(w, ret)
			return
		}
	}

	// Build the player object and add it to our players array.
	player.Username = username
	player.uid = strconv.Itoa(rand.Int())
	player.BingoBoard = make([]bingoCell, len(bingo.phrases))
	for i, phrase := range shufflePhrases(bingo.phrases) {
		player.BingoBoard[i].Phrase = phrase
	}
	bingo.players = append(bingo.players, player)

	ret.UID = player.uid
	fmt.Fprint(w, ret)
	fmt.Printf("Player %s joined with UID %s\n", player.Username, player.uid)
}

// Takes an array of phrases and returns a new array with those phrases randomly
// shuffled.
func shufflePhrases(phrases []string) []string {
	numPhrases := len(phrases)
	shuffledPhrases := make([]string, numPhrases)
	for i, j := range rand.Perm(numPhrases) {
		shuffledPhrases[j] = phrases[i]
	}
	return shuffledPhrases
}

// Loads phrases from a text file with each phrase on its own line.
func loadPhrases(phrasesFile string) ([]string, error) {
	file, err := os.Open(phrasesFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	// Seed the RNG.
	rand.Seed(time.Now().UnixNano())

	// Parse command-line arguments.
	port := flag.String("port", "8080", "the port to listen on")
	help := flag.Bool("help", false, "print usage")
	phrasesFile := flag.String("phrases", "", "the phrases file to use (required)")
	htmlRootDir := flag.String("html", "./html", "path to the html directory for the game")
	room := flag.String("room", "", "the room code players will need to join this game (required)")
	topic := flag.String("topic", "Generic Bingo", "the topic for the game")
	flag.Parse()
	if *help || *phrasesFile == "" || *room == "" {
		flag.Usage()
		return
	}

	bingo.room = *room
	bingo.topic = *topic

	// Load phrases.
	phrases, err := loadPhrases(*phrasesFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	bingo.phrases = phrases

	fmt.Println("Starting up Multiplayer Bingo...")
	fmt.Printf("Room Code: %s\n", bingo.room)
	fmt.Printf("Topic: %s\n", bingo.topic)

	// Handle HTTP requests.
	htmlRoot := http.FileServer(http.Dir(*htmlRootDir))
	http.Handle("/", htmlRoot)
	http.HandleFunc("/join", newPlayer)
	http.HandleFunc("/game", getGameData)
	http.HandleFunc("/cell", updateBingoCell)
	http.ListenAndServe(":"+*port, nil)
}
