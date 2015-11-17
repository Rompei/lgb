package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/kyokomi/emoji.v1"
	"net/url"
)

// Rest is twitter rest api
type Rest struct {
	client *anaconda.TwitterApi
}

// NewRest : Constructor of Twitter rest api
// @Param twInfo ツイッター情報
// return Rest
func NewRest(twInfo TwInfo) *Rest {
	anaconda.SetConsumerKey(twInfo.ConsumerKey)
	anaconda.SetConsumerSecret(twInfo.ConsumerSecret)
	return &Rest{
		client: anaconda.NewTwitterApi(twInfo.AccessToken, twInfo.AccessTokenSecret),
	}
}

// PostTweet posts a tweet
func (r *Rest) PostTweet(text string) (anaconda.Tweet, error) {
	tweet := r.trimTweet(text)
	emoji.Println(":bird:Posting to Twitter...")
	return r.client.PostTweet(tweet, url.Values{})
}

func (r *Rest) trimTweet(original string) string {
	//140字より大きければ縮める
	tweetLen := len(original)
	if tweetLen > 140 {
		emoji.Printf(":scream:Length was over 140 Length=%v\n", tweetLen)
		emoji.Println(":scissors:Tweet will be trimmed\n")
		return original[140:]
	}
	return original
}
