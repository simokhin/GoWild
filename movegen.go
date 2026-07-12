package main

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
	list.Moves[list.Count].Score = 0
	list.Count++
}

// AddEnPassantMove appends an en passant capture move to the move list.
// These are treated as captures but stored with a separate helper for clarity.
func AddEnPassantMove(pos *Board, move int, list *MoveList) {
	list.Moves[list.Count].MoveInt = move
	list.Moves[list.Count].Score = 0
	list.Count++
}

// GenerateAllMoves generates all pseudo-legal moves for the given position and stores
// them in the move list. Currently generates pawn moves (single pushes, double pushes,
// captures, en passant, and promotions) for the side to move. Knight, bishop, rook,
// queen, king, and castling move generation will be added here.
func GenerateAllMoves(pos *Board, list *MoveList) {
	Assert(CheckBoard(pos), "board check failed")

	list.Count = 0

	side := pos.Side
	var sq Square
	var pceNum int

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

			if sq+9 == pos.EnPas {
				AddCaptureMove(pos, EncodeMove(int(sq), int(sq+9), Empty, Empty, MFlagEP), list)
			}
			if sq+11 == pos.EnPas {
				AddCaptureMove(pos, EncodeMove(int(sq), int(sq+11), Empty, Empty, MFlagEP), list)
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

			if sq-9 == pos.EnPas {
				AddCaptureMove(pos, EncodeMove(int(sq), int(sq-9), Empty, Empty, MFlagEP), list)
			}
			if sq-11 == pos.EnPas {
				AddCaptureMove(pos, EncodeMove(int(sq), int(sq-11), Empty, Empty, MFlagEP), list)
			}
		}
	}
}

// AddWhitePawnCapMove generates all capture moves for a white pawn moving from
// `from` to `to`, capturing the piece `cap`. If the pawn is on the promotion rank
// (Rank7), it generates promotion captures for queen, rook, bishop, and knight.
// Otherwise it adds a single standard pawn capture.
func AddWhitePawnCapMove(pos *Board, from, to int, cap Piece, list *MoveList) {
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
	if RanksBrd[from] == Rank2 {
		AddCaptureMove(pos, EncodeMove(from, to, Empty, BQ, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, Empty, BR, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, Empty, BB, 0), list)
		AddCaptureMove(pos, EncodeMove(from, to, Empty, BN, 0), list)
	} else {
		AddCaptureMove(pos, EncodeMove(from, to, Empty, Empty, 0), list)
	}

}

// EncodeMove packs a move's from-square, to-square, captured piece, promoted piece,
// and a flag into a single 28-bit integer using bit shifts:
//   from-square in bits 0–6, to-square in bits 7–13, captured in bits 14–17,
//   promoted piece in bits 20–23, and the move flag (e.g., EP, PS, CA) in bit 18+.
func EncodeMove(f, t int, ca, pro Piece, f1 int) int {
	return f | (t << 7) | (int(ca) << 14) | (int(pro) << 20) | f1
}

// SqOffBoard returns true if the given 120-square mailbox index lies on the
// off-board border padding (file padding). This is used during move generation
// to quickly reject moves that would exit the board.
func SqOffBoard(sq Square) bool {
	return FilesBrd[sq] == FileNone
}
