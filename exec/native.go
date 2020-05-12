package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Native struct{}

func (n Native) exec(command []string) {
	cmd := strings.Join(command, " ")
	fmt.Printf(cmd)
	cmdRunner := exec.Command("bash", "-c", cmd)
	n.run(cmdRunner)
}

func (n Native) run(cmd *exec.Cmd) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	if len(out.String()) > 0 {
		fmt.Println("Result: " + out.String())
	}
}
