package cmd

import (
	"fmt"
	"os"
	"tsuna/models/watanabe"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var jsonLogging = false
var debugOn = false
var rootCmd = &cobra.Command{
	Long:    `プロセスに SIGQUIT => SIGHUP => SIGINT => SIGTERM => SIGKILLの順番でシグナルを送信します。いろんな方法で送信します`,
	Version: "0.0.1",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		if !jsonLogging {
			// jsonLogじゃないときはConsoleLoggerにする
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Logger()
		}

		if debugOn {
			log.Logger = log.Logger.Level(zerolog.DebugLevel)
		} else {
			log.Logger = log.Logger.Level(zerolog.InfoLevel)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("私は %s です\n", watanabe.New())
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonLogging, "json-log", false, "ログ出力をJSON形式にします")
	rootCmd.PersistentFlags().BoolVar(&debugOn, "debug", false, "デバッグ出力をONにします")
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal().
			Err(err).Send()
	}
}
