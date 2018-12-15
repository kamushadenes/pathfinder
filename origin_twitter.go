package pathfinder

import (
	"github.com/caarlos0/env"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const twitterMaxSize = 280

type Twitter struct {
	ConsumerKey    string `env:"TWITTER_CONSUMER_KEY,required"`
	ConsumerSecret string `env:"TWITTER_CONSUMER_SECRET,required"`
	AccessToken    string `env:"TWITTER_ACCESS_TOKEN,required"`
	AccessSecret   string `env:"TWITTER_ACCESS_SECRET,required"`
	httpClient     *http.Client
	client         *twitter.Client
	prefix         string
	suffix         string
	cue            string
}

func NewTwitter() (*Twitter, error) {
	var t Twitter
	err := env.Parse(&t)

	if err != nil {
		return nil, err
	}

	config := oauth1.NewConfig(t.ConsumerKey, t.ConsumerSecret)
	token := oauth1.NewToken(t.AccessToken, t.AccessSecret)
	t.httpClient = config.Client(oauth1.NoContext, token)

	// Twitter client
	t.client = twitter.NewClient(t.httpClient)

	return &t, nil
}

func (t Twitter) GetName() string {
	return "twitter"
}

func (t *Twitter) FindPath() *Path {
	return nil
}

func (t Twitter) GetMaxPayloadSize() int {
	return twitterMaxSize - len(t.prefix) - len(t.suffix)
}

func (t Twitter) GetPrefix() string {
	return t.prefix
}

func (t *Twitter) SetPrefix(s string) error {
	size := len(s) + len(t.suffix)
	if size < (twitterMaxSize - minimumAvailablePayloadSize) {
		t.prefix = s
		return nil
	} else {
		return NewMaxPayloadSizeError(size, twitterMaxSize)
	}
}

func (t Twitter) GetSuffix() string {
	return t.suffix
}

func (t *Twitter) SetSuffix(s string) error {
	size := len(s) + len(t.prefix)
	if size < (twitterMaxSize - minimumAvailablePayloadSize) {
		t.suffix = s
		return nil
	} else {
		return NewMaxPayloadSizeError(size, twitterMaxSize)
	}
}

func (t Twitter) GetCue() string {
	return t.cue
}

func (t *Twitter) SetCue(s string) error {
	t.cue = s
	return nil
}

func (t *Twitter) Start() {
	params := &twitter.StreamFilterParams{
		Track:         []string{t.cue},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := t.client.Streams.Filter(params)

	if err != nil {
		Logger.WithFields(GetLogFields(logrus.Fields{
			"origin": t.GetName(),
			"error":  err.Error(),
		})).Error("failed to start origin")
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		p := NewPath()
		ti, err := tweet.CreatedAtTime()
		if err != nil {
			now := time.Now()
			p.Time = &now
		} else {
			p.Time = &ti
		}
		p.Origin = t.GetName()
		p.OriginPrefix = t.GetPrefix()
		p.OriginSuffix = t.GetSuffix()

		if tweet.ExtendedTweet != nil {
			p.FullText = tweet.ExtendedTweet.FullText
		} else {
			p.FullText = tweet.Text
		}

		p.Metadata["username"] = tweet.User.Name

		pathChan <- *p
	}

	go demux.HandleChan(stream.Messages)

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	stream.Stop()

	wg.Done()
}

func init() {
	t, err := NewTwitter()

	if err != nil {
		Logger.WithFields(GetLogFields(logrus.Fields{
			"origin": t.GetName(),
			"error":  err.Error(),
		})).Error("failed to start origin")
	}

	t.cue = os.Args[1]

	origins = append(origins, t)
}
