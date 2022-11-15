package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("hello world")
  },
}

func Execute() {
  err := rootCmd.Execute()

  if err != nil {
    log.Fatalln(err)
  }
}
