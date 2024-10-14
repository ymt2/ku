package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ymt2/ku/internal/llm"
	"github.com/ymt2/ku/internal/nkf"
)

var rootCmd = &cobra.Command{
	Use: "ku",
}

func init() {
	rootCmd.AddCommand(llm.Cmd)
	rootCmd.AddCommand(nkf.Cmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
