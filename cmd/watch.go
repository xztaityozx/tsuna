package cmd

import (
	"github.com/mitchellh/go-ps"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"sync"
	"xztaityozx/tsuna/services/signal"
)

var watchCmd = &cobra.Command{
	Use: "watch",
	RunE: func(cmd *cobra.Command, args []string) error {
		processes, err := ps.Processes()
		if err != nil {
			return err
		}

		sleepProcesses := lo.Filter(processes, func(process ps.Process, i int) bool { return process.Executable() == "sleep" })
		queue := make(chan int, 1)
		defer close(queue)
		var wg sync.WaitGroup
		go func() {
			wg.Add(1)
			sender := signal.NewSender(&log.Logger, true)
			for pid := range queue {
				if err := sender.Send(pid); err != nil {
					log.Err(err).Send()
				}
			}
		}()

		for _, process := range sleepProcesses {
			queue <- process.Pid()
		}

		wg.Wait()
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.AddCommand(watchCmd)
}
