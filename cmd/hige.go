package cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"xztaityozx/tsuna/services/signal"
)

var higeCmd = &cobra.Command{
	Use:  "hige",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			log.Panic().Err(err).Msg("第一引数は数値であるべきです")
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			log.Panic().Err(err).Msg("討伐対象のプロセスを見つけられませんでした")
		}
		sender := signal.NewSender(&log.Logger)
		if err := sender.Send(proc); err != nil {
			log.Panic().Err(err).Msg("討伐できませんでした")
		}
	},
}

func init() {
	rootCmd.AddCommand(higeCmd)
}
