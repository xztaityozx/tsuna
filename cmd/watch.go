package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/mitchellh/go-ps"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os/signal"
	"syscall"
	"time"
	"tsuna/services/sender"
)

var watchCmd = &cobra.Command{
	Use:  "watch",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		executableName := args[0]
		if len(executableName) == 0 {
			return fmt.Errorf("第一引数には監視したいプロセス名を入力を与えるべきですが、実際には空文字が指定されています")
		}

		interval, err := cmd.Flags().GetDuration("interval")
		if err != nil {
			return err
		}
		attempt, err := cmd.Flags().GetInt("attempt")
		if err != nil {
			return err
		}
		workers, err := cmd.Flags().GetInt("workers")
		if err != nil {
			return err
		}

		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
		queue := make(chan int, 1)
		defer close(queue)

		for i := 0; i < workers; i++ {
			go worker(queue, interval, attempt)
		}

		loop := true
		go func() {
			<-ctx.Done()
			log.Info().Msg("Ctrl+Cが押されました。ワタナベたちが撤退を始めています")
			defer stop()
			loop = false
		}()

		log.Info().Msgf("プロセス名が %v なプロセスの検索を始めます", executableName)
		log.Info().Msgf("ワタナベ(ワーカー)の数は %v 人です", workers)
		log.Warn().Msg("Ctrl+CもしくはSIGINTの送信で停止できます")

		for loop {
			processes, err := ps.Processes()
			if err != nil {
				return err
			}

			pidList := lo.Map(
				lo.Filter(
					processes,
					func(p ps.Process, _ int) bool { return p.Executable() == executableName },
				),
				func(p ps.Process, _ int) int { return p.Pid() },
			)

			for _, pid := range pidList {
				log.Info().Int("pid", pid).Msg("プロセスが見つかりました")
				queue <- pid
			}

			time.Sleep(1 * time.Second)
		}

		return errors.New("終了しました")
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func worker(queue <-chan int, interval time.Duration, maxAttempts int) {
	for pid := range queue {
		s := sender.NewSender(true, interval, maxAttempts)
		result, err := s.Send(pid)
		if err != nil {
			s.Logger.Error().Err(err).Int("pid", pid).Msg("このワタナベはプロセス討伐中にエラーを起こしました")
		} else {
			if result.Ok {
				s.Logger.Info().Int("pid", pid).Msg("このワタナベはプロセスを討伐しました")
			} else {
				s.Logger.Error().Int("pid", pid).Msg("このワタナベはプロセスの討伐に失敗しました")
			}
		}
		s.Logger.Info().Msg("このワタナベは役目を終えました")
	}
}

func init() {
	watchCmd.Flags().DurationP("interval", "i", 1*time.Second, "プロセスが停止したかどうかを調べるまでの時間です")
	watchCmd.Flags().IntP("attempt", "n", 5, "プロセスが停止したかどうかを調べる回数です")
	watchCmd.Flags().IntP("workers", "w", 1, "ワーカーの数です")
	rootCmd.AddCommand(watchCmd)
}
