package main

import (
	"bufio"
	"fmt"
	"github.com/Rompei/lgb/analyzer"
	"github.com/Rompei/lgb/game"
	"github.com/Rompei/lgb/options"
	"github.com/Rompei/lgb/twitter"
	flags "github.com/jessevdk/go-flags"
	"gopkg.in/kyokomi/emoji.v1"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	opts := checkArgs()

	info := twitter.TwInfo{
		ConsumerKey:       opts.ConsumerKey,
		ConsumerSecret:    opts.ConsumerSecret,
		AccessToken:       opts.AccessToken,
		AccessTokenSecret: opts.AccessTokenSecret,
	}

	game := game.NewGame(info)
	err := game.InitGame(opts)
	checkError(err)
	resultCh := make(chan []string, 1024)

	// ライフゲーム開始
	game.StartGame(resultCh)

	// ゲームオーバー
	results := <-resultCh
	emoji.Println(":beer::beer::beer:Game Over:beer::beer::beer:\n")
	if len(results) == 0 {
		// 絶滅
		emoji.Println(":smiling_imp::smiling_imp::smiling_imp:All life was exterminated:smiling_imp::smiling_imp::smiling_imp:")
		os.Exit(0)
	}

	// 生き残ったツイートを処理する
	analyzer := analyzer.NewAnalyzer(results)
	err = analyzer.EscapeTargets()
	if opts.Debug {
		err = analyzer.ShowTargets()
	}
	err = analyzer.AnalizeTargets()
	if opts.Debug {
		err = analyzer.ShowAnalizedTargets()
	}
	checkError(err)
	newText, err := analyzer.Malcov()
	checkError(err)

	emoji.Printf("\n:beer:Generated Text: \n%v\n\n", newText)

	// ツイートするか尋ねる
	isDoTweet, err := tellYesNo("Would you like to tweet?> ")
	checkError(err)
	if !isDoTweet {
		os.Exit(0)
	}

	// ツイートする
	rest := twitter.NewRest(info)
	result, err := rest.PostTweet(newText)
	checkError(err)
	tweetURL := fmt.Sprintf("https://twitter.com/%v/status/%v", result.User.IdStr, result.IdStr)
	emoji.Printf(":beers:Posted the Tweet. You can see from: %v\n", tweetURL)
}

// コマンドライン引数チェック用関数
func checkArgs() *options.Options {
	var opts options.Options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "lgb"
	parser.Usage = "[OPTIONS]"
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.ConsumerKey == "" || opts.ConsumerSecret == "" || opts.AccessToken == "" || opts.AccessTokenSecret == "" {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	if opts.Width < 1 || opts.Height < 1 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	if opts.AliveRate < 1 || opts.AliveRate > 100 || opts.MutRate < 0 || opts.MutRate > 100 || opts.Generation < 0 {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}

	return &opts
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func tellYesNo(text string) (bool, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf(text)
	for scanner.Scan() {
		if scanner.Text() == "Y" || scanner.Text() == "y" || scanner.Text() == "yes" || scanner.Text() == "Yes" {
			return true, nil
		} else if scanner.Text() == "N" || scanner.Text() == "n" || scanner.Text() == "no" || scanner.Text() == "No" {
			return false, nil
		} else {
			fmt.Printf(text)
		}
	}
	err := scanner.Err()
	return false, err
}
