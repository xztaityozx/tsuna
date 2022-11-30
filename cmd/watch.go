package cmd

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"tsuna/services/sender"

	"github.com/mitchellh/go-ps"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
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

		// Ctrl+Cでctx.Done()になるので、各ワーカーはselectでそれを待ち受けつつ処理する
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
		defer stop()
		queue := make(chan int, 1)

		var wg sync.WaitGroup
		wg.Add(workers)
		for i := 0; i < workers; i++ {
			go func() {
				worker(ctx, queue, func() *sender.Sender {
					// watchで使うsender.Senderはシグナルの送信に失敗したら諦めるようにする
					return sender.NewSender(true, interval, attempt)
				}, &wg)
			}()
		}

		log.Info().Msgf("プロセス名が %v なプロセスの検索を始めます", executableName)
		log.Info().Msgf("ワタナベ(ワーカー)の数は %v 人です", workers)
		log.Warn().Msg("Ctrl+CもしくはSIGINTの送信で停止できます")

		wg.Add(1)
		go func() {
			defer wg.Done()
			<-ctx.Done()
			log.Info().Msg("Ctrl+Cが押されました。ワタナベたちが撤退を始めています")
			// dispatcherが詰めようとしているPIDを消費しないとdispatcherがキャンセルされないのでログはいて終わる
			for pid := range queue {
				log.Info().Int("pid", pid).Msg("プロセスは補足されましたが討伐されませんでした")
			}
		}()

		wg.Add(1)
		go func() {
			// dispatcherの起動ここ
			if err := dispatcher(ctx, queue, func(p ps.Process, _ int) bool { return p.Executable() == executableName }, &wg); err != nil {
				log.Error().Err(err).Msg("プロセスの探索をしていたワタナベがエラーを返しました。これ以上プロセスを探索できません")
			}

			log.Info().Msg("プロセスを探索していたワタナベが撤退しました")
		}()

		wg.Wait()

		return errors.New("Watchを終了しました")
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

// dispatcher はProcessを探して queue に詰めるワーカー
func dispatcher(ctx context.Context, queue chan<- int, filterFunc func(process ps.Process, i int) bool, wg *sync.WaitGroup) error {
	dict := map[int]struct{}{}
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			close(queue)
			return nil
		default:
			processes, err := ps.Processes()
			if err != nil {
				return err
			}

			for _, pid := range lo.Uniq(lo.Map(lo.Filter(processes, filterFunc), func(p ps.Process, _ int) int { return p.Pid() })) {
				// 一回見つけたプロセスは送信しない
				if _, ok := dict[pid]; ok {
					log.Info().Int("pid", pid).Msg("プロセスを見つけましたが、すでに補足したものであったためスキップしました")
					continue
				} else {
					log.Info().Int("pid", pid).Msg("プロセスが見つかりました")
					dict[pid] = struct{}{}
				}
				queue <- pid
			}

			time.Sleep(1 * time.Second)
		}
	}
}

// worker は queue からPIDを受け取ってシグナルを送るワーカー
func worker(ctx context.Context, queue <-chan int, senderFactoryFunc func() *sender.Sender, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case pid := <-queue:
			s := senderFactoryFunc()
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
}

func init() {
	watchCmd.Flags().DurationP("interval", "i", 1*time.Second, "プロセスが停止したかどうかを調べるまでの時間です")
	watchCmd.Flags().IntP("attempt", "n", 5, "プロセスが停止したかどうかを調べる回数です")
	watchCmd.Flags().IntP("workers", "w", 1, "ワーカーの数です")
	rootCmd.AddCommand(watchCmd)
}
