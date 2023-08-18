package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	ExampleGeneration = "example"
)

var exampleContent = `name: Mokk Example Server
routes:
  - path: "/"
    method: "GET"
    statusCode: 200
    response: '{"status":"Success","message":"Hello world!"}'
`

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a mokk.yml file in the current directory.",
	Long: `
 _______         __     __
|   |   |.-----.|  |--.|  |--.
|       ||  _  ||    < |    <
|__|_|__||_____||__|__||__|__|

Generate a mokk.yml file in the current directory.

Supported generators:
  - example
`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		if len(args) != 1 {
			panic("Invalid number of args, 1 expected")
		}

		gen := args[0]

		switch gen {
		case ExampleGeneration:
			err = os.WriteFile("example.mokk.yml", []byte(exampleContent), 0666)
			if err != nil {
				err = fmt.Errorf("failed to generate file: %w", err)
			}
		default:
			err = fmt.Errorf("failed to generate file: unsupported generator provided")
		}

		if err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
