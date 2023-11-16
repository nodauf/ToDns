package cmd

import (
	"fmt"
	"os"

	"github.com/nodauf/ToDns/src/helpers/enumflag"
	"github.com/nodauf/ToDns/src/server"
	"github.com/spf13/cobra"
)

var serverOptions server.Options
var numberOfIPsReturned int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ToDns",
	Short: "Transfer over DNS",
	Long: `To download a file, the server will split it in multiple chunks (default max size is 250) and send the corresponding chunk according to the TXT or A query (<numericalValue>.<domainName>.<tld>).

To be more stealthy the "wait" argument could be adjust to wait a specific time before answering the request. This will delay all the requests as the clients are not multi threading based.
The parameter "size" could also be used to send less data for each response.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(serverOptions.QueryType)
		// Few settings before
		if serverOptions.QueryType == "A" {
			serverOptions.Size = 4 * numberOfIPsReturned
		}
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
	rootCmd.Flags().IntVarP(&serverOptions.Size, "size", "s", 250, "Size of returned data per request. Only for TXT query")
	rootCmd.Flags().IntVarP(&numberOfIPsReturned, "number-ip", "n", 10, "Number of IP returned. Only for A query. Max around 2700")
	rootCmd.Flags().BoolVarP(&serverOptions.Verbose, "verbose", "v", false, "verbose")
	rootCmd.Flags().StringVarP(&serverOptions.ListenAddress, "listen", "l", "0.0.0.0", "Address to listen on")
	rootCmd.Flags().VarP(
		enumflag.New(&serverOptions.QueryType, "A", "TXT", "A"),
		"query", "q",
		"Type of DNS query that will be answered. can be 'TXT' or 'A'")
	rootCmd.MarkFlagRequired("file")
}
