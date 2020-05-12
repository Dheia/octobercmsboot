package exec

type Runner interface {
	Run(cmd []string)
}

func NewDocker(container string, workDir string) *Docker {
	return &Docker{
		Container:  container,
		WorkingDir: workDir,
	}
}

func (d Docker) Run(cmd []string) {
	d.exec(cmd)
}

func (n Native) Run(cmd []string) {
	n.exec(cmd)
}
