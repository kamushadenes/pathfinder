package main

import (
	"fmt"
	"github.com/kamushadenes/pathfinder"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) >= 3 {
		switch os.Args[1] {
		case "encode":
			s, err := pathfinder.EncodeString(strings.Join(os.Args[2:], " "))
			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Println(s)

			os.Exit(0)
		case "decode":
			s, err := pathfinder.DecodeString(os.Args[2])
			if err != nil {
				log.Fatal(err.Error())
			}

			fmt.Println(s)

			os.Exit(0)
		}
	}
}
