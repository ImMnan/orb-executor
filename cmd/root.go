/*
Copyright © 2024 Manan Patel immnan333@gmail.com
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/immnan/orca/cmd/run"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "orb",
	Short: "A Distributed Automation testing application",
	Long:  ``,
	// Uncomment the following line if your bare application
	Run: func(cmd *cobra.Command, args []string) {
		version, _ := cmd.Flags().GetBool("version")
		license, _ := cmd.Flags().GetBool("license")
		if version {
			fmt.println("Version: Not-Defined")
			fmt.println("Author:  https://github.com/ImMnan/")
			fmt.Println(`
Orb-executor Copyright (C) 2024 Manan Patel
This program comes with ABSOLUTELY NO WARRANTY; 
This is free software, and you are welcome to redistribute it
under certain conditions; type "bmgo --license" for details.
	`)
		} else if license {
			fmt.Println(`Orca - A distributed Automation testing executor
			Copyright (C) 2024  Manan Patel

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details. 

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.`)
		} else {
			cmd.Help()
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

func addSubCommand() {
	rootCmd.AddCommand(run.RunCmd)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.orca.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().BoolP("version", "v", false, "Version of the installed orb!")
	rootCmd.Flags().BoolP("license", "l", false, "Show license details")
	addSubCommand()
}
