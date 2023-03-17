package utilities

import (
	"golang.org/x/sys/unix"
	"os/exec"
	"strings"
	"syscall"
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

func GetAvailableSpace(mountPath string) uint64 {
	var freeBytes uint64

	var stat syscall.Statfs_t
	err := syscall.Statfs(mountPath, &stat)
	if err != nil {
		fmt.Println(err)
		return 0
	}

	// Available blocks * size per block = available space in bytes
	freeBytes = stat.Bavail * uint64(stat.Bsize)

	return freeBytes
}
