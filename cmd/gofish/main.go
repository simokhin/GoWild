// Command gofish runs the chess engine's UCI loop.
package main

import "github.com/simokhin/gofish/internal/engine"

func main() {
	engine.AllInit()
	engine.UciLoop()
}
