package twitter

import (
	"errors"
	"github.com/ChimeraCoder/anaconda"
	gq "github.com/DamnWidget/goqueue"
	"net/url"
)

// QueueLimit : Limit of Tweet's queue
const QueueLimit = 1024

// Stream is twitter streaming api
type Stream struct {
	client     *anaconda.TwitterApi
	streamAPI  *anaconda.Stream
	tweetQueue *gq.Queue
}

// NewStream : constructor of Stream
// @Param twInfo ツイッター情報
// return Stream
func NewStream(twInfo TwInfo) *Stream {
	anaconda.SetConsumerKey(twInfo.ConsumerKey)
	anaconda.SetConsumerSecret(twInfo.ConsumerSecret)
	return &Stream{
		client:     anaconda.NewTwitterApi(twInfo.AccessToken, twInfo.AccessTokenSecret),
		tweetQueue: gq.New(QueueLimit),
	}
}

// StartPublicFilterStream start stream api
func (s *Stream) StartPublicFilterStream(v url.Values, tweetCh chan string) {
	s.streamAPI = s.client.PublicStreamFilter(v)

	for {
		item := <-s.streamAPI.C
		switch status := item.(type) {
		case anaconda.Tweet:
			tweetCh <- status.Text
		default:
		}
	}
}

// StartPublicSampleStream start stream api
func (s *Stream) StartPublicSampleStream(v url.Values, tweetCh chan string) {
	s.streamAPI = s.client.PublicStreamSample(v)

	for {
		item := <-s.streamAPI.C
		switch status := item.(type) {
		case anaconda.Tweet:
			tweetCh <- status.Text
		default:
		}
	}
}

// CollectTweets collects tweet to s.tweetQueue
func (s *Stream) CollectTweets(tweetCh chan string) {
	for {
		select {
		case tweet := <-tweetCh:
			if err := s.tweetQueue.Push(tweet); err != nil {
				s.tweetQueue.Pop()
				err := s.tweetQueue.Push(tweet)
				if err != nil {
					panic(err)
				}
			}
		default:
		}
	}
}

// GetTweetFromQueue returns tweets required. If tweets was not enough, returns error.
func (s *Stream) GetTweetFromQueue(num int, results []string) (int, error) {
	for i := 0; i < num; i++ {
		switch item := s.tweetQueue.Pop().(type) {
		case string:
			results[i] = item
		default:
			return i, errors.New("Tweet was not enouth")
		}
	}
	return num - 1, nil
}

// StopStream stops streaming
func (s *Stream) StopStream() {
	s.streamAPI.Stop()
}
