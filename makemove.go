package main

// ---------------------------------------------------------------------------
// Zobrist hashing helpers
// ---------------------------------------------------------------------------

// HashPce toggles the Zobrist hash key for the given piece on the given square.
// Used during make/unmake to incrementally update PosKey.
func HashPce(pos *Board, pce Piece, sq Square) {
	pos.PosKey ^= uint64(PieceKeys[pce][sq])
}

// HashCA toggles the Zobrist hash key for the current castling permissions.
// Must be called both before and after changing CastlePerm to keep PosKey in sync.
func HashCA(pos *Board) {
	pos.PosKey ^= uint64(CastleKeys[pos.CastlePerm])
}

// HashSide toggles the Zobrist hash key for the side to move.
// Called when flipping pos.Side during make/unmake.
func HashSide(pos *Board) {
	pos.PosKey ^= uint64(SideKey)
}

// HashEP toggles the Zobrist hash key for the en passant square.
// Must be called both before and after changing pos.EnPas.
func HashEP(pos *Board) {
	pos.PosKey ^= uint64(PieceKeys[Empty][pos.EnPas])
}

// ---------------------------------------------------------------------------
// Castle permission mask
// ---------------------------------------------------------------------------

// CastlePermMask maps each 120-square mailbox index to a bitmask that preserves
// only the castling rights which the given square does NOT invalidate.
// Applied via AND when a piece moves from or to that square.
var CastlePermMask = [120]int{
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 13, 15, 15, 15, 12, 15, 15, 14, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 7, 15, 15, 15, 3, 15, 15, 11, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 15,
}

// ---------------------------------------------------------------------------
// Piece management helpers
// ---------------------------------------------------------------------------

// ClearPiece removes the piece on square sq from the board and updates all
// associated data structures: mailbox (Pieces), piece list (PList), piece counts
// (PceNum, BigPce, MajPce, MinPce), material totals, pawn bitboards, and PosKey.
func ClearPiece(sq Square, pos *Board) {
	Assert(SqOnBoard(sq), "square not on board")
	Assert(CheckBoard(pos), "board check failed")

	pce := pos.Pieces[sq]
	Assert(PieceValid(pce), "invalid piece")

	col := PieceCol[pce]
	var index int
	tPceNum := -1

	Assert(SideValid(col), "invalid side")

	HashPce(pos, pce, sq)

	pos.Pieces[sq] = Empty
	pos.Material[col] -= PieceVal[pce]

	if PieceBig[pce] {
		pos.BigPce[col]--
		if PieceMaj[pce] {
			pos.MajPce[col]--
		} else {
			pos.MinPce[col]--
		}
	} else {
		CLRBIT(&pos.Pawns[col], SQ64(sq))
		CLRBIT(&pos.Pawns[Both], SQ64(sq))
	}

	for index = 0; index < pos.PceNum[pce]; index++ {
		if pos.PList[pce][index] == sq {
			tPceNum = index
			break
		}
	}

	Assert(tPceNum != -1, "piece not found in piece list")
	Assert(tPceNum >= 0 && tPceNum < 10, "piece index out of range")

	pos.PceNum[pce]--
	pos.PList[pce][tPceNum] = pos.PList[pce][pos.PceNum[pce]]
}

// AddPiece places a piece on square sq and updates all associated data structures:
// mailbox (Pieces), piece list (PList), piece counts, material totals, pawn
// bitboards, and PosKey. The inverse of ClearPiece.
func AddPiece(sq Square, pos *Board, pce Piece) {
	Assert(PieceValid(pce), "invalid piece")
	Assert(SqOnBoard(sq), "square not on board")

	col := PieceCol[pce]
	Assert(SideValid(col), "invalid side")

	HashPce(pos, pce, sq)

	pos.Pieces[sq] = pce
	if PieceBig[pce] {
		pos.BigPce[col]++
		if PieceMaj[pce] {
			pos.MajPce[col]++
		} else {
			pos.MinPce[col]++
		}
	} else {
		SETBIT(&pos.Pawns[col], SQ64(sq))
		SETBIT(&pos.Pawns[Both], SQ64(sq))
	}

	pos.Material[col] += PieceVal[pce]
	pos.PList[pce][pos.PceNum[pce]] = sq
	pos.PceNum[pce]++
}

// MovePiece relocates a piece from square from to square to, updating the mailbox,
// piece list, pawn bitboards (if the piece is a pawn), and PosKey. Does NOT handle
// captures — the piece at the destination is simply overwritten.
func MovePiece(from, to Square, pos *Board) {
	Assert(SqOnBoard(from), "from square not on board")
	Assert(SqOnBoard(to), "to square not on board")

	var index int
	pce := pos.Pieces[from]
	col := PieceCol[pce]
	Assert(SideValid(col), "invalid side")
	Assert(PieceValid(pce), "invalid piece")

	HashPce(pos, pce, from)
	pos.Pieces[from] = Empty

	HashPce(pos, pce, to)
	pos.Pieces[to] = pce

	if !PieceBig[pce] {
		CLRBIT(&pos.Pawns[col], SQ64(from))
		CLRBIT(&pos.Pawns[Both], SQ64(from))
		SETBIT(&pos.Pawns[col], SQ64(to))
		SETBIT(&pos.Pawns[Both], SQ64(to))
	}
	for index = 0; index < pos.PceNum[pce]; index++ {
		if pos.PList[pce][index] == from {
			pos.PList[pce][index] = to
			break
		}
	}
}

// ---------------------------------------------------------------------------
// Make / Unmake move
// ---------------------------------------------------------------------------

// MakeMove applies a move to the board and returns true if the resulting position
// is legal (the moving side's king is not left in check). If illegal, the move is
// taken back via TakeMove and false is returned.
//
// Handles: normal moves, captures, en passant, castling, pawn double-push (setting
// the en passant square), and pawn promotion. Incrementally updates the Zobrist
// hash, castling permissions, fifty-move counter, and search ply.
func MakeMove(pos *Board, move int) bool {
	Assert(CheckBoard(pos), "board check failed")

	from := FromSq(move)
	to := ToSq(move)
	side := pos.Side

	Assert(SqOnBoard(Square(from)), "from square not on board")
	Assert(SqOnBoard(Square(to)), "to square not on board")
	Assert(SideValid(side), "invalid side")
	Assert(PieceValid(pos.Pieces[from]), "invalid piece at from")
	Assert(pos.HisPly() >= 0 && pos.HisPly() < MaxGameMoves, "hisPly out of range")
	Assert(pos.Ply >= 0 && pos.Ply < MaxDepth, "ply out of range")

	pos.History = append(pos.History, Undo{PosKey: pos.PosKey})
	hist := len(pos.History) - 1

	if move&MFlagEP != 0 {
		if side == White {
			ClearPiece(Square(to-10), pos)
		} else {
			ClearPiece(Square(to+10), pos)
		}
	} else if move&MFlagCA != 0 {
		switch Square(to) {
		case C1:
			MovePiece(A1, D1, pos)
		case C8:
			MovePiece(A8, D8, pos)
		case G1:
			MovePiece(H1, F1, pos)
		case G8:
			MovePiece(H8, F8, pos)
		default:
			Assert(false, "invalid castle target square")
		}
	}

	if pos.EnPas != NoSquare {
		HashEP(pos)
	}
	HashCA(pos)

	pos.History[hist].Move = move
	pos.History[hist].FiftyMove = pos.FiftyMove
	pos.History[hist].EnPas = pos.EnPas
	pos.History[hist].CastlePerm = pos.CastlePerm

	pos.CastlePerm &= CastlePerm(CastlePermMask[from])
	pos.CastlePerm &= CastlePerm(CastlePermMask[to])
	pos.EnPas = NoSquare

	HashCA(pos)

	captured := Captured(move)
	pos.FiftyMove++

	if captured != Empty {
		Assert(PieceValid(captured), "invalid captured piece")
		ClearPiece(Square(to), pos)
		pos.FiftyMove = 0
	}

	pos.Ply++

	Assert(pos.HisPly() >= 0 && pos.HisPly() < MaxGameMoves, "hisPly out of range")
	Assert(pos.Ply >= 0 && pos.Ply < MaxDepth, "ply out of range")

	if PiecePawn[pos.Pieces[from]] {
		pos.FiftyMove = 0
		if move&MFlagPS != 0 {
			if side == White {
				pos.EnPas = Square(from) + 10
				Assert(RanksBrd[pos.EnPas] == Rank3, "en passant square not on rank 3")
			} else {
				pos.EnPas = Square(from) - 10
				Assert(RanksBrd[pos.EnPas] == Rank6, "en passant square not on rank 6")
			}
			HashEP(pos)
		}
	}

	MovePiece(Square(from), Square(to), pos)

	prPce := Promoted(move)
	if prPce != Empty {
		Assert(PieceValid(prPce) && !PiecePawn[prPce], "invalid promotion piece")
		ClearPiece(Square(to), pos)
		AddPiece(Square(to), pos, prPce)
	}

	if PieceKing[pos.Pieces[to]] {
		pos.KingSq[pos.Side] = Square(to)
	}

	pos.Side ^= 1
	HashSide(pos)
	Assert(CheckBoard(pos), "board check failed")

	if SqAttacked(pos.KingSq[side], pos.Side, pos) {
		TakeMove(pos)
		return false
	}

	return true
}

// TakeMove reverses the last move made on the board, restoring all state from the
// Undo snapshot stored in pos.History. Handles unmaking of en passant captures,
// castling rook moves, regular captures, and promotions.
func TakeMove(pos *Board) {
	Assert(CheckBoard(pos), "board check failed")

	pos.Ply--

	Assert(len(pos.History) > 0, "history empty")
	hist := len(pos.History) - 1
	undo := pos.History[hist]
	pos.History = pos.History[:hist]

	Assert(pos.HisPly() >= 0 && pos.HisPly() < MaxGameMoves, "hisPly out of range")
	Assert(pos.Ply >= 0 && pos.Ply < MaxDepth, "ply out of range")

	move := undo.Move
	from := FromSq(move)
	to := ToSq(move)

	Assert(SqOnBoard(Square(from)), "from square not on board")
	Assert(SqOnBoard(Square(to)), "to square not on board")

	if pos.EnPas != NoSquare {
		HashEP(pos)
	}
	HashCA(pos)

	pos.CastlePerm = undo.CastlePerm
	pos.FiftyMove = undo.FiftyMove
	pos.EnPas = undo.EnPas

	if pos.EnPas != NoSquare {
		HashEP(pos)
	}
	HashCA(pos)

	pos.Side ^= 1
	HashSide(pos)

	if move&MFlagEP != 0 {
		if pos.Side == White {
			AddPiece(Square(to)-10, pos, BP)
		} else {
			AddPiece(Square(to)+10, pos, WP)
		}
	} else if move&MFlagCA != 0 {
		switch Square(to) {
		case C1:
			MovePiece(D1, A1, pos)
		case C8:
			MovePiece(D8, A8, pos)
		case G1:
			MovePiece(F1, H1, pos)
		case G8:
			MovePiece(F8, H8, pos)
		default:
			Assert(false, "invalid castle target square")
		}
	}

	MovePiece(Square(to), Square(from), pos)

	if PieceKing[pos.Pieces[from]] {
		pos.KingSq[pos.Side] = Square(from)
	}

	captured := Captured(move)
	if captured != Empty {
		Assert(PieceValid(captured), "invalid captured piece")
		AddPiece(Square(to), pos, captured)
	}

	if Promoted(move) != Empty {
		promoted := Promoted(move)
		Assert(PieceValid(promoted) && !PiecePawn[promoted], "invalid promoted piece")
		ClearPiece(Square(from), pos)
		if PieceCol[promoted] == White {
			AddPiece(Square(from), pos, WP)
		} else {
			AddPiece(Square(from), pos, BP)
		}
	}

	Assert(CheckBoard(pos), "board check failed")
}
