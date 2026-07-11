package main

func ResetBoard(pos *Board) {
	for index := range 120 {
		pos.Pieces[index] = OffBoard
	}

	for index := range 64 {
		pos.Pieces[SQ120(index)] = Empty
	}

	for index := range 3 {
		pos.BigPce[index] = 0
		pos.MajPce[index] = 0
		pos.MinPce[index] = 0
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
