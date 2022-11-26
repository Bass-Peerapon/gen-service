/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package create

import (
	"fmt"

	subcreate "github.com/Bass-Peerapon/gen-service/cmd/create/sub_create"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Call any sub create",
	Long:  ``,
	Run: func(cmd *cobra.Command, _ []string) {
		fmt.Println(cmd.Help())
	},
}

func addSubcommandPalettes() {
	CreateCmd.AddCommand(subcreate.ServiceCmd)
	CreateCmd.AddCommand(subcreate.ModelCmd)
}

func init() {
	// rootCmd.AddCommand(createCmd)
	addSubcommandPalettes()
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
