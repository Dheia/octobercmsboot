/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"github.com/panakour/octobercmsboot"
	"github.com/panakour/octobercmsboot/exec"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update octobercms, plugins and theme",
	Long: `The update command install/update plugins, themes, dependencies and run any new database migrations`,
	Run: func(cmd *cobra.Command, args []string) {
		env := cmd.Flag("env").Value.String()
		october, _ := octobercmsboot.NewOctober("./october.yaml", env)
		var phpRunner = exec.NewDocker("php-fpm", october.Env[env].WorkingDir)
		october.InstallPlugins(phpRunner)
		phpRunner.Run([]string{"composer", "update", "--no-scripts", "--no-interaction", "--prefer-dist", "--lock"})
		october.InstallThemes(phpRunner)
		phpRunner.Run([]string{"php", "artisan", "october:up"})
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	updateCmd.Flags().StringP("env", "e", "dev", "Use prod for production or dev for development")
}
