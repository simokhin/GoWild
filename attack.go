package main

// KnDir lists the 8 knight move offsets in the 120-square mailbox board.
// For a knight at square S, the reachable squares are S+offsets[i].
// Offsets that would wrap around the board edges are caught by the off-board border.
var KnDir = [8]int{-8, -19, -21, -12, 8, 19, 21, 12}

// RkDir lists the 4 rook move direction offsets (left, down, right, up)
// in the 120-square mailbox board. Used for sliding-piece attack/move generation.
var RkDir = [4]int{-1, -10, 1, 10}

// BiDir lists the 4 bishop move direction offsets (down-left, down-right, up-left, up-right)
// in the 120-square mailbox board. Used for sliding-piece attack/move generation.
var BiDir = [4]int{-9, -11, 11, 9}

// KiDir lists the 8 king move direction offsets (non-sliding, one step in each direction).
// Combines rook directions and bishop directions for the king's full movement set.
var KiDir = [8]int{-1, -10, 1, 10, -9, -11, 11, 9}

// SqAttacked checks whether the given square is attacked by any piece of the given side.
// It tests pawn attacks first (by file offsets), then knights (8 directions), then
// sliding pieces (rooks/queens along ranks/files, bishops/queens along diagonals),
// and finally king attacks (8 directions).
func SqAttacked(sq Square, side Color, pos *Board) bool {
	Assert(SqOnBoard(sq), "square not on board")
	Assert(SideValid(side), "invalid side")
	Assert(CheckBoard(pos), "board check failed")

	var piece Piece
	var index int
	var tSq Square
	var dir int

	// Pawn attacks: white pawns attack diagonally up-left (-11) and up-right (-9);
	// black pawns attack diagonally down-left (+9) and down-right (+11).
	if side == White {
		if pos.Pieces[sq-11] == WP || pos.Pieces[sq-9] == WP {
			return true
		}
	} else {
		if pos.Pieces[sq+11] == BP || pos.Pieces[sq+9] == BP {
			return true
		}
	}

	// Knight attacks: check all 8 L-shaped offsets for a knight of the given side.
	for index = 0; index < 8; index++ {
		piece = pos.Pieces[sq+Square(KnDir[index])]
		if piece != OffBoard && piece != Empty && IsKn(piece) && PieceCol[piece] == side {
			return true
		}
	}

	// Rook/queen attacks: slide along all 4 orthogonal directions until hitting
	// a piece or the board edge. Return true if a rook or queen of the given side is found.
	for index = 0; index < 4; index++ {
		dir = RkDir[index]
		tSq = sq + Square(dir)
		piece = pos.Pieces[tSq]
		for piece != OffBoard {
			if piece != Empty {
				if IsRQ(piece) && PieceCol[piece] == side {
					return true
				}
				break
			}
			tSq += Square(dir)
			piece = pos.Pieces[tSq]
		}
	}

	// Bishop/queen attacks: slide along all 4 diagonal directions until hitting
	// a piece or the board edge. Return true if a bishop or queen of the given side is found.
	for index = 0; index < 4; index++ {
		dir = BiDir[index]
		tSq = sq + Square(dir)
		piece = pos.Pieces[tSq]
		for piece != OffBoard {
			if piece != Empty {
				if IsBQ(piece) && PieceCol[piece] == side {
					return true
				}
				break
			}
			tSq += Square(dir)
			piece = pos.Pieces[tSq]
		}
	}

	// King attacks: check all 8 king-step offsets for the enemy king.
	for index = range 8 {
		piece = pos.Pieces[sq+Square(KiDir[index])]
		if piece != OffBoard && piece != Empty && IsKi(piece) && PieceCol[piece] == side {
			return true
		}
	}

	return false

}
