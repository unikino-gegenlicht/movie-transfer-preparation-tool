//go:build windows

package utilities

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"os/exec"
	"strings"
)

// GetExternalDrives gets all external drives connected to a windows machine but excludes C: as it is the
// main hard drive, and we do not need that
func GetExternalDrives() [][]string {
	var drives [][]string
	cmd := exec.Command("wmic", "logicaldisk", "get", "DeviceID,VolumeName", "/format:csv")
	output, err := cmd.Output()
	if err != nil {
		log.Error().Err(err).Stack().Msg("unable to get connected drives")
		return [][]string{}
	}
	log.Debug().Msg("wmic query successful")
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			log.Debug().Msg("skipping empty line of wmic output")
			continue
		}
		if line == "Node,DeviceID,VolumeName" {
			log.Debug().Msg("skipping header of wmic output")
			continue
		}
		fields := strings.Split(line, ",")
		if len(fields) != 3 {
			log.Error().Msg("wmic output in unexpected format")
			return [][]string{}
		}
		if fields[1] == "C:" {
			continue
		}
		drives = append(drives, []string{fields[1], fields[2]})
	}
	return drives
}
