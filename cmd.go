package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"

	"github.com/spf13/cobra"
)

func cmdExecInit() *cobra.Command {
	var directory string
	var command = &cobra.Command{
		Use:   "init",
		Short: "Initializes the ADR configurations",
		Long:  "Initializes the ADR configuration with an optional ADR base directory\n This is a a prerequisite to running any other adr sub-command",
		Run: func(cmd *cobra.Command, args []string) {
			color.Green("Initializing ADR base at " + directory)
			helper := NewAdrHelper(directory)

			if err := helper.InitBaseDir(directory); err != nil {
				msg := fmt.Sprintf("ops failed to init dir: %s\n", err)
				color.Red(msg)
				return
			}

			if err := helper.InitConfig(); err != nil {
				msg := fmt.Sprintf("ops failed to init config: %s\n", err)
				color.Red(msg)
				return
			}

			if err := helper.InitTemplate(); err != nil {
				msg := fmt.Sprintf("ops failed to init template: %s\n", err)
				color.Red(msg)
				return
			}

			return

		},
	}
	command.Flags().StringVarP(
		&directory,
		"directory", "d", "",
		"adr directory path")

	return command
}

func cmdExecNew() *cobra.Command {
	var directory, title string
	var command = &cobra.Command{
		Use:   "new",
		Short: "Create a new ADR",
		Long:  "",
		Run: func(cmd *cobra.Command, args []string) {
			color.Green("Initializing ADR base at " + directory)
			helper := NewAdrHelper(directory)

			currentConfig := helper.GetConfig()
			currentConfig.CurrentAdr++
			if err := helper.UpdateConfig(currentConfig); err != nil {
				msg := fmt.Sprintf("ops failed to init template: %s\n", err)
				color.Red(msg)
			}
			helper.NewAdr(currentConfig, title)

		},
	}
	command.Flags().StringVarP(&directory, "directory", "d", adrDefaultBaseDirName, "adr directory path")
	command.Flags().StringVarP(&title, "title", "t", "", "adr title")
	return command
}

func runCmd() {
	var rootCmd = &cobra.Command{Use: "adr"}
	rootCmd.AddCommand(cmdExecInit())
	rootCmd.AddCommand(cmdExecNew())

	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("ops.. failed to execute command: %s\n", err)
		os.Exit(1)
	}
	return
}
