/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/Bass-Peerapon/gen-service/cmd/create"
	"github.com/spf13/cobra"
)

const (
	VERSION string = "0.0.2"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gen-service",
	Short: "gen-service is cmd for generate innovasive service directory template",
	Long:  `gen-service is cmd for generate innovasive service directory template 🚀`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, _ []string) {
		v, err := cmd.Flags().GetBool("version")
		if err != nil {
			fmt.Println(err)
			return
		}
		if v {
			fmt.Println(VERSION)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommandPalettes() {
	rootCmd.AddCommand(create.CreateCmd)
}
func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gen-service.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "version")
	addSubcommandPalettes()
}
