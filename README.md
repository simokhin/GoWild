# GoWild

GoWild is a UCI-compatible chess engine written in Go.

It started as a Go port of [VICE](https://github.com/bluefeversoft/vice) (Video Instructional Chess Engine) — the C engine built by Richard "Bluefever Software" Allbert in his chess programming video series — and has now reached feature parity with it: a 10x12 board with bitboard piece sets, PV/transposition tables, and alpha-beta search, reimplemented idiomatically in Go. From here, GoWild will continue to evolve independently of VICE.

## Features

- Board representation combining a 120-square array with bitboards per piece/side
- Pseudo-legal move generation for all piece types, including castling and en passant
- Zobrist hashing for position keys
- Alpha-beta search with:
  - Iterative deepening
  - Quiescence search
  - Null-move pruning
  - MVV-LVA move ordering and PV-move ordering
  - Principal variation / transposition hash table
- Evaluation function covering material, piece-square tables, isolated/passed pawns, and open/semi-open files for rooks
- Perft move-generation testing
- UCI protocol support, so it can be used with any UCI-compatible GUI (Arena, CuteChess, Banksia, etc.)

## Status

GoWild is a work in progress. It plays legal chess, speaks UCI, and has the full VICE feature set in place. Strength, evaluation, and search are now being developed further as independent, original work.

## Building

Requires Go 1.26 or later.

```sh
git clone https://github.com/simokhin/gowild.git
cd gowild
go build -o gowild ./cmd/gowild
```

## Usage

GoWild speaks the [UCI protocol](https://www.chessprogramming.org/UCI) over stdin/stdout, so it isn't meant to be played directly from a terminal. Point a UCI-compatible GUI at the built binary, or drive it by hand:

```sh
./gowild
uci
isready
position startpos
go depth 6
```

## Testing

```sh
go test ./...
```

## Acknowledgments

- [VICE](https://github.com/bluefeversoft/vice) by Richard Allbert (Bluefever Software) — the engine GoWild was originally ported from before growing into its own project.
- [Chess Programming Wiki](https://www.chessprogramming.org/) — reference for algorithms and techniques used throughout.

## License

GoWild is licensed under the [MIT License](LICENSE).
