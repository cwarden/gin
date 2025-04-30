package gin

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Builder interface {
	Build() error
	Binary() string
	Errors() string
}

type builder struct {
	dir         string
	binary      string
	errors      string
	useGodep    bool
	wd          string
	preBuildCmd string
	buildArgs   []string
}

// NewBuilder creates a new Builder. preBuildCmd is the command to run before building (empty means skip pre-build step).
func NewBuilder(dir string, bin string, useGodep bool, wd string, buildArgs []string, preBuildCmd string) Builder {
	if len(bin) == 0 {
		bin = "bin"
	}

	// does not work on Windows without the ".exe" extension
	if runtime.GOOS == "windows" {
		if !strings.HasSuffix(bin, ".exe") { // check if it already has the .exe extension
			bin += ".exe"
		}
	}

	preBuildCmd = strings.TrimSpace(preBuildCmd)
	return &builder{dir: dir, binary: bin, useGodep: useGodep, wd: wd, preBuildCmd: preBuildCmd, buildArgs: buildArgs}
}

func (b *builder) Binary() string {
	return b.binary
}

func (b *builder) Errors() string {
	return b.errors
}

func (b *builder) Build() error {
	// Pre-build step: run the specified command and propagate errors
	if b.preBuildCmd != "" {
		var preCmd *exec.Cmd
		if runtime.GOOS == "windows" {
			preCmd = exec.Command("cmd.exe", "/C", b.preBuildCmd)
		} else {
			preCmd = exec.Command("sh", "-c", b.preBuildCmd)
		}
		preCmd.Dir = b.dir
		preOut, preErr := preCmd.CombinedOutput()
		if len(preOut) > 0 {
			fmt.Printf("%s", preOut)
		}
		if preErr != nil {
			b.errors = string(preOut)
			return fmt.Errorf(b.errors)
		}
	}

	// Build step
	args := append([]string{"go", "build", "-o", filepath.Join(b.wd, b.binary)}, b.buildArgs...)
	if b.useGodep {
		args = append([]string{"godep"}, args...)
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = b.dir
	output, err := cmd.CombinedOutput()
	if cmd.ProcessState.Success() {
		b.errors = ""
	} else {
		b.errors = string(output)
	}
	if len(b.errors) > 0 {
		return fmt.Errorf(b.errors)
	}
	return err
}
