package watanabe

import (
	"fmt"
	"github.com/mattn/go-gimei"
	"math/rand"
	"time"
)

type Metadata struct {
	Last  string
	First string
	yomi  string
}

var watanabeLastNames = []string{"渡辺", "渡邊", "渡邉"}

func (w Metadata) String() string {
	return fmt.Sprintf("%s %s (ワタナベ %s)", w.Last, w.First, w.yomi)
}

func (w Metadata) FullName() string {
	return w.Last + w.First
}

// New は新しい WatanabeMetadata を作って返す
func New() Metadata {
	rand.Seed(time.Now().UnixMicro())
	g := gimei.NewName().First
	return Metadata{yomi: g.Katakana(), First: g.Kanji(), Last: watanabeLastNames[rand.Intn(3)]}
}
