package main

func GeneratePosKey(pos *Board) uint64 {
	var finalKey Bitboard = 0
	var piece Piece = Empty

	// pieces
	for sq := range 120 {
		piece = pos.Pieces[sq]
		if piece != Empty {
			Assert(piece >= WP && piece <= BK, "piece out of valid range")
			finalKey ^= PieceKeys[piece][sq]
		}
	}

	if pos.Side == White {
		finalKey ^= SideKey
	}

	if pos.EnPas != NoSquare {
		Assert(pos.EnPas >= 0 && int(pos.EnPas) < 120, "en passant square out of range")
		finalKey ^= PieceKeys[Empty][pos.EnPas]
	}

	Assert(pos.CastlePerm >= 0 && pos.CastlePerm <= 15, "castle perm out of range")

	finalKey ^= CastleKeys[pos.CastlePerm]

	return uint64(finalKey)
}
