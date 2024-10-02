package nkf

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use: "nkf",
}

func init() {
	Cmd.AddCommand(guess)
	Cmd.AddCommand(euc)
	Cmd.AddCommand(sjis)
	Cmd.AddCommand(utf8)
}

var guess = &cobra.Command{
	Use: "guess",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("nkf: failed to open file: %w", err)
			}
			defer f.Close()
			cmd.SetIn(f)
		}

		nkf := exec.Command("nkf", "--guess")
		nkf.Stdin = cmd.InOrStdin()
		nkf.Stdout = os.Stdout

		if err := nkf.Run(); err != nil {
			return fmt.Errorf("nkf: failed to run nkf: %w", err)
		}

		return nil
	},
}

func convertFunc(nkfArgs []string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		f, err := os.CreateTemp("", "nkf")
		if err != nil {
			return fmt.Errorf("nkf: failed to create temp file: %w", err)
		}
		defer f.Close()

		nkf := exec.Command("nkf", append(nkfArgs, args[0])...)
		nkf.Stdout = f

		if err := nkf.Run(); err != nil {
			return fmt.Errorf("nkf: failed to run nkf: %w", err)
		}

		if err := os.Rename(args[0], args[0]+".bak"); err != nil {
			return fmt.Errorf("nkf: failed to rename original file: %w", err)
		}

		if err := os.Rename(f.Name(), args[0]); err != nil {
			return fmt.Errorf("nkf: failed to rename temp file: %w", err)
		}

		return nil
	}
}

var euc = &cobra.Command{
	Use:  "euc",
	Args: cobra.ExactArgs(1),
	RunE: convertFunc([]string{"-e", "-Lu"}),
}

var sjis = &cobra.Command{
	Use:  "sjis",
	Args: cobra.ExactArgs(1),
	RunE: convertFunc([]string{"-s", "-Lw"}),
}

var utf8 = &cobra.Command{
	Use:  "utf8",
	Args: cobra.ExactArgs(1),
	RunE: convertFunc([]string{"-w", "-Lu"}),
}
