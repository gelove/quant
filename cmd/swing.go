/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"quant/internal/app/logic"
	"quant/pkg/utils/json"

	"github.com/spf13/cobra"
)

var count int

// swingCmd represents the swing command
var swingCmd = &cobra.Command{
	Use:   "swing",
	Short: "获取振幅最优的标的",
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("swing called, count: %d", count)
		list := logic.GetTopAmplitude(count)
		log.Printf("振幅最优的标的 => %s", json.MustToString(list))
	},
}

func init() {
	swingCmd.Flags().IntVarP(&count, "number", "n", 20, "The quantity of best amplitude")
	rootCmd.AddCommand(swingCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// swingCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// swingCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
