package main

type Piece int8

const (
	Empty Piece = iota
	WP
	WN
	WB
	WR
	WQ
	WK
	BP
	BN
	BB
	BR
	BQ
	BK
)

type File int8

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
	FileNone
)

type Rank int8

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
	RankNone
)

type Color int8

const (
	White Color = iota
	Black
	Both
)

type Square int8

const (
	A1 Square = 21
	B1 Square = 22
	C1 Square = 23
	D1 Square = 24
	E1 Square = 25
	F1 Square = 26
	G1 Square = 27
	H1 Square = 28

	A2 Square = 31
	B2 Square = 32
	C2 Square = 33
	D2 Square = 34
	E2 Square = 35
	F2 Square = 36
	G2 Square = 37
	H2 Square = 38

	A3 Square = 41
	B3 Square = 42
	C3 Square = 43
	D3 Square = 44
	E3 Square = 45
	F3 Square = 46
	G3 Square = 47
	H3 Square = 48

	A4 Square = 51
	B4 Square = 52
	C4 Square = 53
	D4 Square = 54
	E4 Square = 55
	F4 Square = 56
	G4 Square = 57
	H4 Square = 58

	A5 Square = 61
	B5 Square = 62
	C5 Square = 63
	D5 Square = 64
	E5 Square = 65
	F5 Square = 66
	G5 Square = 67
	H5 Square = 68

	A6 Square = 71
	B6 Square = 72
	C6 Square = 73
	D6 Square = 74
	E6 Square = 75
	F6 Square = 76
	G6 Square = 77
	H6 Square = 78

	A7 Square = 81
	B7 Square = 82
	C7 Square = 83
	D7 Square = 84
	E7 Square = 85
	F7 Square = 86
	G7 Square = 87
	H7 Square = 88

	A8 Square = 91
	B8 Square = 92
	C8 Square = 93
	D8 Square = 94
	E8 Square = 95
	F8 Square = 96
	G8 Square = 97
	H8 Square = 98

	NoSquare Square = 99
)

var board [120]Piece
