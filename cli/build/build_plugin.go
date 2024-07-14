package build

import (
	"os"
	"os/exec"
)

type pluginBuilder struct {
	Dir       string
	DebugMode bool
}

func (builder *pluginBuilder) Build(source_file string, destination_file string) error {
	//err := builder.initModule("my-plugin")

	// if err != nil {
	// 	return err
	// }

	err := builder.downloadDependencies()

	if err != nil {
		return err
	}

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", destination_file, source_file)
	cmd.Dir = builder.Dir

	if builder.DebugMode {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	return cmd.Run()
}

func (builder *pluginBuilder) downloadDependencies() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = builder.Dir

	if builder.DebugMode {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	return cmd.Run()
}

func (builder *pluginBuilder) initModule(package_name string) error {
	cmd := exec.Command("go", "mod", "init", package_name)
	cmd.Dir = builder.Dir

	if builder.DebugMode {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
	}

	return cmd.Run()
}
