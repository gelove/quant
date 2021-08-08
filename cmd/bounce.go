package cmd

import (
	"log"
	"quant/internal/app/logic"
	"quant/pkg/utils/json"

	"github.com/spf13/cobra"
)

var number int

//
var bounceCmd = &cobra.Command{
	Use:   "bounce [start time]",
	Short: "获取某个时间大跌之后反弹最多的币种",
	Long: `echo is for echoing anything back.
Echo works a lot like print, except it has a child command.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Printf("bounce called, time: %s", args[0])
		res := logic.GetTopBounce(number, args[0])
		log.Printf("%s之后反弹最多的币种 => %s", args[0], json.MustToString(res))
		return nil
	},
}

func init() {
	bounceCmd.Flags().IntVarP(&number, "number", "n", 20, "The quantity of top bounce")
	rootCmd.AddCommand(bounceCmd)
}
