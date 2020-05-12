package cmd

import (
	"github.com/panakour/octobercmsboot"
	"github.com/panakour/octobercmsboot/exec"
	"github.com/spf13/cobra"
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install october",
	Long:  `This command is to quickly create fresh installation of octobercms with plugins and theme.`,
	Run: func(cmd *cobra.Command, args []string) {
		env := cmd.Flag("env").Value.String()
		installOctober(env)
	},
}

func init() {
	//*** new implementation => command flag that let users choose the runner ***
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:

	installCmd.Flags().StringP("env", "e", "dev", "Use prod for production or dev for development")
}

func installOctober(env string) {
	//var runner octobercmsboot.Docker
	october, _ := octobercmsboot.NewOctober("./october.yaml", env)
	october.Download()
	var phpRunner = exec.NewDocker("php-fpm", october.Env[env].WorkingDir)
	var mysqlRunner = exec.NewDocker("mysql", "")
	phpRunner.Run([]string{"composer", "install", "--no-scripts", "--no-interaction", "--prefer-dist"})
	october.Env.Generate(october, phpRunner)
	createSchemaCommand := exec.CreateSchemaCommand(october.Env[env].Db.Database, october.Env[env].Db.Username, october.Env[env].Db.Password)
	mysqlRunner.Run(createSchemaCommand)
	phpRunner.Run([]string{"php", "artisan", "october:up"})
	october.InstallPlugins(phpRunner)
	phpRunner.Run([]string{"composer", "update", "--no-scripts", "--no-interaction", "--prefer-dist", "--lock"})
	october.InstallThemes(phpRunner)
	phpRunner.Run([]string{"php", "artisan", "october:fresh"})
	phpRunner.Run([]string{"php", "artisan", "cache:clear"})
	phpRunner.Run([]string{"php", "artisan", "october:util", "set", "build"})
	phpRunner.Run([]string{"php", "artisan", "october:up"})

}
