package engine

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

// Win at Chess (WAC) is a classic 300-position tactical test suite
// (https://www.chessprogramming.org/Win_at_Chess). Each position has one or
// more known correct "best moves" in SAN. Running the engine's search
// against it is a standard way to regression-test AlphaBeta's tactical
// strength: if a future change to move ordering, pruning, or evaluation
// quietly weakens the search, the solve rate drops.
const (
	wacTimePerPositionMs = 1000
	wacMaxSearchDepth    = 30
	// Search is time-budgeted, not depth-budgeted, so the solve rate varies
	// a few points run to run with machine load (observed 78-82% across
	// repeated runs on the same machine). Set well below that to avoid
	// flaking on noise while still catching a genuine search regression.
	wacMinSolveRate = 0.65
)

var (
	wacBestMoveRe = regexp.MustCompile(`bm\s+([^;]+);`)
	wacIDRe       = regexp.MustCompile(`id\s+"([^"]+)"`)
)

func TestWinAtChessAlphaBeta(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Win at Chess suite in -short mode")
	}

	data, err := os.ReadFile("testdata/win_at_chess_positions.epd")
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	solved := 0
	total := 0

	// Shared across positions and cleared between them: allocating a fresh
	// PvSize-d table (128MB) per position, 300 times over, makes the suite
	// far slower than the search itself warrants.
	hashTable := &HashTable{}
	InitHashTable(hashTable)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		bmMatch := wacBestMoveRe.FindStringSubmatch(line)
		if bmMatch == nil {
			t.Fatalf("line missing bm field: %s", line)
		}
		wantSAN := strings.Fields(bmMatch[1])

		id := line
		if idMatch := wacIDRe.FindStringSubmatch(line); idMatch != nil {
			id = idMatch[1]
		}

		ClearHashTable(hashTable)
		board := &Board{HashTable: hashTable}

		if ParseFEN(line, board) != 0 {
			t.Fatalf("%s: failed to parse FEN", id)
		}

		wantMoves := make(map[int]bool, len(wantSAN))
		for _, san := range wantSAN {
			move, err := parseSAN(board, san)
			if err != nil {
				t.Fatalf("%s: %v", id, err)
			}
			wantMoves[move] = true
		}

		info := &SearchInfo{}
		bestMove := searchBestMove(board, info, wacMaxSearchDepth, wacTimePerPositionMs)

		total++
		if wantMoves[bestMove] {
			solved++
		} else {
			t.Logf("%s: expected one of %v, got %s", id, wantSAN, PrMove(bestMove))
		}
	}

	rate := float64(solved) / float64(total)
	t.Logf("Win at Chess: solved %d/%d (%.1f%%)", solved, total, rate*100)

	if rate < wacMinSolveRate {
		t.Errorf("solve rate %.1f%% below minimum %.1f%%", rate*100, wacMinSolveRate*100)
	}
}

// searchBestMove runs iterative deepening over AlphaBeta, stopping at
// maxDepth or once timeMs milliseconds have elapsed, and returns the best
// move found for the position currently on pos. It mirrors SearchPosition's
// loop but skips the UCI-style progress output.
func searchBestMove(pos *Board, info *SearchInfo, maxDepth int, timeMs int64) int {
	ClearForSearch(pos, info)
	info.TimeSet = true
	info.StartTime = GetTimeMs()
	info.StopTime = info.StartTime + timeMs
	info.Depth = maxDepth

	bestMove := NoMove
	for currentDepth := 1; currentDepth <= info.Depth; currentDepth++ {
		AlphaBeta(-Infinite, Infinite, currentDepth, pos, info, true)

		if info.Stopped.Load() {
			break
		}

		GetPvLine(currentDepth, pos)
		bestMove = pos.PvArray[0]
	}

	return bestMove
}

// sanMovePattern matches a (non-castling) SAN move, e.g. "Qxh7+", "Nc3",
// "dxe6", "e8=Q". Capture groups: piece letter, disambiguating file,
// disambiguating rank, destination square, promotion piece.
var sanMovePattern = regexp.MustCompile(`^([NBRQK]?)([a-h]?)([1-8]?)x?([a-h][1-8])(?:=?([NBRQ]))?[+#]?$`)

// parseSAN resolves a SAN move string (as used in EPD "bm"/"am" fields) to
// the engine's encoded move integer for the given position, by matching it
// against the position's legal moves. Returns an error if the move can't be
// parsed or doesn't identify exactly one legal move.
func parseSAN(pos *Board, san string) (int, error) {
	san = strings.TrimSpace(san)

	if san == "O-O" || san == "O-O-O" {
		return findCastleMove(pos, san == "O-O")
	}

	m := sanMovePattern.FindStringSubmatch(san)
	if m == nil {
		return NoMove, fmt.Errorf("cannot parse SAN move %q", san)
	}

	pieceLetter, fromFile, fromRank, dest, promoLetter := m[1], m[2], m[3], m[4], m[5]

	destSq := FR2SQ(File(dest[0]-'a'), Rank(dest[1]-'1'))
	wantPiece := sanPieceForSide(pieceLetter, pos.Side)

	var wantPromoted Piece
	if promoLetter != "" {
		wantPromoted = sanPieceForSide(promoLetter, pos.Side)
	}

	list := &MoveList{}
	GenerateAllMoves(pos, list)

	var matches []int

	for i := 0; i < list.Count; i++ {
		move := list.Moves[i].MoveInt
		from := Square(FromSq(move))
		to := Square(ToSq(move))

		if to != destSq || pos.Pieces[from] != wantPiece {
			continue
		}
		if promoLetter != "" && Promoted(move) != wantPromoted {
			continue
		}
		if promoLetter == "" && Promoted(move) != Empty {
			continue
		}
		if fromFile != "" && FilesBrd[from] != File(fromFile[0]-'a') {
			continue
		}
		if fromRank != "" && RanksBrd[from] != Rank(fromRank[0]-'1') {
			continue
		}

		if !MakeMove(pos, move) {
			continue
		}
		TakeMove(pos)

		matches = append(matches, move)
	}

	if len(matches) != 1 {
		return NoMove, fmt.Errorf("SAN move %q matched %d legal moves, want exactly 1", san, len(matches))
	}

	return matches[0], nil
}

// findCastleMove returns the legal castling move for pos matching the
// requested side (kingside or queenside), if any.
func findCastleMove(pos *Board, kingside bool) (int, error) {
	list := &MoveList{}
	GenerateAllMoves(pos, list)

	for i := 0; i < list.Count; i++ {
		move := list.Moves[i].MoveInt
		if move&MFlagCA == 0 {
			continue
		}

		to := Square(ToSq(move))
		if (kingside && FilesBrd[to] == FileG) || (!kingside && FilesBrd[to] == FileC) {
			if !MakeMove(pos, move) {
				continue
			}
			TakeMove(pos)
			return move, nil
		}
	}

	side := "queenside"
	if kingside {
		side = "kingside"
	}
	return NoMove, fmt.Errorf("no legal %s castle", side)
}

// sanPieceForSide maps a SAN piece letter ("" for pawn, "N", "B", "R", "Q",
// "K") to the concrete Piece constant for the given side to move.
func sanPieceForSide(letter string, side Color) Piece {
	white := map[string]Piece{"": WP, "N": WN, "B": WB, "R": WR, "Q": WQ, "K": WK}
	black := map[string]Piece{"": BP, "N": BN, "B": BB, "R": BR, "Q": BQ, "K": BK}

	if side == White {
		return white[letter]
	}
	return black[letter]
}
