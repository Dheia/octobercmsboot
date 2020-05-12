package octobercmsboot

import (
	"github.com/panakour/octobercmsboot/exec"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type OctoberCMS struct {
	currentEnv string
	Env
	Plugins
	Themes
}

type Plugins struct {
	Marketplace  []string
	Repositories []Repository
}

type ThemeMarketplace struct {
	Name, Path string
}

type Themes struct {
	Use          string
	Marketplace  []ThemeMarketplace
	Repositories []Repository
}

func NewOctober(octoberYamlConfigFile, env string) (OctoberCMS, error) {
	var octobercms OctoberCMS
	octobercms.currentEnv = env
	yamlFile, err := ioutil.ReadFile(octoberYamlConfigFile)
	if err != nil {
		log.Fatalf("Reading october.yaml err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &octobercms)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return octobercms, err
}

func (o OctoberCMS) Download() {
	if o.isInstalled() {
		Info("October is already downloaded. Remove modules directory to download it again.")
		return
	}
	Info("Download OctoberCMS")
	err := downloadFile("october.zip", "https://github.com/octobercms/october/archive/master.zip")
	if err != nil {
		log.Fatalf("error download october: %v", err)
	}
	err = unzip("october.zip", ".")
	if err != nil {
		log.Fatalf("error on unzip october: %v", err)
	}
	os.Remove("october.zip")
	os.Remove("README.md")
	os.Remove("LICENSE")
	os.Remove("CHANGELOG.md")
	os.Remove("CODE_OF_CONDUCT.md")
	os.Remove("SECURITY.md")
	os.RemoveAll(".github")
}

func (o OctoberCMS) InstallPlugins(phpRunner exec.Runner) {
	var wg sync.WaitGroup
	for _, plugin := range o.Plugins.Marketplace {
		wg.Add(1)
		go func(plugin string) {
			defer wg.Done()
			phpRunner.Run([]string{"php", "artisan", "plugin:Install", plugin})

		}(plugin)
	}
	for _, plugin := range o.Plugins.Repositories {
		wg.Add(1)
		go func(plugin Repository) {
			defer wg.Done()
			plugin.Path = rootPath() + "/plugins/" + plugin.Path
			git := newGit(plugin)
			git.install()
		}(plugin)
	}
	wg.Wait()
}

func (o OctoberCMS) InstallThemes(phpRunner exec.Runner) {
	var wg sync.WaitGroup
	for _, theme := range o.Themes.Marketplace {
		wg.Add(1)
		go func(theme ThemeMarketplace) {
			defer wg.Done()
			phpRunner.Run([]string{"php", "artisan", "theme:Install", theme.Name, theme.Path})
		}(theme)
	}
	for _, theme := range o.Themes.Repositories {
		wg.Add(1)
		go func(theme Repository) {
			defer wg.Done()
			theme.Path = rootPath() + "/themes/" + theme.Path
			git := newGit(theme)
			git.install()
		}(theme)
	}
	wg.Wait()
}

func (o OctoberCMS) isInstalled() bool {
	_, err := os.Stat(rootPath() + "/modules")
	if err == nil {
		return true
	}
	return false
}
