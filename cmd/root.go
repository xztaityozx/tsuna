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
		fmt.Printf("私は %s です\n", watanabeMetadata)
	},
}

var watanabeMetadata watanabe.Metadata

func init() {
	watanabeMetadata = watanabe.New()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Str("ワタナベ名", watanabeMetadata.String()).Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func Execute() {
	err := rootCmd.Execute()

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("討伐できませんでした")
	}

	log.Info().Msg("このワタナベは役目を終えました")
}
