package config

import (
	"bufio"
	"fmt"
	"github.com/rohanthewiz/serr"
	"os"
	"strings"
)

func ReadEnvFile() (err error) {
	envf, err := os.Open("env.vars")
	if err != nil {
		return serr.Wrap(err, "Error opening env file")
	}

	scnr := bufio.NewScanner(envf)

	for scnr.Scan() {
		line := scnr.Text()
		if strings.TrimSpace(line) != "" {
			if strings.HasPrefix(line, "#") {
				continue
			}
			arr := strings.SplitN(line, "=", 2)
			if len(arr) == 2 {
				key := strings.TrimSpace(arr[0])
				val := strings.TrimSpace(arr[1])
				er := os.Setenv(key, val)
				if er != nil {
					fmt.Println("Failed to set env var key:", key, " value:", val)
				} else {
					fmt.Println("Successfully set env var key:", key, " value:", val)
				}
			} else {
				fmt.Println("Not setting environment vars from this line:", line)
			}
		}
	}

	return
}
