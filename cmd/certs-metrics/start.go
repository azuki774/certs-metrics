package main

import (
	"certs-metrics/internal/factory"
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

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
		ms := factory.NewMetricsServer(l, us, args)
		err = ms.Start(ctx)
		return err
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
