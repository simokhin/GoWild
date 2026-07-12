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

// GenerateAllMoves generates all legal moves for the given position and stores
// them in the move list. Currently a stub that resets the list; move generation
// logic for each piece type will be added here.
func GenerateAllMoves(pos *Board, list *MoveList) {
	list.Count = 0
}
