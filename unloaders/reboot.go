package unloaders

import (
	"github.com/rohanthewiz/serr"
	"log"
	"os"
	"os/exec"
)

func Reboot() (err error) {
	rebootPath, err := exec.LookPath("reboot")
	if err != nil {
		return serr.Wrap(err, "Could not find `reboot` in path")
	}

	cmdReboot := &exec.Cmd{
		Path:   rebootPath,
		Args:   []string{rebootPath},
		Env:    nil,
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}

	err = cmdReboot.Start()
	if err != nil {
		return serr.Wrap(err, "error initiating a reboot")
	}

	log.Println("Reboot initiated...")
	return
}
