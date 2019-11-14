package main

import (
	"fmt"
	"github.com/google/shlex"
	"github.com/kamushadenes/pathfinder"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os/exec"
)

func Handle(p pathfinder.Path) {
	for k := range p.DecodedEntities {
		en := p.DecodedEntities[k]

		switch en.Type {
		case "command":
			pathfinder.Logger.WithFields(pathfinder.GetLogFields(logrus.Fields{
				"time":            p.Time,
				"origin":          p.Origin,
				"fullText":        p.FullText,
				"decodedEntities": p.DecodedEntities,
				"command":         fmt.Sprintf("%s = %s\n", en.Tag, en.Value),
			})).Info("message received")

			c, err := shlex.Split(en.Value)

			if err == nil {
				cmd := exec.Command(c[0], c[1:]...)
				_ = cmd.Run()
			}

			break
		}
	}
}

func main() {
	pathfinder.Handle = Handle

	data, err := ioutil.ReadFile("config.yaml")

	if err != nil {
		log.Fatal(err.Error())
	}

	config, err := pathfinder.ReadConfig(data)

	pathfinder.Run(config)
}
