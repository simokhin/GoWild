package main

const offBoard = 65 // Sentinel value used to mark off-board squares in the 120-to-64 mapping

// InitSq120ToSq64 builds the lookup tables that translate between the 120-square
// mailbox board indices and the compact 64-square array indices.
func InitSq120ToSq64() {
	for index := range 120 {
		Sq120ToSq64[index] = offBoard
	}
	for index := range 64 {
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
