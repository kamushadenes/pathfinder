package main

import (
	"fmt"
	"github.com/kamushadenes/pathfinder"
	"os"
)

func Handle(p pathfinder.Path) {
	fmt.Printf("command received: %s\n", p.DecodedEntities)
}

func main() {
	pathfinder.Handle = Handle

	pathfinder.RegPrefix = "HYX"
	pathfinder.RegSuffix = "XYH"

	pathfinder.TagMap["potato"] = "shutdown -h now"

	if len(os.Args) >= 2 {
		pathfinder.Run(os.Args[1])
	}
}
