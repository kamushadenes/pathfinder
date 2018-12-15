package pathfinder

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var Handle func(p Path)

var RegPrefix = defaultPrefix
var RegSuffix = defaultSuffix

var TagMap = make(map[string]string)
var TagSep = " "
var TagTrim = "#"

func HandleArbitrary(job Path) (Path, error) {
	re := regexp.MustCompile(fmt.Sprintf("%s(.*?)%s", job.OriginPrefix, job.OriginSuffix))

	match := re.FindStringSubmatch(job.FullText)

	if len(match) != 2 {
		Logger.WithFields(GetLogFields(logrus.Fields{
			"time":     job.Time,
			"origin":   job.Origin,
			"fullText": job.FullText,
		})).Debug("regex match failed")
		return job, NewRegexMatchFailed()
	}

	job.MatchText = match[1]

	text, err := DecodeString(job.MatchText)
	if err != nil {
		Logger.WithFields(GetLogFields(logrus.Fields{
			"time":      job.Time,
			"origin":    job.Origin,
			"fullText":  job.FullText,
			"matchText": job.MatchText,
		})).Error("failed to decode text")
		return job, NewDecodeFailed()
	} else {
		job.DecodedEntities = append(job.DecodedEntities, text)
	}

	return job, nil
}

func HandleTagged(job Path) (Path, error) {
	tags := strings.Split(job.FullText, TagSep)

	for _, tag := range tags {
		tag = strings.Trim(tag, TagTrim)
		if _, ok := TagMap[tag]; ok {
			job.DecodedEntities = append(job.DecodedEntities, TagMap[tag])
		}
	}
	return job, nil
}

func worker(jobChan <-chan Path) {
	for job := range jobChan {
		var err error

		job, err = HandleArbitrary(job)
		if err != nil {
			job, err = HandleTagged(job)
			if err != nil {
				Logger.WithFields(GetLogFields(logrus.Fields{
					"time":     job.Time,
					"origin":   job.Origin,
					"fullText": job.FullText,
				})).Error("failed to decode text")
				continue
			}
		}

		Logger.WithFields(GetLogFields(logrus.Fields{
			"time":            job.Time,
			"origin":          job.Origin,
			"fullText":        job.FullText,
			"decodedEntities": job.DecodedEntities,
		})).Info("message received")

		Handle(job)
	}
}

var pathChan = make(chan Path, 100)

func Run(cue string) {
	go worker(pathChan)

	for k := range origins {
		Logger.WithFields(GetLogFields(logrus.Fields{
			"origin": origins[k].GetName(),
		})).Info("starting origin")

		err := origins[k].SetPrefix(RegPrefix)
		if err != nil {
			Logger.WithFields(GetLogFields(logrus.Fields{
				"origin": origins[k].GetName(),
				"error":  err.Error(),
			})).Error("failed to initialize origin")
			continue
		}
		err = origins[k].SetSuffix(RegSuffix)
		if err != nil {
			Logger.WithFields(GetLogFields(logrus.Fields{
				"origin": origins[k].GetName(),
				"error":  err.Error(),
			})).Error("failed to initialize origin")
			continue
		}
		err = origins[k].SetCue(cue)
		if err != nil {
			Logger.WithFields(GetLogFields(logrus.Fields{
				"origin": origins[k].GetName(),
				"error":  err.Error(),
			})).Error("failed to initialize origin")
			continue
		}

		wg.Add(1)
		go origins[k].Start()

		Logger.WithFields(GetLogFields(logrus.Fields{
			"origin":         origins[k].GetName(),
			"maxPayloadSize": origins[k].GetMaxPayloadSize(),
		})).Info("origin started")
	}

	wg.Wait()
}
