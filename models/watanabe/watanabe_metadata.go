package watanabe

import (
	"fmt"
	"github.com/mattn/go-gimei"
	"math/rand"
	"time"
)

type Watanabe struct {
	Last  string
	First string
	yomi  string
}

var watanabeLastNames = []string{"渡辺", "渡邊", "渡邉"}

func (w Watanabe) String() string {
	return fmt.Sprintf("私は %s %s (ワタナベ %s) です", w.Last, w.First, w.yomi)
}

func (w Watanabe) FullName() string {
	return w.Last + w.First
}

// New は新しい WatanabeMetadata を作って返す
func New() Watanabe {
	rand.Seed(time.Now().UnixMicro())
	g := gimei.NewName().First
	return Watanabe{yomi: g.Katakana(), First: g.Kanji(), Last: watanabeLastNames[rand.Intn(3)]}
}
