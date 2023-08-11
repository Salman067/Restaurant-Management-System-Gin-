package command

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"os/signal"
	"pi-inventory/common/logger"
	"pi-inventory/dic"
	commonRoute "pi-inventory/route"
	"syscall"
	"time"
)

func init() {
	var serverPort string
	defaultServerPort := viper.GetString("SERVER_PORT")
	serverCmd.PersistentFlags().StringVar(&serverPort, "port", defaultServerPort, "Server port")
	viper.BindPFlag("SERVER_PORT", serverCmd.PersistentFlags().Lookup("port"))
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server",
	Run: func(cmd *cobra.Command, args []string) {
		port := ":" + viper.GetString("SERVER_PORT")
		router := commonRoute.Setup(dic.CommonBuilder)
		logger.LogInfo("Running server on port " + port)
		server := &http.Server{
			Addr:    port,
			Handler: router,
		}

		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logger.LogError("failed to start server", err)
				os.Exit(1)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logger.LogInfo("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.LogError("Server forced to shutdown: ", err)
			os.Exit(1)
		}

		logger.LogInfo("Server exiting")
	},
}
