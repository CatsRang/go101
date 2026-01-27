package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"example.com/rest-core-arch/pkg/config"
	"example.com/rest-core-arch/pkg/server"
	"example.com/rest-core-arch/pkg/util/log"
)

var (
	confPath string
	rootCmd  = &cobra.Command{
		Use:   "app",
		Short: "Core REST API Example",
		Run:   rootRun,
	}
)

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&confPath, "config", "c", "config.yaml", "config file path")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	appConf := config.SharedAppConf()

	// 1. Init Config
	if err := appConf.Init(confPath); err != nil {
		// If config file is missing, we might want to fail or warn depending on requirements.
		// For now, we print to stderr as logger isn't ready.
		fmt.Fprintf(os.Stderr, "Warning: Failed to load config: %v. Using defaults.\n", err)
	}

	// 2. Init Logger
	log.InitLogger(appConf.Log.Level, appConf.Log.Format)

	log.L().Info("Config loaded", zap.String("path", confPath), zap.String("app_id", appConf.AppId))
}

func rootRun(cmd *cobra.Command, args []string) {
	appConf := config.SharedAppConf()
	defer log.Sync()

	// 3. Start Server
	srv := server.New(appConf)
	srv.Run() // Blocks until shutdown
}
