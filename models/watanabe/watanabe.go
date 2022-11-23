package watanabe

import (
	"fmt"
	"github.com/mattn/go-gimei"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

type Watanabe struct {
	Logger *zerolog.Logger
	Last   string
	First  string
	Yomi   string
}

var watanabeLastNames = []string{"渡辺", "渡邊", "渡邉"}

func (w Watanabe) String() string {
	return fmt.Sprintf("%s %s (ワタナベ %s)", w.Last, w.First, w.Yomi)
}

func (w Watanabe) FullName() string {
	return w.Last + w.First
}

func New() Watanabe {
	g := gimei.NewName().First
	w := Watanabe{
		Last:  watanabeLastNames[rand.Intn(3)],
		First: g.Kanji(),
		Yomi:  g.Katakana(),
	}
	logger := log.Logger.With().Str("ワタナベ名", w.String()).Logger()
	w.Logger = &logger

	return w
}

func init() {
	rand.Seed(time.Now().UnixMicro())
}
