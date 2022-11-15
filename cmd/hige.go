package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strconv"
	"xztaityozx/tsuna/services/signal"
)

var higeCmd = &cobra.Command{
	Use:  "hige [pid]",
	Long: `プロセスに SIGHUP => SIGQUIT => SIGINT => SIGKILL の順番でシグナルを送信します`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			log.Fatal().Err(err).Msg("第一引数は数値であるべきです")
		}

		failOnError, _ := cmd.Flags().GetBool("failOnError")
		sender := signal.NewSender(&log.Logger, failOnError)
		return sender.Send(pid)
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	higeCmd.Flags().BoolP("failOnError", "f", false, "シグナル送信時に失敗したとき、処理を中断します")
	rootCmd.AddCommand(higeCmd)
}
