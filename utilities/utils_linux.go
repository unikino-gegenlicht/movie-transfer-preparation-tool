package utilities

import (
	"golang.org/x/sys/unix"
	"os/exec"
	"strings"
)

// GetExternalDrives gets all mounted drives on a linux machine
func GetExternalDrives() [][]string {
	var drives [][]string
	cmd := exec.Command("mount")
	out, err := cmd.Output()
	if err != nil {
		return drives
	}

	for _, line := range strings.Split(string(out), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		drive := fields[0]
		mountPath := fields[2]

		if unix.Access(mountPath, unix.W_OK) == nil {
			drives = append(drives, []string{mountPath, drive})
		} else {
			continue
		}
	}
	return drives
}
