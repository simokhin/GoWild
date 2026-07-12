package main

import "fmt"

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
