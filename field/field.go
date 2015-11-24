package field

import (
	"errors"
	"fmt"
	"github.com/Rompei/lgb/point"
	"github.com/Rompei/lgb/utils"
	"github.com/ikawaha/kagome/tokenizer"
	"github.com/olekukonko/tablewriter"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"
)

// Field object
type Field struct {
	Points [][]point.Point
	SizeX  int
	SizeY  int
}

// NewField : Constructor of Field
// @Param size      フィールドサイズ
// @Param aliveRate 初期化時の生存率
// @Param strs      文字
// return Field
func NewField(sizeX, sizeY int, aliveRate float64) (*Field, int, error) {
	aliveNum := 0
	f := initField(sizeX, sizeY)
	f.Points = make([][]point.Point, sizeY)
	for y := 0; y < sizeY; y++ {
		cells := make([]point.Point, sizeX)
		for x := 0; x < sizeX; x++ {
			isAlive := utils.CheckRate(aliveRate)
			if isAlive {
				aliveNum++
			}
			cells[x] = point.NewPoint(x, y, "", isAlive)
		}
		f.Points[y] = cells
	}
	return f, aliveNum, nil
}

// InitField initialize empty field
func InitField(sizeX, sizeY int) *Field {
	f := initField(sizeX, sizeY)

	f.Points = make([][]point.Point, sizeY)
	for y := 0; y < sizeY; y++ {
		cells := make([]point.Point, sizeX)
		for x := 0; x < sizeX; x++ {
			cells[x] = point.NewPoint(x, y, "", false)
		}
		f.Points[y] = cells
	}
	return f
}

func initField(sizeX, sizeY int) *Field {
	return &Field{
		SizeX: sizeX,
		SizeY: sizeY,
	}
}

// AddPoint add point to empty field
func (f *Field) AddPoint(x, y int, str string) error {
	if f.Points[y][x].IsAlive {
		return errors.New("Life already existed")
	}
	f.Points[y][x].Str = str
	f.Points[y][x].IsAlive = true
	return nil
}

// AddRandomPoint add point for field
func (f *Field) AddRandomPoint(str string) {
	isFind := false
	for !isFind {
		point := getRandomPoint(f.SizeX, f.SizeY)
		if !f.Points[point.Y][point.X].IsAlive {
			f.AddPoint(point.X, point.Y, str)
			isFind = true
		}
	}
}

// DelPoint kill the point
func (f *Field) DelPoint(x, y int) {
	f.Points[y][x].IsAlive = false
}

// IsAlive retuens whether the cell alive or not
func (f *Field) IsAlive(x, y int) bool {
	return f.Points[y][x].IsAlive
}

// GetNumberOfAlive returns the nubmer of alive cells
func (f *Field) GetNumberOfAlive() int {
	sum := 0
	for y := 0; y < f.SizeY; y++ {
		for x := 0; x < f.SizeY; x++ {
			if f.IsAlive(x, y) {
				sum++
			}
		}
	}
	return sum
}

// DrawWorld draws world
func (f *Field) DrawWorld(isTable bool) {

	// オプションにより切り替え
	if isTable {
		f.drawTable()
	} else {
		f.drawPoint()
	}
}

// DrawPoint draws point
func (f *Field) drawPoint() {
	for y := 0; y < f.SizeY; y++ {
		for x := 0; x < f.SizeX; x++ {
			if f.Points[y][x].IsAlive {
				//fmt.Printf("%v:%v is alive ", x, y)
				fmt.Printf("%v", string([]rune(f.Points[y][x].Str)[0]))
			} else {
				//fmt.Printf("%v:%v is dead ", x, y)
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}

// DrawTable draws table of the world
func (f *Field) drawTable() {
	firstRow := make([]string, f.SizeX+1)
	firstRow[0] = ""
	for x := 0; x < f.SizeX; x++ {
		index := strconv.FormatInt(int64(x), 10)
		firstRow[x+1] = index
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.Append(firstRow)
	for y := 0; y < f.SizeY; y++ {
		col := make([]string, f.SizeX+1)
		col[0] = strconv.FormatInt(int64(y), 10)
		for x := 0; x < f.SizeX; x++ {
			if f.Points[y][x].IsAlive {
				col[x+1] = string([]rune(f.Points[y][x].Str)[0])
			} else {
				col[x+1] = ""
			}
		}
		table.Append(col)
	}
	table.Render()

}

// Reset screen
func (f *Field) Reset() {
	for y := 0; y < f.SizeY; y++ {
		fmt.Println("")
	}
}

// GetAliveCells returns the number of neighbor cell
func (f *Field) GetAliveCells(x, y int) int {
	if x == 0 || x == f.SizeX-1 || y == 0 || y == f.SizeY-1 {
		return 0
	}

	sum := 0
	if f.Points[y-1][x-1].IsAlive {
		sum++
	}
	if f.Points[y-1][x].IsAlive {
		sum++
	}
	if f.Points[y-1][x+1].IsAlive {
		sum++
	}
	if f.Points[y][x-1].IsAlive {
		sum++
	}
	if f.Points[y][x+1].IsAlive {
		sum++
	}
	if f.Points[y+1][x-1].IsAlive {
		sum++
	}
	if f.Points[y+1][x].IsAlive {
		sum++
	}
	if f.Points[y+1][x+1].IsAlive {
		sum++
	}
	return sum

}

// CrossParents cross parent cell's tweet
func (f *Field) CrossParents(x, y int, newTweet string) (string, error) {
	if x == 0 || x == f.SizeX-1 || y == 0 || y == f.SizeY-1 {
		return "", errors.New("Invalid cell.")
	}

	var tweets []string
	if f.Points[y-1][x-1].IsAlive {
		tweets = append(tweets, f.Points[y-1][x-1].Str)
	}
	if f.Points[y-1][x].IsAlive {
		tweets = append(tweets, f.Points[y-1][x].Str)
	}
	if f.Points[y-1][x+1].IsAlive {
		tweets = append(tweets, f.Points[y-1][x+1].Str)
	}
	if f.Points[y][x-1].IsAlive {
		tweets = append(tweets, f.Points[y][x-1].Str)
	}
	if f.Points[y][x+1].IsAlive {
		tweets = append(tweets, f.Points[y][x+1].Str)
	}
	if f.Points[y+1][x-1].IsAlive {
		tweets = append(tweets, f.Points[y+1][x-1].Str)
	}
	if f.Points[y+1][x].IsAlive {
		tweets = append(tweets, f.Points[y+1][x].Str)
	}
	if f.Points[y+1][x+1].IsAlive {
		tweets = append(tweets, f.Points[y+1][x+1].Str)
	}

	re1, err := regexp.Compile(`(^|\s)(@|https?://)\S+`)
	if err != nil {
		return "", err
	}
	re2, err := regexp.Compile(`^\s*|\s*$`)
	if err != nil {
		return "", err
	}

	t := tokenizer.New()
	newTweet = re2.ReplaceAllString(re1.ReplaceAllString(newTweet, ""), "")
	originTweetTokens := t.Tokenize(newTweet)
	var parentTweetsTokens [][]tokenizer.Token
	for i, v := range tweets {
		tweets[i] = re2.ReplaceAllString(re1.ReplaceAllString(v, ""), "")
		parentTweetsTokens = append(parentTweetsTokens, t.Tokenize(tweets[i]))
	}

	parentPtr := 0
	for i, ot := range originTweetTokens {
		if ot.Class == tokenizer.DUMMY {
			continue
		}
		for _, t := range parentTweetsTokens[parentPtr] {
			if t.Class != tokenizer.DUMMY && ot.Features()[0] == t.Features()[0] && ot.Features()[1] == t.Features()[1] {
				originTweetTokens[i] = t
				parentPtr++
				break
			}
		}
		if parentPtr == len(parentTweetsTokens)-1 {
			parentPtr = 0
		}
	}

	generatedTweet := ""
	for _, t := range originTweetTokens {
		if t.Class != tokenizer.DUMMY {
			generatedTweet += t.Surface
		}
	}
	if generatedTweet == "" {
		generatedTweet = "からの"
	}

	return generatedTweet, nil
}

//ShowFieldInfo : Debug function
func (f *Field) ShowFieldInfo() {
	fmt.Printf("Field width: %v\n", len(f.Points[0]))
	fmt.Printf("Field height: %v\n", len(f.Points))
}

func getRandomPoint(sizeX, sizeY int) point.Point {
	rand.Seed(time.Now().UnixNano())
	return point.Point{
		X:       rand.Intn(sizeX),
		Y:       rand.Intn(sizeY),
		Str:     "",
		IsAlive: true,
	}
}
