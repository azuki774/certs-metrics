package main

import (
	"certs-metrics/internal/factory"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

var startOpt startOption

type startOption struct {
	Port string
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l, err := factory.NewLogger()
		if err != nil {
			return err
		}
		if len(args) == 0 {
			l.Error("required a certication file")
			return fmt.Errorf("required a certication file")
		}

		ctx := context.Background()
		us := factory.NewUsecase(l)
		ms := factory.NewMetricsServer(l, us, args, startOpt.Port)
		err = ms.Start(ctx)
		return err
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVar(&startOpt.Port, "port", "8334", "listen port")
}
