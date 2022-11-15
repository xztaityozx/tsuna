package cmd

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"xztaityozx/tsuna/models/watanabe"
)

var rootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(watanabeMetadata)
	},
}

var watanabeMetadata watanabe.Watanabe

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	watanabeMetadata = watanabe.New()
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal().
			Err(err).
			Str("ワタナベ名", watanabeMetadata.FullName()).
			Msg("急に死んでしまった！")
	}

	log.Info().Str("ワタナベ名", watanabeMetadata.FullName()).Msg("このワタナベは役目を終えました")
}
