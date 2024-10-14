package llm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/ymt2/ku/internal/llm/perplexity"
)

const PERPLEXITY_API_TOKEN_KEY = "PERPLEXITY_API_TOKEN"

var Cmd = &cobra.Command{
	Use: "llm",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			cmd.SetIn(strings.NewReader(args[0]))
		}

		token := os.Getenv(PERPLEXITY_API_TOKEN_KEY)
		if token == "" {
			return fmt.Errorf("llm: %s is not set", PERPLEXITY_API_TOKEN_KEY)
		}

		cli, err := perplexity.NewClient(token)
		if err != nil {
			return fmt.Errorf("llm: failed to create client: %w", err)
		}

		return chatCompletions(cli, cmd.InOrStdin())
	},
}

var interactive bool

func init() {
	Cmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Continue the conversation")
}

func chatCompletions(cli *perplexity.Client, in io.Reader) error {
	content, err := io.ReadAll(in)
	if err != nil {
		return fmt.Errorf("llm: failed to read input: %w", err)
	}

	histories := []perplexity.ChatCompletionMessage{
		{
			Role:    perplexity.RoleSystem,
			Content: "Be precise and concise.",
		},
		{
			Role:    perplexity.RoleUser,
			Content: string(content),
		},
	}

	for {
		req := perplexity.NewChatCompletionRequest()
		req.Messages = histories

		s := spinner.New(spinner.CharSets[2], 100*time.Millisecond)
		s.Start()

		res, err := cli.ChatCompletions(req)
		if err != nil {
			return fmt.Errorf("llm: failed to get chat completions: %w", err)
		}

		s.Stop()

		histories = append(histories, res.Choices[0].Message)

		fmt.Printf("[%s]: %s\n", res.Choices[0].Message.Role, res.Choices[0].Message.Content)

		if !interactive {
			break
		}

		input := bufio.NewScanner(os.Stdin)

		fmt.Printf("[%s]: ", perplexity.RoleUser)
		input.Scan()

		histories = append(histories, perplexity.ChatCompletionMessage{
			Role:    perplexity.RoleUser,
			Content: input.Text(),
		})
	}

	return nil
}
