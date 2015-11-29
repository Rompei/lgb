package twitter

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"gopkg.in/kyokomi/emoji.v1"
	"net/url"
	"regexp"
	"unicode/utf8"
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
func (r *Rest) ConvertFujiwara() (string, error) {
	r.trimTweet()
	re, err := regexp.Compile(`(.)`)
	if err != nil {
		return "", err
	}
	r.text = re.ReplaceAllString(r.text, "${1}゛")
	return r.text, nil
}

// ConvertTNOK makes tweet TNOK
func (r *Rest) ConvertTNOK() string {
	ru := []rune(r.text)
	if utf8.RuneCountInString(r.text) > 65 {
		r.text = string(ru[:65])
	}
	r.text = fmt.Sprintf("%v疲れからか不幸にも、黒塗りの高級車に追突してしまう後輩をかばいすべての責任を負った三浦に対し、車の主、暴力団員谷岡に言い渡された示談の条件とは...。", r.text)
	return r.text

}

func (r *Rest) trimTweet() bool {
	tweetLength := utf8.RuneCountInString(r.text)
	ru := []rune(r.text)
	if tweetLength > 140 {
		emoji.Printf(":scissors:Tweet is over 140 strings %v will be trimmed...\n", tweetLength)
		r.text = string(ru[:140])
		return true
	}
	return false
}
