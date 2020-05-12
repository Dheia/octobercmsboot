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
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize root directory for octobercms",
	Long:  `Init command will prepared the root directory of the octobercms with the default october.yaml and .gitignore files.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires 1 argument for the root directory. If current dir specify '.' (dot) ")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		projectRootDir := args[0]
		initializeOctoberProject(projectRootDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func initializeOctoberProject(projectName string) {
	if projectName != "." {
		err := os.Mkdir(projectName, 0755)
		if err != nil {
			fmt.Print(err)
		}
	}
	var projectPath strings.Builder
	projectPath.WriteString(projectName)
	projectPath.WriteString("/")
	viper.SetConfigName("octobercms")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath("$HOME/.octobercms")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	err = viper.WriteConfigAs(projectPath.String() + "october.yaml")
	if err != nil {
		panic(fmt.Errorf("Fatal error: %s \n", err))
	}
	fmt.Println("october.yaml has been created.")
	_, err = copyGitignore(projectPath.String() + ".gitignore")
	if err != nil {
		panic(fmt.Errorf("Fatal error: %s \n", err))
	}
	fmt.Println(".gitignore has been created.")
}


func copyGitignore(dst string) (int64, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}
	fileSrc := usr.HomeDir + "/.octobercms/gitignore.txt"
	sourceFileStat, err := os.Stat(fileSrc)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", fileSrc)
	}

	source, err := os.Open(fileSrc)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
