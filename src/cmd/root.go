package cmd

import (
	"fmt"
	"os"

	"github.com/nodauf/ToDns/src/server"
	"github.com/spf13/cobra"
)

var serverOptions server.Options

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ToDns",
	Short: "Transfer over DNS",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		serverOptions.Serve()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.Flags().StringVarP(&serverOptions.File, "file", "f", "", "File to send")
	rootCmd.Flags().IntVarP(&serverOptions.Wait, "wait", "w", 0, "Wait before sending the response (in ms)")
	rootCmd.Flags().IntVarP(&serverOptions.Size, "size", "s", 250, "Size of returned data per request")
	rootCmd.Flags().BoolVarP(&serverOptions.Verbose, "verbose", "v", false, "verbose")
	rootCmd.MarkFlagRequired("file")
}
