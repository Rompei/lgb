package twitter

import (
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/kyokomi/emoji.v1"
	"net/url"
	"regexp"
)

// Rest is twitter rest api
type Rest struct {
	client *anaconda.TwitterApi
	text   string
}

// NewRest : Constructor of Twitter rest api
// @Param twInfo ツイッター情報
// return Rest
func NewRest(twInfo TwInfo, tweet string) *Rest {
	anaconda.SetConsumerKey(twInfo.ConsumerKey)
	anaconda.SetConsumerSecret(twInfo.ConsumerSecret)
	return &Rest{
		client: anaconda.NewTwitterApi(twInfo.AccessToken, twInfo.AccessTokenSecret),
		text:   tweet,
	}
}

// PostTweet posts a tweet
func (r *Rest) PostTweet() (anaconda.Tweet, error) {
	r.trimTweet()
	emoji.Println(":bird:Posting to Twitter...")
	return r.client.PostTweet(r.text, url.Values{})
}

// ConvertFujiwara converts tweet to Fujiwara Tatsuya
func (r *Rest) ConvertFujiwara() error {
	r.trimTweet()
	re, err := regexp.Compile(`(.)`)
	if err != nil {
		return err
	}
	r.text = re.ReplaceAllString(r.text, "${1}゛")
	return nil
}

func (r *Rest) trimTweet() {
	tweetLength := len(r.text)
	if tweetLength > 140 {
		emoji.Printf(":scissors:Tweet is over 140 strings %v will be trimmed...\n", tweetLength)
		r.text = r.text[:140]
	}
}
