package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log.level", "l", "info", "log level (debug, info, warn, error)")
}

var logLevel = "info"

var rootCmd = &cobra.Command{
	Use:          "dping_exporter",
	Version:      version,
	Short:        "Distributed Ping Exporter",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		initLogger()
		return nil
	},
}

// initLogger initializes logger
func initLogger() {
	level := &slog.LevelVar{}

	switch logLevel {
	default:
		level.Set(slog.LevelInfo)
	case "debug":
		level.Set(slog.LevelDebug)
	case "info":
		level.Set(slog.LevelInfo)
	case "warn":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
