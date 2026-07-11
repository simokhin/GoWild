package main

// PceChar maps each internal Piece constant to its printable character.
// Index by Piece value: Empty('.') -> King('k').
// The order is: ., P, N, B, R, Q, K, p, n, b, r, q, k
var PceChar = ".PNBRQKpnbrqk"

// SideChar maps each Color constant to its single-character representation.
// Index by Color: White -> 'w', Black -> 'b', Both -> '-'.
var SideChar = "wb-"

// RankChar maps each Rank constant (0-7) to its chess notation digit (1-8).
var RankChar = "12345678"

// FileChar maps each File constant (0-7) to its chess notation letter (a-h).
var FileChar = "abcdefgh"
