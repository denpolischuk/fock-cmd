package utils

import (
	"strings"

	"github.com/shirou/gopsutil/process"
)

// IsProcessRunning - retrieves the list of running processes and checks whether cmd line includes given string
func IsProcessRunning(substr string) (bool, *process.Process) {
	processes, _ := process.Processes()
	for _, p := range processes {
		pName, _ := p.Name()
		if pName == "node" {
			line, _ := p.Cmdline()
			if strings.Contains(line, substr) {
				return true, p
			}
		}
	}
	return false, nil
}
