package exec

import (
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"strings"
)

type Docker struct {
	Container  string
	WorkingDir string
}

func (d Docker) exec(cmd []string) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(strings.Join(cmd, " ")))

	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		panic(err)
	}

	filter := filters.NewArgs()
	filter.Add("name", d.Container)
	filter.Add("status", "running")

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Limit: 1, Filters: filter})
	if err != nil {
		panic(err)
	}

	execConfig := types.ExecConfig{
		User:         "www-data",
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}

	if len(d.WorkingDir) > 0 {
		execConfig.WorkingDir = d.WorkingDir
	}

	cec, _ := cli.ContainerExecCreate(ctx, containers[0].ID, execConfig)
	res, err := cli.ContainerExecAttach(ctx, cec.ID, types.ExecStartCheck{})

	if err != nil {
		panic(err)
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	_, err = stdcopy.StdCopy(stdout, stderr, res.Reader)
	if err != nil {
		panic(err)
	}
	s := stdout.String()
	fmt.Println(s)
	i := stderr.String()
	fmt.Println(i)
}
