package engine

var PawnTable = [64]int{
	0, 0, 0, 0, 0, 0, 0, 0,
	10, 10, 0, -10, -10, 0, 10, 10,
	5, 0, 0, 5, 5, 0, 0, 5,
	0, 0, 10, 20, 20, 10, 0, 0,
	5, 5, 5, 10, 10, 5, 5, 5,
	10, 10, 10, 20, 20, 10, 10, 10,
	20, 20, 20, 30, 30, 20, 20, 20,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var KnightTable = [64]int{
	0, -10, 0, 0, 0, 0, -10, 0,
	0, 0, 0, 5, 5, 0, 0, 0,
	0, 0, 10, 10, 10, 10, 0, 0,
	0, 0, 10, 20, 20, 10, 0, 0,
	5, 10, 15, 20, 20, 15, 10, 5,
	5, 10, 10, 20, 20, 10, 10, 5,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var BishopTable = [64]int{
	0, 0, -10, 0, 0, -10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 10, 15, 20, 20, 15, 10, 0,
	0, 0, 10, 15, 15, 10, 0, 0,
	0, 0, 0, 10, 10, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var RookTable = [64]int{
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	0, 0, 5, 10, 10, 5, 0, 0,
	25, 25, 25, 25, 25, 25, 25, 25,
	0, 0, 5, 10, 10, 5, 0, 0,
}

var KingE = [64]int{
	-50, -10, 0, 0, 0, 0, -10, -50,
	-10, 0, 10, 10, 10, 10, 0, -10,
	0, 10, 20, 20, 20, 20, 10, 0,
	0, 10, 20, 40, 40, 20, 10, 0,
	0, 10, 20, 40, 40, 20, 10, 0,
	0, 10, 20, 20, 20, 20, 10, 0,
	-10, 0, 10, 10, 10, 10, 0, -10,
	-50, -10, 0, 0, 0, 0, -10, -50,
}

var KingO = [64]int{
	0, 5, 5, -10, -10, 0, 10, 5,
	-30, -30, -30, -30, -30, -30, -30, -30,
	-50, -50, -50, -50, -50, -50, -50, -50,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
	-70, -70, -70, -70, -70, -70, -70, -70,
}

var Mirror64 = [64]int{
	56, 57, 58, 59, 60, 61, 62, 63,
	48, 49, 50, 51, 52, 53, 54, 55,
	40, 41, 42, 43, 44, 45, 46, 47,
	32, 33, 34, 35, 36, 37, 38, 39,
	24, 25, 26, 27, 28, 29, 30, 31,
	16, 17, 18, 19, 20, 21, 22, 23,
	8, 9, 10, 11, 12, 13, 14, 15,
	0, 1, 2, 3, 4, 5, 6, 7,
}

const PawnIsolated = -10

var PawnPassed = [8]int{0, 5, 10, 20, 35, 60, 100, 200}

const RookOpenFile = 10

func EvalPosition(pos *Board) int {
	score := pos.Material[White] - pos.Material[Black]

	pce := WP
	for pceNum := 0; pceNum < pos.PceNum[pce]; pceNum++ {
		sq := pos.PList[pce][pceNum]
		Assert(SqOnBoard(sq), "square not on board")
		score += PawnTable[SQ64(sq)]

		if IsolatedMask[SQ64(sq)]&pos.Pawns[White] == 0 {
			score += PawnIsolated
		}

		if WhitePassedMask[SQ64(sq)]&pos.Pawns[Black] == 0 {
			score += PawnPassed[RanksBrd[sq]]
		}
	}

	pce = BP
	for pceNum := 0; pceNum < pos.PceNum[pce]; pceNum++ {
		sq := pos.PList[pce][pceNum]
		Assert(SqOnBoard(sq), "square not on board")
		score -= PawnTable[Mirror64[SQ64(sq)]]

		if IsolatedMask[SQ64(sq)]&pos.Pawns[Black] == 0 {
			score -= PawnIsolated
		}

		if BlackPassedMask[SQ64(sq)]&pos.Pawns[White] == 0 {
			score -= PawnPassed[7-RanksBrd[sq]]
		}
	}

	pce = WN
	for pceNum := 0; pceNum < pos.PceNum[pce]; pceNum++ {
		sq := pos.PList[pce][pceNum]
		Assert(SqOnBoard(sq), "square not on board")
		score += KnightTable[SQ64(sq)]
	}

	pce = BN
	for pceNum := 0; pceNum < pos.PceNum[pce]; pceNum++ {
		sq := pos.PList[pce][pceNum]
		Assert(SqOnBoard(sq), "square not on board")
		score -= KnightTable[Mirror64[SQ64(sq)]]
	}

	pce = WB
	for pceNum := 0; pceNum < pos.PceNum[pce]; pceNum++ {
		sq := pos.PList[pce][pceNum]
		Assert(SqOnBoard(sq), "square not on board")
		score += BishopTable[SQ64(sq)]
	}

	pce = BB
	for pceNum := 0; pceNum < pos.PceNum[pce]; pceNum++ {
		sq := pos.PList[pce][pceNum]
		Assert(SqOnBoard(sq), "square not on board")
		score -= BishopTable[Mirror64[SQ64(sq)]]
	}

	pce = WR
	for pceNum := 0; pceNum < pos.PceNum[pce]; pceNum++ {
		sq := pos.PList[pce][pceNum]
		Assert(SqOnBoard(sq), "square not on board")
		score += RookTable[SQ64(sq)]
	}

	pce = BR
	for pceNum := 0; pceNum < pos.PceNum[pce]; pceNum++ {
		sq := pos.PList[pce][pceNum]
		Assert(SqOnBoard(sq), "square not on board")
		score -= RookTable[Mirror64[SQ64(sq)]]
	}

	if pos.Side == White {
		return score
	}

	return -score
}
