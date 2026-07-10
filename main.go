package main

const offBoard = 65

func InitSq120ToSq64() {
	for index := 0; index < 120; index++ {
		Sq120ToSq64[index] = offBoard
	}
	for index := 0; index < 64; index++ {
		Sq64ToSq120[index] = 120
	}

	sq64 := 0
	for rank := Rank1; rank <= Rank8; rank++ {
		for file := FileA; file <= FileH; file++ {
			sq := FR2SQ(file, rank)
			Sq64ToSq120[sq64] = sq
			Sq120ToSq64[sq] = sq64
			sq64++
		}
	}
}

func main() {
	InitSq120ToSq64()
}
