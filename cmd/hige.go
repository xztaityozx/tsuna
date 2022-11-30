package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strconv"
	"time"
	"tsuna/services/sender"
)

var higeCmd = &cobra.Command{
	Use:  "hige [pid]",
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal().Err(err).Msg("第一引数は数値であるべきです")
		}

		giveUpOnError, _ := cmd.Flags().GetBool("giveUpOnError")
		interval, err := cmd.Flags().GetDuration("interval")
		if err != nil {
			return err
		}
		attempt, err := cmd.Flags().GetInt("attempt")
		if err != nil {
			return err
		}
		sender := sender.NewSender(giveUpOnError, interval, attempt)
		result, err := sender.Send(pid)

		logger := log.Logger.With().Int("pid", result.Process.Pid).Str("シグナル", result.LastSignal.String()).Logger()

		if result.Ok {
			logger.Info().Msg("討伐が完了しました")
		} else {
			logger.Error().Msg("討伐に失敗しました")
		}

		return err
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	higeCmd.Flags().BoolP("giveUpOnError", "f", false, "シグナル送信時に失敗したとき、処理を中断します")
	higeCmd.Flags().DurationP("interval", "i", 1*time.Second, "プロセスが停止したかどうかを調べるまでの時間です")
	higeCmd.Flags().IntP("attempt", "n", 5, "プロセスが停止したかどうかを調べる回数です")
	rootCmd.AddCommand(higeCmd)
}
