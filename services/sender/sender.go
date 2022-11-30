package sender

import (
	"os"
	"syscall"
	"time"
	"tsuna/models/watanabe"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Sender はプロセスに順番にシグナルを送るための構造体
type Sender struct {
	// シグナルの送信に失敗した段階で後続のシグナル送信諦めるかどうか
	giveUpOnError bool
	// プロセスの生存確認インターバル
	watchInterval time.Duration
	// プロセスの生存確認の回数
	watchAttempt int
	// extends watanabe.Watanabe 。カスタマイズされた zerolog.Logger を使うため
	watanabe.Watanabe
}

// NewSender は sender.Sender を作って返す
func NewSender(giveUpOnError bool, interval time.Duration, attempt int) *Sender {
	return &Sender{
		giveUpOnError: giveUpOnError,
		watchInterval: interval,
		watchAttempt:  attempt,
		// watanabe.Watanabe はここでテキトウに作っておく
		Watanabe: watanabe.New(),
	}
}

// SendResult はシグナル送信の結果を表す構造体
type SendResult struct {
	// 送信先のプロセス
	Process *os.Process
	// 最後に送ったシグナル
	LastSignal os.Signal
	// 停止を確認できたかどうか
	Ok bool
}

// Send はプロセスにシグナルを順番に送る。順番は sendOrder の通り
func (s *Sender) Send(pid int) (SendResult, error) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return SendResult{proc, nil, false}, err
	}

	logger := s.Logger.With().Int("pid", proc.Pid).Logger()
	for _, sig := range sendOrder {
		logger.Info().Str("シグナル", sig.String()).Msg("シグナルを送信しようとしています")
		if err := proc.Signal(sig); err != nil {
			log.Error().Err(err).Int("pid", proc.Pid).Msg("プロセスにシグナルを送信できませんでした")
			if s.giveUpOnError {
				return SendResult{proc, sig, false}, errors.WithMessage(err, "中断しました")
			}
			continue
		}

		logger.Info().Str("シグナル", sig.String()).Msg("シグナルを送信できました")
		// 生存確認ここから
		// watchAttempt * watchInterval 秒待機することになる
		for i := 0; i < s.watchAttempt; i++ {
			logger.Info().Int("監視回数", i+1).Int("最大監視回数", s.watchAttempt).Msg("プロセスの終了を監視しています...")
			time.Sleep(s.watchInterval)
			if err := proc.Signal(syscall.Signal(0)); err != nil {
				logger.Info().Int("監視回数", i+1).Msg("プロセスの停止を確認しました")
				return SendResult{proc, sig, true}, nil
			}
		}
	}

	logger.Warn().Msg("監視中にプロセスが終了しませんでした")
	return SendResult{nil, nil, false}, errors.New("いくつかのシグナルを送信しましたが、終了できませんでした")
}

var sendOrder = []os.Signal{
	syscall.SIGQUIT,
	syscall.SIGINT,
	syscall.SIGHUP,
	syscall.SIGPIPE,
	syscall.SIGSEGV,
	syscall.SIGALRM,
	syscall.SIGTERM,
	syscall.SIGKILL,
}
