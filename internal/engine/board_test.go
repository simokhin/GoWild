package engine

import (
	"os"
	"strings"
	"testing"
)

func TestMirrorBoard(t *testing.T) {
	data, err := os.ReadFile("testdata/mirror_positions.epd")
	if err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		fen := strings.TrimSpace(line)
		if fen == "" {
			continue
		}

		board := &Board{}
		board.PvTable = &PVTable{}
		InitPvTable(board.PvTable)

		if ParseFEN(fen, board) != 0 {
			t.Errorf("line %d: failed to parse FEN: %s", i, fen)
			continue
		}

		eval1 := EvalPosition(board)
		MirrorBoard(board)
		eval2 := EvalPosition(board)

		if eval1 != eval2 {
			t.Errorf("line %d: mirror eval mismatch for %s: got %d, mirrored %d", i, fen, eval1, eval2)
		}
	}
}
