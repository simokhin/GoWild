package main

import "fmt"

// ResetBoard clears a Board to its initial empty state: all squares set to
// OffBoard (border) or Empty (inner 64), piece counts zeroed, and history reset.
func ResetBoard(pos *Board) {
	for index := range 120 {
		pos.Pieces[index] = OffBoard
	}

	for index := range 64 {
		pos.Pieces[SQ120(index)] = Empty
	}

	for index := range 2 {
		pos.BigPce[index] = 0
		pos.MajPce[index] = 0
		pos.MinPce[index] = 0
		pos.Material[index] = 0
	}

	for index := range 3 {
		pos.Pawns[index] = 0
	}

	for index := range 13 {
		pos.PceNum[index] = 0
	}

	pos.KingSq[White] = NoSquare
	pos.KingSq[Black] = NoSquare

	pos.Side = Both
	pos.EnPas = NoSquare
	pos.FiftyMove = 0

	pos.Ply = 0

	pos.CastlePerm = 0

	pos.PosKey = 0

	pos.History = pos.History[:0]
}

// ParseFEN parses a Forsyth–Edwards Notation string into a Board struct.
// It clears the board, places pieces according to the FEN, and sets the
// side to move, castling rights, en passant square, and computes the Zobrist key.
// Returns 0 on success, -1 on parse error.
func ParseFEN(fen string, pos *Board) int {
	Assert(pos != nil, "pos must not be nil")

	rank := Rank8
	file := FileA
	var piece Piece = 0
	count := 0
	i := 0

	ResetBoard(pos)

	for rank >= Rank1 && i < len(fen) {
		count = 1

		switch fen[i] {
		case 'p':
			piece = BP
		case 'r':
			piece = BR
		case 'n':
			piece = BN
		case 'b':
			piece = BB
		case 'k':
			piece = BK
		case 'q':
			piece = BQ
		case 'P':
			piece = WP
		case 'R':
			piece = WR
		case 'N':
			piece = WN
		case 'B':
			piece = WB
		case 'K':
			piece = WK
		case 'Q':
			piece = WQ
		case '1', '2', '3', '4', '5', '6', '7', '8':
			piece = Empty
			count = int(fen[i] - '0')
		case '/', ' ':
			rank--
			file = FileA
			i++
			continue
		default:
			fmt.Println("FEN error")
			return -1
		}

		for j := 0; j < count; j++ {
			sq64 := int(rank)*8 + int(file)
			sq120 := SQ120(sq64)
			if piece != Empty {
				pos.Pieces[sq120] = piece
			}
			file++
		}
		i++
	}

	Assert(fen[i] == 'w' || fen[i] == 'b', "expected side to move")
	if fen[i] == 'w' {
		pos.Side = White
	} else {
		pos.Side = Black
	}
	i += 2

	for range 4 {
		if fen[i] == ' ' {
			break
		}
		switch fen[i] {
		case 'K':
			pos.CastlePerm |= WKCA
		case 'Q':
			pos.CastlePerm |= WQCA
		case 'k':
			pos.CastlePerm |= BKCA
		case 'q':
			pos.CastlePerm |= BQCA
		}
		i++
	}
	i++

	Assert(pos.CastlePerm >= 0 && pos.CastlePerm <= 15, "castle perm out of range")

	if fen[i] != '-' {
		file = File(fen[i] - 'a')
		rank = Rank(fen[i+1] - '1')
		Assert(file >= FileA && file <= FileH, "file out of range")
		Assert(rank >= Rank1 && rank <= Rank8, "rank out of range")
		pos.EnPas = FR2SQ(file, rank)
	}

	pos.PosKey = GeneratePosKey(pos)

	UpdateListsMaterial(pos)

	return 0
}

// PrintBoard prints a visual representation of the board state to stdout.
// It displays pieces by their character codes, along with side to move,
// en passant square, castling rights, and the Zobrist position key.
func PrintBoard(pos *Board) {
	var sq Square
	var file File
	var rank Rank
	var piece Piece

	fmt.Println("\nGameBoard:")

	for rank = Rank8; rank >= Rank1; rank-- {
		fmt.Printf("%d ", int(rank)+1)
		for file = FileA; file <= FileH; file++ {
			sq = FR2SQ(file, rank)
			piece = pos.Pieces[sq]
			fmt.Printf("%3c", PceChar[piece])
		}
		fmt.Println()
	}

	fmt.Print("\n   ")
	for file = FileA; file <= FileH; file++ {
		fmt.Printf("%3c", FileChar[file])
	}
	fmt.Println()

	fmt.Printf("side:%c\n", SideChar[pos.Side])
	fmt.Printf("enPas:%d\n", pos.EnPas)

	castleWK := byte('-')
	if pos.CastlePerm&WKCA != 0 {
		castleWK = 'K'
	}
	castleWQ := byte('-')
	if pos.CastlePerm&WQCA != 0 {
		castleWQ = 'Q'
	}
	castleBK := byte('-')
	if pos.CastlePerm&BKCA != 0 {
		castleBK = 'k'
	}
	castleBQ := byte('-')
	if pos.CastlePerm&BQCA != 0 {
		castleBQ = 'q'
	}
	fmt.Printf("castle:%c%c%c%c\n", castleWK, castleWQ, castleBK, castleBQ)

	fmt.Printf("PosKey:%X\n", pos.PosKey)
}

// UpdateListsMaterial walks the entire 120-square board and recomputes the
// per-side piece lists (PList), piece counts (PceNum, BigPce, MajPce, MinPce),
// material totals (Material), and king squares (KingSq).
// Must be called after loading a position (e.g., after ParseFEN).
func UpdateListsMaterial(pos *Board) {
	for index := range 120 {
		sq := index
		piece := pos.Pieces[sq]

		if piece != OffBoard && piece != Empty {
			colour := PieceCol[piece]

			if PieceBig[piece] {
				pos.BigPce[colour]++
			}
			if PieceMin[piece] {
				pos.MinPce[colour]++
			}
			if PieceMaj[piece] {
				pos.MajPce[colour]++
			}

			pos.Material[colour] += PieceVal[piece]

			pos.PList[piece][pos.PceNum[piece]] = Square(sq)
			pos.PceNum[piece]++

			if piece == WK {
				pos.KingSq[White] = Square(sq)
			}
			if piece == BK {
				pos.KingSq[Black] = Square(sq)
			}
		}
	}
}
