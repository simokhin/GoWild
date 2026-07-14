package engine

// SqOnBoard returns true if the given 120-square mailbox index corresponds to
// a real square on the board (i.e., its file is not FileNone, meaning it is not
// part of the off-board border padding).
func SqOnBoard(sq Square) bool {
	return FilesBrd[sq] != FileNone
}

// SideValid returns true if the given side is either White or Black.
// Used to validate the side-to-move or piece colour before operations.
func SideValid(side Color) bool {
	return side == White || side == Black
}

// FileRankValid returns true if the given integer is a valid file or rank index
// (0 through 7, corresponding to A–H or 1–8).
func FileRankValid(fr int) bool {
	return fr >= 0 && fr <= 7
}

// PieceValidEmpty returns true if the given piece value is within the valid
// range including Empty (Empty through BK). Use this when a square may be empty.
func PieceValidEmpty(pce Piece) bool {
	return pce >= Empty && pce <= BK
}

// PieceValid returns true if the given piece value is a real piece (WP through BK).
// Unlike PieceValidEmpty, this excludes Empty so it only matches actual piece types.
func PieceValid(pce Piece) bool {
	return pce >= WP && pce <= BK
}
