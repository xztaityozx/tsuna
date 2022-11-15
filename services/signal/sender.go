package signal

import (
	"github.com/rs/zerolog"
	"os"
	"syscall"
	"xztaityozx/tsuna/models/watanabe"
)

type Sender struct {
	logger *zerolog.Logger
	watanabe.Watanabe
}

func NewSender(logger *zerolog.Logger) *Sender {
	return &Sender{logger: logger}
}

func (s *Sender) Send(proc *os.Process) error {
	for _, signal := range sendOrder {
		if err := proc.Signal(signal); err != nil {
			s.logger.Error().Err(err).Str("ワタナベ名", s.First).Int("pid", proc.Pid).Msg("プロセスにシグナルを送信できませんでした")
		}
	}

	return nil
}

var sendOrder = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGQUIT,
	syscall.SIGINT,
	syscall.SIGKILL,
}
