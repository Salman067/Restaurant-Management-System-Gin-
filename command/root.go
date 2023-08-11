package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var rootCmd = &cobra.Command{
	Short: "Application description",
	Long: `Long
                application
                description`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please use command server. Like go run main.go server")
	},
}

func Execute() {
	rootCmd.Use = viper.GetString("APP_COMMAND")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
