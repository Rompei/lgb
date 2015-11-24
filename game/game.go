package game

import (
	"github.com/Rompei/lgb/field"
	"github.com/Rompei/lgb/options"
	"github.com/Rompei/lgb/twitter"
	"github.com/Rompei/lgb/utils"
	"github.com/cheggaaa/pb"
	"gopkg.in/kyokomi/emoji.v1"
	"net/url"
	"time"
)

// Game object
type Game struct {
	stream     *twitter.Stream
	field      *field.Field
	tweetCh    chan string
	finishFlag bool
	mutRate    float64
	regGen     int
	gameStatus *GameStatus
	debug      bool
	isTable    bool
	isSpeed    bool
	isChaos    bool
}

// GameStatus object
type GameStatus struct {
	Generation int
	DeadNum    int
	BornNum    int
	MutNum     int
}

// NewGame : Constructor of Game
// @Param twInfo ツイッター情報
// return Game
func NewGame(twInfo twitter.TwInfo) *Game {
	return &Game{
		stream:     twitter.NewStream(twInfo),
		gameStatus: &GameStatus{},
	}
}

// InitGame : Initialize field and twitter streaming
func (g *Game) InitGame(opts *options.Options) error {

	// 突然変異確率と世代制限
	g.mutRate = opts.MutRate
	g.regGen = opts.Generation

	// デバッグ
	g.debug = opts.Debug

	// 描画タイプ
	g.isTable = opts.Table

	// スピード
	g.isSpeed = opts.Speed

	// カオスモード
	g.isChaos = opts.Chaos

	// ゲーム情報表示
	g.showGameInfo(opts)
	time.Sleep(3000 * time.Millisecond)

	// Fieldを作る
	var err error
	var aliveNum int
	g.field, aliveNum, err = field.NewField(opts.Width, opts.Height, opts.AliveRate)

	// Twitter Streamingを開始する
	emoji.Println(":santa::santa::santa:Starting Twitter streaming API:christmas_tree::christmas_tree::christmas_tree:")
	v := url.Values{}
	v.Set("track", opts.Keyword)
	if opts.Keyword == "" {
		v.Set("locations", opts.Location)
	}
	g.tweetCh = make(chan string)
	go g.stream.StartPublicFilterStream(v, g.tweetCh)
	g.getInitTweets(aliveNum)

	emoji.Println(":santa::santa::santa:Game was initialized:santa::santa::santa:\n")

	// キャッシュ用キュー稼働
	go g.stream.CollectTweets(g.tweetCh)

	return err
}

// StartGame start main loop of the life game
func (g *Game) StartGame(resultCh chan []string) {
	g.field.Reset()
	g.field.DrawWorld(g.isTable)
	g.field.Reset()
	for {
		// メインループ
		if g.finishFlag || g.regGen != 0 && g.gameStatus.Generation == g.regGen {
			break
		}

		// ステータス表示
		g.initGameStatus()

		// 1ステップすすめる
		g.step()
		g.field.Reset()
		g.showGameStatus()
		g.field.DrawWorld(g.isTable)
		if !g.isSpeed {
			time.Sleep(1500 * time.Millisecond)
		}
	}
	g.StopGame(resultCh)
}

// StopGame stops game and retuns result
func (g *Game) StopGame(resultCh chan []string) {
	g.showLastStatus()
	time.Sleep(1000 * time.Millisecond)
	resultCh <- g.collectTweets()
	g.stream.StopStream()
	close(resultCh)
}

func (g *Game) step() {
	g.finishFlag = true
	nextField := field.InitField(g.field.SizeX, g.field.SizeY)
	for y := 0; y < g.field.SizeY; y++ {
		for x := 0; x < g.field.SizeX; x++ {
			if utils.CheckRate(g.mutRate) {
				// 突然変異
				g.mutation(nextField, x, y)
				g.gameStatus.MutNum++
				g.finishFlag = false
				continue
			}
			aliveCells := g.field.GetAliveCells(x, y)
			if aliveCells == 2 {
				// 変化なし
				if g.debug {
					emoji.Printf(":grin:Points[%v][%v] was not changed\n", x, y)
				}
				nextField.Points[y][x] = g.field.Points[y][x]
			} else if aliveCells == 3 {
				// 創発
				if g.field.Points[y][x].IsAlive {
					if g.debug {
						emoji.Printf(":smiley:Points[%v][%v] is already existed\n", x, y)
					}
					nextField.Points[y][x] = g.field.Points[y][x]
					continue
				}
				tweet := g.getTweets(1)[0]
				if g.debug {
					emoji.Printf(":baby:Points[%v][%v] generated and Fetched tweet: %v\n", x, y, tweet)
				}
				if g.isChaos {
					crossedTweet, err := g.field.CrossParents(x, y, tweet)
					if err != nil {
						panic(err)
					}
					if g.debug {
						emoji.Printf(":baby:Points[%v][%v] generated and Crossed tweet: %v\n", x, y, crossedTweet)
					}
					nextField.AddPoint(x, y, crossedTweet)
					g.gameStatus.BornNum++
					g.finishFlag = false
					continue
				}
				nextField.AddPoint(x, y, tweet)
				g.gameStatus.BornNum++
				g.finishFlag = false
			} else {
				// 絶命
				if !g.field.Points[y][x].IsAlive {
					if g.debug {
						emoji.Printf(":angel:Point[%v][%v] already died\n", x, y)
					}
					continue
				}
				if g.debug {
					emoji.Printf(":dizzy_face:Points[%v][%v] died\n", x, y)
				}
				nextField.DelPoint(x, y)
				g.gameStatus.DeadNum++
				g.finishFlag = false
			}
		}
	}

	// フィールド切り替え
	g.field = nextField

	// 世代更新
	g.gameStatus.Generation++
}

func (g *Game) mutation(nextField *field.Field, x, y int) {
	tweet := g.getTweets(1)[0]
	nextField.AddPoint(x, y, tweet)
	if g.debug {
		emoji.Printf(":boom:Points[%v][%v] generated spontaneously!!! with %v\n", x, y, tweet)
	}
}

func (*Game) showGameInfo(opts *options.Options) {
	emoji.Println(":bell::bell::bell:Game Info:bell::bell::bell::")
	emoji.Printf(":seedling:Initial Alive Rate: %v%%\n", opts.AliveRate)
	emoji.Printf(":family:Reguration of Generations: %v\n", opts.Generation)
	emoji.Printf(":earth_asia:World Size (Width*height): %v*%v\n", opts.Width, opts.Height)
	emoji.Printf(":boom:Mutation Rate: %v%%\n", opts.MutRate)
	emoji.Printf(":globe_with_meridians:Search Keyword: %v\n", opts.Keyword)
	emoji.Printf(":sun_with_face:Tweet locations: %v\n", opts.Location)
	if opts.Chaos {
		emoji.Println(":smiling_imp:Chaos mode ON\n")
	}
}

func (g *Game) initGameStatus() {
	g.gameStatus.DeadNum = 0
	g.gameStatus.BornNum = 0
	g.gameStatus.MutNum = 0
}

func (g *Game) showGameStatus() {
	emoji.Printf(":bamboo::bamboo::bamboo:Game Status of %v:bamboo::bamboo::bamboo:\n", g.gameStatus.Generation)
	emoji.Printf(":innocent:Dead: %v\n", g.gameStatus.DeadNum)
	emoji.Printf(":baby:Born: %v\n", g.gameStatus.BornNum)
	emoji.Printf(":sunglasses:Mutation: %v\n\n", g.gameStatus.MutNum)
}

func (g *Game) showLastStatus() {
	aliveNum := g.field.GetNumberOfAlive()
	emoji.Println(":snowman::snowman::snowman:Last Game Status:snowman::snowman::snowman:")
	emoji.Printf(":family:Generations: %v\n", g.gameStatus.Generation)
	emoji.Printf(":neutral_face:The number of alive cells: %v\n", aliveNum)
	emoji.Printf(":skull:The number of dead cells: %v\n\n", g.field.SizeX*g.field.SizeY-aliveNum)
}

func (g *Game) getTweets(num int) []string {
	results := make([]string, num)

	// キューから取得して足りなければChannelから取得
	if addNum, err := g.stream.GetTweetFromQueue(num, results); err != nil {
		for i := addNum; i < num; i++ {
			results[i] = <-g.tweetCh
		}
		return results
	}
	return results
}

func (g *Game) getInitTweets(aliveNum int) {
	var progress *pb.ProgressBar
	if !g.debug {
		progress = pb.StartNew(aliveNum)
	}
	for y := 0; y < g.field.SizeY; y++ {
		for x := 0; x < g.field.SizeX; x++ {
			if g.field.Points[y][x].IsAlive {
				tweet := <-g.tweetCh
				if g.debug {
					emoji.Printf(":bird:Points[%v][%v]: %v\n", x, y, tweet)
				}
				g.field.Points[y][x].Str = tweet
				if !g.debug {
					progress.Increment()
				}
			}
		}
	}
	if g.debug {
		emoji.Println(":smile::smile::smile:Collected initial tweets:smile::smile::smile:")
	} else {
		e := emoji.Sprint(":smile::smile::smile:")
		progress.FinishPrint(e + "Collected initial tweets" + e)
	}
}

func (g *Game) collectTweets() []string {
	var results []string
	for y := 0; y < g.field.SizeY; y++ {
		for x := 0; x < g.field.SizeX; x++ {
			if g.field.IsAlive(x, y) {
				results = append(results, g.field.Points[y][x].Str)
			}
		}
	}
	return results
}
