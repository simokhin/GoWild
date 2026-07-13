package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseGo(line string, info *SearchInfo, pos *Board) {
	depth := -1
	movesToGo := 30
	moveTime := -1
	time := -1
	inc := 0

	info.TimeSet = false

	tokens := strings.Fields(line)

	for i, tok := range tokens {
		switch tok {
		case "infinite":
		case "binc":
			if pos.Side == Black && i+1 < len(tokens) {
				inc, _ = strconv.Atoi(tokens[i+1])
			}
		case "winc":
			if pos.Side == White && i+1 < len(tokens) {
				inc, _ = strconv.Atoi(tokens[i+1])
			}
		case "wtime":
			if pos.Side == White && i+1 < len(tokens) {
				time, _ = strconv.Atoi(tokens[i+1])
			}
		case "btime":
			if pos.Side == Black && i+1 < len(tokens) {
				time, _ = strconv.Atoi(tokens[i+1])
			}
		case "movestogo":
			if i+1 < len(tokens) {
				movesToGo, _ = strconv.Atoi(tokens[i+1])
			}
		case "movetime":
			if i+1 < len(tokens) {
				moveTime, _ = strconv.Atoi(tokens[i+1])
			}
		case "depth":
			if i+1 < len(tokens) {
				depth, _ = strconv.Atoi(tokens[i+1])
			}
		}
	}

	if moveTime != -1 {
		time = moveTime
		movesToGo = 1
	}

	info.StartTime = GetTimeMs()
	info.Depth = depth

	if time != -1 {
		info.TimeSet = true
		time /= movesToGo
		time -= 50
		info.StopTime = info.StartTime + int64(time) + int64(inc)
	}

	if depth == -1 {
		info.Depth = MaxDepth
	}

	go SearchPosition(pos, info)
}

func ParsePosition(lineIn string, pos *Board) {
	tokens := strings.Fields(lineIn)

	if len(tokens) < 2 {
		return
	}

	index := 1 // tokens[0] == "position"

	if tokens[index] == "startpos" {
		ParseFEN(START_FEN, pos)
		index++
	} else if tokens[index] == "fen" {
		index++
		if index+6 > len(tokens) {
			ParseFEN(START_FEN, pos)
			return
		}
		fenTokens := tokens[index : index+6]
		fen := strings.Join(fenTokens, " ")
		ParseFEN(fen, pos)
		index += 6
	} else {
		ParseFEN(START_FEN, pos)
	}

	if index < len(tokens) && tokens[index] == "moves" {
		index++
		for ; index < len(tokens); index++ {
			move := ParseMove(tokens[index], pos)
			if move == NoMove {
				break
			}
			MakeMove(pos, move)
			pos.Ply = 0
		}
	}

	PrintBoard(pos)
}

func UciLoop() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("id name %s\n", Name)
	fmt.Println("id author Nikita Simokhin")
	fmt.Println("uciok")

	pos := &Board{}
	pos.PvTable = &PVTable{}
	info := &SearchInfo{}

	InitPvTable(pos.PvTable)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "isready"):
			fmt.Println("readyok")
		case strings.HasPrefix(line, "position"):
			ParsePosition(line, pos)
		case strings.HasPrefix(line, "ucinewgame"):
			ParsePosition("position startpos\n", pos)
		case strings.HasPrefix(line, "go"):
			ParseGo(line, info, pos)
		case strings.HasPrefix(line, "quit"):
			info.Quit = true
			return
		case strings.HasPrefix(line, "uci"):
			fmt.Printf("id name %s\n", Name)
			fmt.Println("id author Nikita Simokhin")
			fmt.Println("uciok")
		case strings.HasPrefix(line, "stop"):
			info.Stopped.Store(true)
		default:
			if info.Quit {
				return
			}
		}
	}
}
