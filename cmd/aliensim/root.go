package aliensim

import (
	"alien-invasion-simulator/pkg/aliemsim"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strconv"
)

var rootCmd = &cobra.Command{
	Use:   "aliensim",
	Short: "aliensim - a simple CLI to simulate alien invasions",
	Long:  `Provide a sample .txt file (arg[0]) with cities and a number of aliens (arg[1]). Aliensim will simulate the invasion.`,
	Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		filePath := args[0]
		numAliens, err := strconv.Atoi(args[1])

		if err != nil {
			log.Fatalf("Invalid number of aliens provided: %v", args[1])
			os.Exit(1)

		}
		verbose, _ := cmd.Flags().GetBool("verbose")
		err = aliemsim.StartSimulation(filePath, numAliens, aliemsim.OsFS, 10000, verbose)
		if err != nil {
			log.Fatalf("Simulation Failed %v", err)
			os.Exit(1)

		}
	},
}

func Init() {
	rootCmd.PersistentFlags().Bool("verbose", false, "A print map stats on every iteration")
}
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)

		os.Exit(1)
	}
}
