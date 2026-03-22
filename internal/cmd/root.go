package cmd

import (
	"fmt"
	"os"

	"github.com/exiquo/zipper/internal/archiver"
	"github.com/spf13/cobra"
)

var src string
var out string

var rootCmd = &cobra.Command{
	Use:   "zipper",
	Short: "Archive a directory into a zip file",
	RunE: func(cmd *cobra.Command, args []string) error {
		return archiver.CreateArchive(src, out)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&src, "src", "", "path to the source directory")
	rootCmd.Flags().StringVar(&out, "out", "", "path to the output zip file")

	_ = rootCmd.MarkFlagRequired("src")
	_ = rootCmd.MarkFlagRequired("out")
}
