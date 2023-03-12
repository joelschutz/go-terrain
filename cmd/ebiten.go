/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"plugin"

	"github.com/spf13/cobra"
)

// ebitenCmd represents the ebiten command
var ebitenCmd = &cobra.Command{
	Use:   "ebiten",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ebiten called")
		// load module
		// 1. open the so file to load the symbols
		plug, err := plugin.Open("./engine/ebiten/plugin.so")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 2. look up a symbol (an exported function or variable)
		// in this case, variable Greeter
		app, err := plug.Lookup("EbitenApp")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// 3. Assert that loaded symbol is of a desired type
		// in this case interface type Greeter (defined above)
		entryPoint, ok := app.(func())
		if !ok {
			fmt.Println("unexpected type from module symbol")
			os.Exit(1)
		}

		// 4. use the module
		entryPoint()
	},
}

func init() {
	rootCmd.AddCommand(ebitenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ebitenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ebitenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
