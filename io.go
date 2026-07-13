package main

import "fmt"

// NoMove is the sentinel value returned by ParseMove when the input string does
// not represent a valid move in the current position.
const NoMove = 0

// ParseMove parses a user-supplied move string in algebraic notation (e.g., "e2e4")
// and returns the corresponding encoded move integer. Promotion moves may include
// a fifth character ('n', 'b', 'r', or 'q') to specify the promoted piece. Returns
// NoMove if the input is invalid or the move is not found in the generated move list.
func ParseMove(ptrChar string, pos *Board) int {
	Assert(CheckBoard(pos), "board check failed")

	// Validate that file/rank characters are within bounds (a1–h8)
	if ptrChar[1] > '8' || ptrChar[1] < '1' {
		return NoMove
	}
	if ptrChar[3] > '8' || ptrChar[3] < '1' {
		return NoMove
	}
	if ptrChar[0] > 'h' || ptrChar[0] < 'a' {
		return NoMove
	}
	if ptrChar[2] > 'h' || ptrChar[2] < 'a' {
		return NoMove
	}

	from := FR2SQ(File(ptrChar[0]-'a'), Rank(ptrChar[1]-'1'))
	to := FR2SQ(File(ptrChar[2]-'a'), Rank(ptrChar[3]-'1'))

	fmt.Printf("%d %d %s\n", from, to, ptrChar)

	Assert(SqOnBoard(from) && SqOnBoard(to), "square not on board")

	list := &MoveList{}
	GenerateAllMoves(pos, list)

	var move int
	var promPce Piece

	// Iterate over all pseudo-legal moves and try to match the user input
	for moveNum := 0; moveNum < list.Count; moveNum++ {
		move = list.Moves[moveNum].MoveInt
		if FromSq(move) == int(from) && ToSq(move) == int(to) {
			promPce = Promoted(move)
			if promPce != Empty {
				// Promotion move: require a 5th character to disambiguate
				if len(ptrChar) < 5 {
					continue
				}
				// Match the promoted piece type against the user-supplied character
				if IsRQ(promPce) && !IsBQ(promPce) && ptrChar[4] == 'r' {
					return move
				} else if !IsRQ(promPce) && IsBQ(promPce) && ptrChar[4] == 'b' {
					return move
				} else if IsRQ(promPce) && IsBQ(promPce) && ptrChar[4] == 'q' {
					return move
				} else if IsKn(promPce) && ptrChar[4] == 'n' {
					return move
				}
				continue
			}
			// Non-promotion move: direct match
			return move
		}
	}
	return NoMove
}

// PrSq converts a 120-square mailbox index into human-readable algebraic notation
// (e.g., A1 → "a1", H8 → "h8"). Uses the precomputed FilesBrd and RanksBrd tables.
func PrSq(sq Square) string {
	file := FilesBrd[sq]
	rank := RanksBrd[sq]
	return fmt.Sprintf("%c%c", 'a'+byte(file), '1'+byte(rank))
}

// PrMove converts an encoded move integer into standard algebraic notation
// (e.g., from A2 to H7 → "a2h7"). For promotion moves, the promoted piece
// is appended as a lowercase character ('n', 'b', 'r', or 'q').
func PrMove(move int) string {
	from := Square(FromSq(move))
	to := Square(ToSq(move))
	promoted := Promoted(move)

	moveStr := PrSq(from) + PrSq(to)

	if promoted != Empty {
		pchar := byte('q')
		if IsKn(promoted) {
			pchar = 'n'
		} else if IsRQ(promoted) && !IsBQ(promoted) {
			pchar = 'r'
		} else if !IsRQ(promoted) && IsBQ(promoted) {
			pchar = 'b'
		}
		moveStr += string(pchar)
	}

	return moveStr
}

// PrintMoveList prints all moves in a MoveList to stdout, one per line.
// Each line shows the move index, algebraic notation (via PrMove), and the
// search score. The total move count is printed at the end.
func PrintMoveList(list *MoveList) {
	fmt.Println("MoveList:")

	for index := 0; index < list.Count; index++ {
		move := list.Moves[index].MoveInt
		score := list.Moves[index].Score

		fmt.Printf("Move:%d > %s (score:%d)\n", index+1, PrMove(move), score)
	}
	fmt.Printf("MoveList Total %d Moves:\n\n", list.Count)
}
