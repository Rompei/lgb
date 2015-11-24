package options

// Options : options of this program
type Options struct {
	Width             int     `short:"w" long:"width" description:"Field width" default:"160"`
	Height            int     `short:"h" long:"height" description:"Field height" default:"90"`
	AliveRate         float64 `short:"a" long:"alive-rate" description:"The rate of alive cells in initialization" default:"50"`
	MutRate           float64 `short:"m" long:"mut-rate" description:"The rate of mutation" default:"1.0"`
	Generation        int     `short:"g" long:"generation" description:"Reguration of generations"`
	Keyword           string  `short:"k" long:"keyword" description:"Keyword for twitter"`
	Location          string  `short:"l" long:"location" description:"Tweet location(Default: Japan)" default:"132.2,29.9,146.2,39.0,138.4,33.5,146.1,46.20"`
	Debug             bool    `short:"d" long:"debug" description:"Debug mode"`
	Table             bool    `short:"t" long:"table" description:"View type"`
	Speed             bool    `short:"s" long:"speed" description:"Spead mode"`
	Chaos             bool    `short:"c" long:"chaos" description:"Chaos mode"`
	ConsumerKey       string  `long:"consumer-key" description:"Twitter consumer key"`
	ConsumerSecret    string  `long:"consumer-secret" description:"Twitter consumer secret"`
	AccessToken       string  `long:"access-token" description:"Twitter access token"`
	AccessTokenSecret string  `long:"access-token-secret" description:"Twitter access token secret"`
}
