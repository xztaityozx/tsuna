package sender

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"syscall"
	"time"
	"tsuna/models/watanabe"
)

type Sender struct {
	failOnError   bool
	watchInterval time.Duration
	watchAttempt  int
	watanabe.Watanabe
}

func NewSender(failOnError bool, interval time.Duration, attempt int) *Sender {
	return &Sender{
		failOnError:   failOnError,
		watchInterval: interval,
		watchAttempt:  attempt,
		Watanabe:      watanabe.New(),
	}
}

type SendResult struct {
	Process    *os.Process
	SentSignal os.Signal
	Ok         bool
}

func (s *Sender) Send(pid int) (SendResult, error) {
	proc, err := os.FindProcess(pid)
	logger := s.Logger.With().Int("pid", proc.Pid).Logger()
	for _, sig := range sendOrder {
		if err != nil {
			return SendResult{proc, sig, false}, err
		}

		logger.Info().Str("シグナル", sig.String()).Msg("シグナルを送信しようとしています")
		err = proc.Signal(sig)
		if err != nil {
			log.Error().Err(err).Int("pid", proc.Pid).Msg("プロセスにシグナルを送信できませんでした")
			if s.failOnError {
				return SendResult{proc, sig, false}, errors.WithMessage(err, "中断しました")
			}
			continue
		}

		logger.Info().Str("シグナル", sig.String()).Msg("シグナルを送信できました")
		for i := 0; i < s.watchAttempt; i++ {
			log.Info().Int("監視回数", i+1).Int("最大監視回数", s.watchAttempt).Msg("プロセスの終了を監視しています...")
			time.Sleep(s.watchInterval)
			if err := proc.Signal(syscall.Signal(0)); err != nil {
				log.Info().Int("監視回数", i+1).Msg("プロセスが終了したようです")
				return SendResult{proc, sig, true}, nil
			}
		}
	}

	logger.Warn().Msg("監視中にプロセスが終了しませんでした")
	return SendResult{nil, nil, false}, errors.New("いくつかのシグナルを送信しましたが、終了できませんでした")
}

var sendOrder = []os.Signal{
	syscall.SIGQUIT,
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGKILL,
}
