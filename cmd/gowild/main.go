// Command gowild runs the chess engine's UCI loop.
package main

import "github.com/simokhin/gowild/internal/engine"

func main() {
	engine.AllInit()
	engine.UciLoop()
}
