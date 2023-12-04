/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"

	ingestemail "github.com/McFlip/edisco-test/cmd/ingestEmail"
	"github.com/spf13/cobra"
)

// ingestEmailCmd represents the ingestEmail command
var ingestEmailCmd = &cobra.Command{
	Use:   "ingestEmail",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ingestEmail called")
		if err := ingestemail.Ingest(*inDir, *outFile); err != nil {
			log.Fatal(err)
		}
	},
}

var (
	inDir, outFile *string
)

func init() {
	rootCmd.AddCommand(ingestEmailCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	inDir = ingestEmailCmd.PersistentFlags().String("in-dir", "", "input directory containing eml files")
	outFile = ingestEmailCmd.PersistentFlags().String("out", "", "output json file")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ingestEmailCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
