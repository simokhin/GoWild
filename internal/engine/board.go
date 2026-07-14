package engine

import "fmt"

// CheckBoard verifies the integrity of a Board struct. It cross-checks the piece lists
// (PList) against the mailbox board (Pieces), counts all pieces by type/colour/material,
// validates pawn bitboards, and checks the Zobrist position key, en passant square,
// king squares, and side to move. Returns true on success; panics on any mismatch
// (when Debug is enabled).
func CheckBoard(pos *Board) bool {
	var tPceNum [13]int
	var tBigPce [2]int
	var tMajPce [2]int
	var tMinPce [2]int
	var tMaterial [2]int

	var sq64 int
	var tPiece Piece
	var sq120 Square
	var colour Color
	var pcount int

	tPawns := [3]Bitboard{0, 0, 0}

	tPawns[White] = pos.Pawns[White]
	tPawns[Black] = pos.Pawns[Black]
	tPawns[Both] = pos.Pawns[Both]

	for tPiece := WP; tPiece <= BK; tPiece++ {
		for tPceNum := 0; tPceNum < pos.PceNum[tPiece]; tPceNum++ {
			sq120 = pos.PList[tPiece][tPceNum]
			Assert(pos.Pieces[sq120] == tPiece, "piece list inconsistent with board")
		}
	}

	for sq64 = range 64 {
		sq120 = SQ120(sq64)
		tPiece = pos.Pieces[sq120]
		if tPiece == Empty {
			continue
		}
		tPceNum[tPiece]++
		colour = PieceCol[tPiece]
		if PieceBig[tPiece] {
			tBigPce[colour]++
		}
		if PieceMin[tPiece] {
			tMinPce[colour]++
		}
		if PieceMaj[tPiece] {
			tMajPce[colour]++
		}
		tMaterial[colour] += PieceVal[tPiece]
	}

	for tPiece = WP; tPiece <= BK; tPiece++ {
		Assert(tPceNum[tPiece] == pos.PceNum[tPiece], "piece count mismatch")
	}

	// Verify that pawn bitboard population counts match the piece-list counts
	pcount = CNT(tPawns[White])
	Assert(pcount == pos.PceNum[WP], "white pawn bitboard count mismatch")

	pcount = CNT(tPawns[Black])
	Assert(pcount == pos.PceNum[BP], "black pawn bitboard count mismatch")

	pcount = CNT(tPawns[Both])
	Assert(pcount == pos.PceNum[BP]+pos.PceNum[WP], "combined pawn bitboard count mismatch")

	// Verify that each set bit in the White pawn bitboard actually holds a White pawn
	for tPawns[White] != 0 {
		sq64 = POP(&tPawns[White])
		Assert(pos.Pieces[SQ120(sq64)] == WP, "white pawn bitboard square mismatch")
	}

	// Verify that each set bit in the Black pawn bitboard actually holds a Black pawn
	for tPawns[Black] != 0 {
		sq64 = POP(&tPawns[Black])
		Assert(pos.Pieces[SQ120(sq64)] == BP, "black pawn bitboard square mismatch")
	}

	// Verify that each set bit in the combined pawn bitboard actually holds a pawn of either colour
	for tPawns[Both] != 0 {
		sq64 = POP(&tPawns[Both])
		Assert(pos.Pieces[SQ120(sq64)] == BP || pos.Pieces[SQ120(sq64)] == WP, "combined pawn bitboard square mismatch")
	}

	Assert(tMaterial[White] == pos.Material[White] && tMaterial[Black] == pos.Material[Black], "material count mismatch")
	Assert(tMinPce[White] == pos.MinPce[White] && tMinPce[Black] == pos.MinPce[Black], "minor count mismatch")
	Assert(tMajPce[White] == pos.MajPce[White] && tMajPce[Black] == pos.MajPce[Black], "major count mismatch")
	Assert(tBigPce[White] == pos.BigPce[White] && tBigPce[Black] == pos.BigPce[Black], "big count mismatch")

	Assert(pos.Side == White || pos.Side == Black, "side must be white or black")
	Assert(GeneratePosKey(pos) == pos.PosKey, "position key mismatch")

	Assert(pos.EnPas == NoSquare ||
		(RanksBrd[pos.EnPas] == Rank6 && pos.Side == White) ||
		(RanksBrd[pos.EnPas] == Rank3 && pos.Side == Black), "en passant square inconsistent with side to move")

	Assert(pos.Pieces[pos.KingSq[White]] == WK, "white king square mismatch")
	Assert(pos.Pieces[pos.KingSq[Black]] == BK, "black king square mismatch")

	return true
}

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

	// Lazily allocate the PV table on first reset so callers don't have to set
	// it up themselves, then clear it since any previous position's PV entries
	// are no longer valid.
	if pos.HashTable == nil {
		pos.HashTable = &HashTable{}
	}
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
		sq := Square(index)
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

			switch piece {
			case WP:
				SETBIT(&pos.Pawns[White], SQ64(sq))
				SETBIT(&pos.Pawns[Both], SQ64(sq))
			case BP:
				SETBIT(&pos.Pawns[Black], SQ64(sq))
				SETBIT(&pos.Pawns[Both], SQ64(sq))
			}
		}
	}
}

func MirrorBoard(pos *Board) {
	var tempPiecesArray [64]Piece
	tempSide := pos.Side ^ 1

	swapPieces := [13]Piece{Empty, BP, BN, BB, BR, BQ, BK, WP, WN, WB, WR, WQ, WK}

	var tempCastlePerm CastlePerm
	tempEnPas := NoSquare

	if pos.CastlePerm&WKCA != 0 {
		tempCastlePerm |= BKCA
	}
	if pos.CastlePerm&WQCA != 0 {
		tempCastlePerm |= BQCA
	}
	if pos.CastlePerm&BKCA != 0 {
		tempCastlePerm |= WKCA
	}
	if pos.CastlePerm&BQCA != 0 {
		tempCastlePerm |= WQCA
	}

	if pos.EnPas != NoSquare {
		tempEnPas = SQ120(Mirror64[SQ64(pos.EnPas)])
	}

	for sq := 0; sq < 64; sq++ {
		tempPiecesArray[sq] = pos.Pieces[SQ120(Mirror64[sq])]
	}

	ResetBoard(pos)

	for sq := 0; sq < 64; sq++ {
		tp := swapPieces[tempPiecesArray[sq]]
		pos.Pieces[SQ120(sq)] = tp
	}

	pos.Side = tempSide
	pos.CastlePerm = tempCastlePerm
	pos.EnPas = tempEnPas
	pos.PosKey = GeneratePosKey(pos)

	UpdateListsMaterial(pos)

	Assert(CheckBoard(pos), "board check failed")
}
