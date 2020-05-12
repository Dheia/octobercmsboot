package octobercmsboot

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/panakour/octobercmsboot/exec"
	"io/ioutil"
	"os"
	"strings"
)

type Env map[string]*struct {
	WorkingDir string
	App        struct {
		Debug string
		Url   string
	}
	Db struct {
		Connection, Host, Port, Database, Username, Password string
	}
	PhpContainer   string
	MysqlContainer string
}

const envFile = ".env"

func (e Env) Generate(october OctoberCMS, phpRunner exec.Runner) {
	_, err := os.Stat(rootPath() + "/" + envFile)
	if err == nil {
		Info("env already exists")
		return
	}
	e.databaseQuestions(october)
	phpRunner.Run([]string{"php", "artisan", "october:env"})
	e.replace("APP_DEBUG=true", "APP_DEBUG="+e[october.currentEnv].App.Debug)
	e.replace("APP_URL=http://localhost", "APP_URL="+e[october.currentEnv].App.Url)
	e.replace("APP_KEY=CHANGE_ME!!!!!!!!!!!!!!!!!!!!!!!", "APP_KEY="+e.generateKey())
	e.replace("DB_CONNECTION=mysql", "DB_CONNECTION="+e[october.currentEnv].Db.Connection)
	e.replace("DB_HOST=localhost", "DB_HOST="+e[october.currentEnv].Db.Host)
	e.replace("DB_PORT=3306", "DB_PORT="+e[october.currentEnv].Db.Port)
	e.replace("DB_DATABASE=database", "DB_DATABASE="+e[october.currentEnv].Db.Database)
	e.replace("DB_USERNAME=root", "DB_USERNAME="+e[october.currentEnv].Db.Username)
	e.replace("DB_PASSWORD=", "DB_PASSWORD="+e[october.currentEnv].Db.Password)
}

func (e Env) replace(old, new string) {
	read, err := ioutil.ReadFile(envFile)
	if err != nil {
		panic(err)
	}
	newContent := strings.Replace(string(read), old, new, -1)
	err = ioutil.WriteFile(envFile, []byte(newContent), 0)
	if err != nil {
		panic(err)
	}
}

func (e Env) generateKey() string {
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", key)
}

func (e Env) databaseQuestions(october OctoberCMS) {
	if len(e[october.currentEnv].Db.Connection) == 0 {
		prompt := promptui.Select{
			Label: "Choose a DB_CONNECTION",
			Items: []string{"mysql", "postgresql", "sqlite"},
		}
		_, result, _ := prompt.Run()
		e[october.currentEnv].Db.Connection = result
	}

	if len(e[october.currentEnv].Db.Host) == 0 {
		prompt := promptui.Prompt{
			Label:   "DB_HOST",
			Default: "127.0.0.1",
		}
		result, _ := prompt.Run()
		e[october.currentEnv].Db.Host = result
	}

	if len(e[october.currentEnv].Db.Port) == 0 {
		prompt := promptui.Prompt{
			Label:   "DB_PORT",
			Default: "3306",
		}
		result, _ := prompt.Run()
		e[october.currentEnv].Db.Port = result
	}

	if len(e[october.currentEnv].Db.Database) == 0 {
		validate := func(input string) error {
			if len(strings.TrimSpace(input)) == 0 {
				return errors.New("Database is required.")
			}
			return nil
		}

		prompt := promptui.Prompt{
			Label:    "DB_DATABASE",
			Validate: validate,
		}
		result, _ := prompt.Run()
		e[october.currentEnv].Db.Database = result
	}

	if len(e[october.currentEnv].Db.Username) == 0 {
		validate := func(input string) error {
			if len(strings.TrimSpace(input)) == 0 {
				return errors.New("Username is required.")
			}
			return nil
		}

		prompt := promptui.Prompt{
			Label:    "DB_USERNAME",
			Validate: validate,
		}
		result, _ := prompt.Run()
		e[october.currentEnv].Db.Username = result
	}

	if len(e[october.currentEnv].Db.Password) == 0 {
		validate := func(input string) error {
			if len(strings.TrimSpace(input)) == 0 {
				return errors.New("Password is required.")
			}
			return nil
		}

		prompt := promptui.Prompt{
			Label:    "DB_PASSWORD",
			Validate: validate,
		}
		result, _ := prompt.Run()
		e[october.currentEnv].Db.Password = result
	}
}
