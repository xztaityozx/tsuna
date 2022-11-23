package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"tsuna/models/watanabe"
)

var jsonLogging = false
var rootCmd = &cobra.Command{
	Long:    `プロセスに SIGQUIT => SIGHUP => SIGINT => SIGTERM => SIGKILLの順番でシグナルを送信します。いろんな方法で送信します`,
	Version: "0.0.1",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		if !jsonLogging {
			// jsonLogじゃないときはConsoleLoggerにする
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Logger()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("私は %s です\n", watanabe.New())
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonLogging, "json-log", false, "ログ出力をJSON形式にします")
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal().
			Err(err).Send()
	}
}
