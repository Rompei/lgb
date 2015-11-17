package analizer

import (
	"errors"
	"github.com/ikawaha/kagome/tokenizer"
	"gopkg.in/kyokomi/emoji.v1"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// Analizer object
type Analizer struct {
	targets         []string
	analizedTargets []string
	table           [][]string
}

// NewAnalizer : Constructor of Analizer
// @Param targets ツイートリスト
// return Analizer
func NewAnalizer(targets []string) *Analizer {
	return &Analizer{
		targets: targets,
	}
}

// EscapeTargets escape Twitter meta strings
func (a *Analizer) EscapeTargets() error {

	// リプライとURLをエスケープ

	if len(a.targets) == 0 {
		return errors.New("There was not targets")
	}

	re1, err := regexp.Compile(`(^|\s)(@|https?://)\S+`)
	if err != nil {
		return err
	}
	re2, err := regexp.Compile(`^\s*|\s*$`)
	if err != nil {
		return err
	}

	for i, t := range a.targets {
		a.targets[i] = re2.ReplaceAllString(re1.ReplaceAllString(t, ""), "")
		a.targets[i] = strings.Trim(a.targets[i], "\n\r")
	}
	return err
}

// GetTarget returns a target
func (a *Analizer) GetTarget(index int) string {
	return a.targets[index]
}

// ShowTargets shows all strings
func (a *Analizer) ShowTargets() error {

	if len(a.targets) == 0 {
		return errors.New("There was not targets")
	}

	for i, t := range a.targets {
		emoji.Printf(":star:Index: %v, Taeget: %v\n", i, t)
	}

	return nil
}

//ShowAnalizedTargets shows all analized target
func (a *Analizer) ShowAnalizedTargets() error {

	if len(a.targets) == 0 || len(a.analizedTargets) == 0 {
		return errors.New("Targets was not be anaziled.")
	}

	for i, t := range a.analizedTargets {
		emoji.Printf(":star2:Index: %v, Analized target: %v\n", i, t)
	}
	return nil
}

// AnalizeTargets analize target strings
func (a *Analizer) AnalizeTargets() error {

	// 分かち書き

	if len(a.targets) == 0 {
		return errors.New("There are no targets")
	}

	t := tokenizer.New()
	for _, v := range a.targets {
		tokens := t.Tokenize(v)
		for _, token := range tokens {
			if token.Class != tokenizer.DUMMY {
				a.analizedTargets = append(a.analizedTargets, token.Surface)
			}
		}
	}
	return nil
}

// Malcov generates sentences with malcov chain
func (a *Analizer) Malcov() (string, error) {

	// 新文章生成

	if len(a.targets) == 0 || len(a.analizedTargets) == 0 {
		return "", errors.New("Targets was not be anaziled.")
	}

	// テーブル作成
	for i := 0; i < len(a.analizedTargets)-2; i++ {
		j := i
		cell := make([]string, 3)
		for k := 0; k < 3; k++ {
			cell[k] = a.analizedTargets[j]
			j++
		}
		a.table = append(a.table, cell)
	}

	keys := make([]string, 2)
	keys[0] = a.table[0][0]
	keys[1] = a.table[0][1]

	result := keys[0] + keys[1]
	for {
		values, _ := a.findNext(keys)
		value := ""
		if len(values) != 0 {
			value = a.getRandomValue(values)
		}
		if value == "" {
			break
		}
		result += value
		keys[0] = keys[1]
		keys[1] = value
	}

	return result, nil
}

func (a *Analizer) findNext(keys []string) ([]string, error) {
	var results = []string{}
	for _, v := range a.table {
		if v[0] == keys[0] && v[1] == keys[1] {
			results = append(results, v[2])
		}
	}

	if len(results) == 0 {
		return nil, errors.New("Not Matched")
	}

	return results, nil
}

func (a *Analizer) getRandomValue(values []string) string {
	rand.Seed(time.Now().UnixNano())
	return values[rand.Intn(len(values))]
}
