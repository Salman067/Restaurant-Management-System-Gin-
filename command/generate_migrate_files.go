package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(migrationGenerateCmd)
}

var migrationGenerateCmd = &cobra.Command{
	Use:   "generate-migrate-file",
	Short: "Generate Migration File",
	Run: func(cmd *cobra.Command, args []string) {
		execute := exec.Command("goose", "-dir", "db/migrations", "create", args[0], "sql")
		execute.Stdout = os.Stdout
		execute.Stderr = os.Stderr

		// Run still runs the command and waits for completion
		// but the output is instantly piped to Stdout
		if err := execute.Run(); err != nil {
			fmt.Println("could not run command: ", err)
		}
	},
}
