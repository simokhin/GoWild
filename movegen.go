package main

// ---------------------------------------------------------------------------
// Move generation direction tables
// ---------------------------------------------------------------------------

// LoopSlidePce is the iteration order for sliding pieces, arranged as white bishop,
// rook, queen followed by black bishop, rook, queen. Empty entries act as sentinel
// terminators (the loop stops when piece == Empty).
var LoopSlidePce = [8]Piece{WB, WR, WQ, Empty, BB, BR, BQ, Empty}

// LoopNonSlidePce is the iteration order for non-sliding pieces, arranged as white
// knight, king followed by black knight, king. Empty entries act as terminators.
var LoopNonSlidePce = [6]Piece{WN, WK, Empty, BN, BK, Empty}

// LoopSlideIndex stores the start index into LoopSlidePce for each side:
// [White]=0, [Black]=4 (skipping the white block and its Empty terminator).
var LoopSlideIndex = [2]int{0, 4}

// LoopNonSlideIndex stores the start index into LoopNonSlidePce for each side:
// [White]=0, [Black]=3 (skipping the white block and its Empty terminator).
var LoopNonSlideIndex = [2]int{0, 3}

// PceDir maps each piece type (indexed by Piece constant) to its possible move
// direction offsets on the 120-square mailbox board. The offsets represent the
// delta added to a square index to move one step in a given direction:
//
//	-1 = west, +1 = east, -10 = north, +10 = south,
//	-11 = north-west, -9 = north-east, +11 = south-east, +9 = south-west.
//
// Knights use the 8 L-shaped offsets; bishops use 4 diagonals; rooks use 4
// orthogonals; queens and kings use all 8; pawns and Empty use zero-filled rows.
var PceDir = [13][8]int{
	{0, 0, 0, 0, 0, 0, 0, 0},           // Empty
	{0, 0, 0, 0, 0, 0, 0, 0},           // WP (handled separately)
	{-8, -19, -21, -12, 8, 19, 21, 12}, // WN
	{-9, -11, 11, 9, 0, 0, 0, 0},       // WB
	{-1, -10, 1, 10, 0, 0, 0, 0},       // WR
	{-1, -10, 1, 10, -9, -11, 11, 9},   // WQ
	{-1, -10, 1, 10, -9, -11, 11, 9},   // WK
	{0, 0, 0, 0, 0, 0, 0, 0},           // BP (handled separately)
	{-8, -19, -21, -12, 8, 19, 21, 12}, // BN
	{-9, -11, 11, 9, 0, 0, 0, 0},       // BB
	{-1, -10, 1, 10, 0, 0, 0, 0},       // BR
	{-1, -10, 1, 10, -9, -11, 11, 9},   // BQ
	{-1, -10, 1, 10, -9, -11, 11, 9},   // BK
}

// NumDir stores the number of valid direction offsets for each piece type.
// Knights have 8, bishops 4, rooks 4, queens and kings 8; pawns and Empty have 0.
var NumDir = [13]int{
	0, // Empty
	0, // WP
	8, // WN
	4, // WB
	4, // WR
	8, // WQ
	8, // WK
	0, // BP
	8, // BN
	4, // BB
	4, // BR
	8, // BQ
	8, // BK
}

var VictimScore = [13]int{0, 100, 200, 300, 400, 500, 600, 100, 200, 300, 400, 500, 600}
var MvvLvaScores [13][13]int

func InitMvvLva() {
	for attacker := WP; attacker <= BK; attacker++ {
		for victim := WP; victim <= BK; victim++ {
			MvvLvaScores[victim][attacker] = VictimScore[victim] + 6 - (VictimScore[attacker] / 100)
		}
	}
}

// ---------------------------------------------------------------------------
// Move list helpers
// ---------------------------------------------------------------------------

// AddQuietMove appends a non-capture move to the move list. The score is set to
// 0 by default (move ordering is handled later by the search).
func AddQuietMove(pos *Board, move int, list *MoveList) {
	list.Moves[list.Count].MoveInt = move
	list.Moves[list.Count].Score = 0
	list.Count++
}

// AddCaptureMove appends a capture move to the move list. The score is set to
// 0 by default; the search will later assign an MVV-LVA or SEE score for ordering.
func AddCaptureMove(pos *Board, move int, list *MoveList) {
	list.Moves[list.Count].MoveInt = move
	list.Moves[list.Count].Score = MvvLvaScores[Captured(move)][pos.Pieces[FromSq(move)]]
	list.Count++
}

// AddEnPassantMove appends an en passant capture move to the move list.
// These are treated as captures but stored with a separate helper for clarity.
func AddEnPassantMove(pos *Board, move int, list *MoveList) {
	list.Moves[list.Count].MoveInt = move
	list.Moves[list.Count].Score = 105
	list.Count++
}

// MoveExist reports whether move is a legal move in the current position.
// It generates all pseudo-legal moves and makes/unmakes each one to filter
// out moves that leave the moving side's king in check.
func MoveExist(pos *Board, move int) bool {
	list := &MoveList{}
	GenerateAllMoves(pos, list)

	for moveNum := 0; moveNum < list.Count; moveNum++ {
		if !MakeMove(pos, list.Moves[moveNum].MoveInt) {
			continue
		}
		TakeMove(pos)
		if list.Moves[moveNum].MoveInt == move {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Move generation
// ---------------------------------------------------------------------------

// GenerateAllMoves generates all pseudo-legal moves for the given position and stores
// them in the move list. Generates pawn moves (single pushes, double pushes, captures,
// en passant, and promotions), sliding piece moves (bishops, rooks, queens), and
// non-sliding piece moves (knights, kings). Castling move generation will be added here.
func GenerateAllMoves(pos *Board, list *MoveList) {
	Assert(CheckBoard(pos), "board check failed")

	list.Count = 0

	side := pos.Side
	var sq, tSq Square
	var pceNum int

	dir := 0
	index := 0
	pceIndex := 0

	if side == White {
		for pceNum = 0; pceNum < pos.PceNum[WP]; pceNum++ {
			sq = pos.PList[WP][pceNum]
			Assert(SqOnBoard(sq), "square not on board")

			if pos.Pieces[sq+10] == Empty {
				AddWhitePawnMove(pos, int(sq), int(sq+10), list)
				if RanksBrd[sq] == Rank2 && pos.Pieces[sq+20] == Empty {
					AddQuietMove(pos, EncodeMove(int(sq), int(sq+20), Empty, Empty, MFlagPS), list)
				}
			}

			if !SqOffBoard(sq+9) && PieceCol[pos.Pieces[sq+9]] == Black {
				AddWhitePawnCapMove(pos, int(sq), int(sq+9), pos.Pieces[sq+9], list)
			}

			if !SqOffBoard(sq+11) && PieceCol[pos.Pieces[sq+11]] == Black {
				AddWhitePawnCapMove(pos, int(sq), int(sq+11), pos.Pieces[sq+11], list)
			}

			if pos.EnPas != NoSquare {
				if sq+9 == pos.EnPas {
					AddEnPassantMove(pos, EncodeMove(int(sq), int(sq+9), Empty, Empty, MFlagEP), list)
				}
				if sq+11 == pos.EnPas {
					AddEnPassantMove(pos, EncodeMove(int(sq), int(sq+11), Empty, Empty, MFlagEP), list)
				}
			}
		}

		// Castle

		if pos.CastlePerm&WKCA != 0 {
			if pos.Pieces[F1] == Empty && pos.Pieces[G1] == Empty {
				if !SqAttacked(E1, Black, pos) && !SqAttacked(F1, Black, pos) {
					AddQuietMove(pos, EncodeMove(int(E1), int(G1), Empty, Empty, MFlagCA), list)
				}
			}
		}

		if pos.CastlePerm&WQCA != 0 {
			if pos.Pieces[D1] == Empty && pos.Pieces[C1] == Empty && pos.Pieces[B1] == Empty {
				if !SqAttacked(E1, Black, pos) && !SqAttacked(D1, Black, pos) {
					AddQuietMove(pos, EncodeMove(int(E1), int(C1), Empty, Empty, MFlagCA), list)
				}
			}
		}
	} else {
		for pceNum = 0; pceNum < pos.PceNum[BP]; pceNum++ {
			sq = pos.PList[BP][pceNum]
			Assert(SqOnBoard(sq), "square not on board")

			if pos.Pieces[sq-10] == Empty {
				AddBlackPawnMove(pos, int(sq), int(sq-10), list)
				if RanksBrd[sq] == Rank7 && pos.Pieces[sq-20] == Empty {
					AddQuietMove(pos, EncodeMove(int(sq), int(sq-20), Empty, Empty, MFlagPS), list)
				}
			}

			if !SqOffBoard(sq-9) && PieceCol[pos.Pieces[sq-9]] == White {
				AddBlackPawnCapMove(pos, int(sq), int(sq-9), pos.Pieces[sq-9], list)
			}

			if !SqOffBoard(sq-11) && PieceCol[pos.Pieces[sq-11]] == White {
				AddBlackPawnCapMove(pos, int(sq), int(sq-11), pos.Pieces[sq-11], list)
			}

			if pos.EnPas != NoSquare {
				if sq-9 == pos.EnPas {
					AddEnPassantMove(pos, EncodeMove(int(sq), int(sq-9), Empty, Empty, MFlagEP), list)
				}
				if sq-11 == pos.EnPas {
					AddEnPassantMove(pos, EncodeMove(int(sq), int(sq-11), Empty, Empty, MFlagEP), list)
				}
			}
		}

		// Castle
		if pos.CastlePerm&BKCA != 0 {
			if pos.Pieces[F8] == Empty && pos.Pieces[G8] == Empty {
				if !SqAttacked(E8, White, pos) && !SqAttacked(F8, White, pos) {
					AddQuietMove(pos, EncodeMove(int(E8), int(G8), Empty, Empty, MFlagCA), list)
				}
			}
		}

		if pos.CastlePerm&BQCA != 0 {
			if pos.Pieces[D8] == Empty && pos.Pieces[C8] == Empty && pos.Pieces[B8] == Empty {
				if !SqAttacked(E8, White, pos) && !SqAttacked(D8, White, pos) {
					AddQuietMove(pos, EncodeMove(int(E8), int(C8), Empty, Empty, MFlagCA), list)
				}
			}
		}
	}

	// Generate sliding piece moves (bishops, rooks, queens) by ray-casting along each
	// direction offset until hitting a board edge or another piece. Captures are added
	// for enemy pieces; friendly pieces block further movement on that ray.

	pceIndex = LoopSlideIndex[side]
	piece := LoopSlidePce[pceIndex]
	pceIndex++
	for piece != Empty {
		Assert(PieceValid(piece), "invalid piece")

		for pceNum = 0; pceNum < pos.PceNum[piece]; pceNum++ {
			sq = pos.PList[piece][pceNum]
			Assert(SqOnBoard(sq), "square not on board")

			for index = 0; index < NumDir[piece]; index++ {
				dir = PceDir[piece][index]
				tSq = sq + Square(dir)

				for !SqOffBoard(tSq) {
					if pos.Pieces[tSq] != Empty {
						if PieceCol[pos.Pieces[tSq]] == side^1 {
							AddCaptureMove(pos, EncodeMove(int(sq), int(tSq), pos.Pieces[tSq], Empty, 0), list)
						}
						break
					}
					AddQuietMove(pos, EncodeMove(int(sq), int(tSq), Empty, Empty, 0), list)
					tSq += Square(dir)
				}
			}
		}

		piece = LoopSlidePce[pceIndex]
		pceIndex++
	}

	// Generate non-sliding piece moves (knights, kings) by testing each single-step
	// direction offset. Only one step is taken per direction; captures are added for
	// enemy pieces, blocked squares are skipped.

	pceIndex = LoopNonSlideIndex[side]
	piece = LoopNonSlidePce[pceIndex]
	pceIndex++
	for piece != Empty {
		Assert(PieceValid(piece), "invalid piece")

		for pceNum = 0; pceNum < pos.PceNum[piece]; pceNum++ {
			sq = pos.PList[piece][pceNum]
			Assert(SqOnBoard(sq), "square not on board")

			for index = 0; index < NumDir[piece]; index++ {
				dir = PceDir[piece][index]
				tSq = sq + Square(dir)

				if SqOffBoard(tSq) {
					continue
				}

				// BLACK ^ 1 == WHITE WHITE ^ 1 == BLACK
				if pos.Pieces[tSq] != Empty {
					if PieceCol[pos.Pieces[tSq]] == side^1 {
						AddCaptureMove(pos, EncodeMove(int(sq), int(tSq), pos.Pieces[tSq], Empty, 0), list)
					}
					continue
				}
				AddQuietMove(pos, EncodeMove(int(sq), int(tSq), Empty, Empty, 0), list)
			}
		}

		piece = LoopNonSlidePce[pceIndex]
		pceIndex++
	}
}

// ---------------------------------------------------------------------------
// Pawn move helpers
// ---------------------------------------------------------------------------

// AddWhitePawnCapMove generates all capture moves for a white pawn moving from
// `from` to `to`, capturing the piece `cap`. If the pawn is on the promotion rank
// (Rank7), it generates promotion captures for queen, rook, bishop, and knight.
// Otherwise it adds a single standard pawn capture.
func AddWhitePawnCapMove(pos *Board, from, to int, cap Piece, list *MoveList) {

	Assert(PieceValidEmpty(cap), "capture piece invalid or off range")
	Assert(SqOnBoard(Square(from)), "from square not on board")
	Assert(SqOnBoard(Square(to)), "to square not on board")

	if RanksBrd[from] == Rank7 {
		AddCaptureMove(pos, EncodeMove(from, to, cap, WQ, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, cap, WR, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, cap, WB, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, cap, WN, 0), list)
	} else {
		AddCaptureMove(pos, EncodeMove(from, to, cap, Empty, 0), list)
	}
}

// AddBlackPawnCapMove generates all capture moves for a black pawn moving from
// `from` to `to`, capturing the piece `cap`. If the pawn is on the promotion rank
// (Rank2), it generates promotion captures for queen, rook, bishop, and knight.
// Otherwise it adds a single standard pawn capture.
func AddBlackPawnCapMove(pos *Board, from, to int, cap Piece, list *MoveList) {

	Assert(PieceValidEmpty(cap), "capture piece invalid or off range")
	Assert(SqOnBoard(Square(from)), "from square not on board")
	Assert(SqOnBoard(Square(to)), "to square not on board")

	if RanksBrd[from] == Rank2 {
		AddCaptureMove(pos, EncodeMove(from, to, cap, BQ, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, cap, BR, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, cap, BB, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, cap, BN, 0), list)
	} else {
		AddCaptureMove(pos, EncodeMove(from, to, cap, Empty, 0), list)
	}
}

// AddWhitePawnMove generates all non-capture (quiet) pawn moves for a white pawn
// moving from `from` to `to`. If the pawn is on the promotion rank (Rank7), it
// generates four promotion variants (queen, rook, bishop, knight). Otherwise it
// adds a single forward push.
func AddWhitePawnMove(pos *Board, from, to int, list *MoveList) {

	Assert(SqOnBoard(Square(from)), "from square not on board")
	Assert(SqOnBoard(Square(to)), "to square not on board")

	if RanksBrd[from] == Rank7 {
		AddCaptureMove(pos, EncodeMove(from, to, Empty, WQ, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, Empty, WR, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, Empty, WB, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, Empty, WN, 0), list)
	} else {
		AddCaptureMove(pos, EncodeMove(from, to, Empty, Empty, 0), list)
	}

}

// AddBlackPawnMove generates all non-capture (quiet) pawn moves for a black pawn
// moving from `from` to `to`. If the pawn is on the promotion rank (Rank2), it
// generates four promotion variants (queen, rook, bishop, knight). Otherwise it
// adds a single forward push.
func AddBlackPawnMove(pos *Board, from, to int, list *MoveList) {

	Assert(SqOnBoard(Square(from)), "from square not on board")
	Assert(SqOnBoard(Square(to)), "to square not on board")

	if RanksBrd[from] == Rank2 {
		AddCaptureMove(pos, EncodeMove(from, to, Empty, BQ, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, Empty, BR, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, Empty, BB, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, Empty, BN, 0), list)
	} else {
		AddCaptureMove(pos, EncodeMove(from, to, Empty, Empty, 0), list)
	}

}

// ---------------------------------------------------------------------------
// Encoding / board helpers
// ---------------------------------------------------------------------------

// EncodeMove packs a move's from-square, to-square, captured piece, promoted piece,
// and a flag into a single 28-bit integer using bit shifts:
//
//	from-square in bits 0–6, to-square in bits 7–13, captured in bits 14–17,
//	promoted piece in bits 20–23, and the move flag (e.g., EP, PS, CA) in bit 18+.
func EncodeMove(f, t int, ca, pro Piece, f1 int) int {
	return f | (t << 7) | (int(ca) << 14) | (int(pro) << 20) | f1
}

// SqOffBoard returns true if the given 120-square mailbox index lies on the
// off-board border padding (file padding). This is used during move generation
// to quickly reject moves that would exit the board.
func SqOffBoard(sq Square) bool {
	return FilesBrd[sq] == FileNone
}
