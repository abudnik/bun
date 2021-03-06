package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"

	"github.com/mesosphere/bun/v2/tools"
)

var toolCmd = &cobra.Command{
	Use:   "tool",
	Short: "Bun development tool",
	Long:  "Contains subcommands which help to add new file types and checks.",
}

func findFiles(cmd *cobra.Command, args []string) {
	fileTypes, err := tools.FindFiles(bundlePath)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	y, err := yaml.Marshal(&fileTypes)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	escape, err := cmd.Flags().GetBool("escape")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	if escape {
		y = bytes.ReplaceAll(y, []byte("`"), []byte("`+ \"`\" +`"))
	}
	fmt.Println(string(y))
}

func init() {
	rootCmd.AddCommand(toolCmd)
	var findFileCmd = &cobra.Command{
		Use:   "find-files",
		Short: "Finds all file types in a given bundle",
		Long: "Finds all file types in a given bundle, suggests names, and" +
			" renders it in a YAML format to the stdout.",
		Run: findFiles,
	}
	findFileCmd.Flags().BoolP("escape", "e", false, "Escape back ticks for using in the files_yaml.go")
	toolCmd.AddCommand(findFileCmd)
}
