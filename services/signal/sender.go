package signal

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"syscall"
)

type Sender struct {
	logger      *zerolog.Logger
	failOnError bool
}

func NewSender(logger *zerolog.Logger, failOnError bool) *Sender {
	return &Sender{logger: logger, failOnError: failOnError}
}

func (s *Sender) send(proc *os.Process, signal os.Signal) error {
	s.logger.Info().Str("シグナル", signal.String()).Int("pid", proc.Pid).Msg("シグナルを送信しようとしています")
	if err := proc.Signal(signal); err != nil {
		s.logger.Error().Err(err).Int("pid", proc.Pid).Msg("プロセスにシグナルを送信できませんでした")
		return err
	}
	s.logger.Info().Str("シグナル", signal.String()).Int("pid", proc.Pid).Msg("シグナルを送信できました")

	return nil
}

func (s *Sender) Send(pid int) error {
	for _, signal := range sendOrder {
		proc, err := os.FindProcess(pid)
		if err != nil {
			return err
		}
		// 失敗したとき failOnError ならそこで止まる
		if err := s.send(proc, signal); err != nil && s.failOnError {
			return errors.WithMessage(err, "--failOnErrorが指定されたため中断しました")
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
