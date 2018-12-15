package pathfinder

import (
	"sync"
	"time"
)

const (
	minimumAvailablePayloadSize = 5
	defaultPrefix               = "HYX"
	defaultSuffix               = "XYH"
)

var wg sync.WaitGroup

type Path struct {
	Origin          string
	OriginPrefix    string
	OriginSuffix    string
	Time            *time.Time
	FullText        string
	MatchText       string
	DecodedEntities []string
	Metadata        map[string]string
}

func NewPath() *Path {
	var p Path
	p.Metadata = make(map[string]string)

	return &p
}

type Origin interface {
	GetName() string
	FindPath() *Path
	GetMaxPayloadSize() int
	GetPrefix() string
	SetPrefix(s string) error
	GetSuffix() string
	SetSuffix(s string) error
	GetCue() string
	SetCue(s string) error
	Start()
}

var origins []Origin
