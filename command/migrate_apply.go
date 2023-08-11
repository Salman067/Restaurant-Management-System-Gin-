package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
)

func init() {
	rootCmd.AddCommand(migrationApplyCmd)
}

var migrationApplyCmd = &cobra.Command{
	Use:   "apply-migrate-files",
	Short: "Apply Migration Files",
	Run: func(cmd *cobra.Command, args []string) {
		db := viper.GetString("DB")
		dbUrl := viper.GetString("DB_URL")
		execute := exec.Command("goose", "-dir", "db/migrations", db, dbUrl, args[0])
		execute.Stdout = os.Stdout
		execute.Stderr = os.Stderr

		// Run still runs the command and waits for completion
		// but the output is instantly piped to Stdout
		if err := execute.Run(); err != nil {
			fmt.Println("could not run command: ", err)
		}
	},
}
