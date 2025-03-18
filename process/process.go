package process

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/amjadcp/process/ai"
	"github.com/shirou/gopsutil/v3/process"
)

// ProcessInfo stores information about a single process.
type ProcessInfo struct {
	PID     int32
	Name    string
	Status  string
	CPU     float64
	Memory  float32
	Command string
}

// TrackProcesses monitors process changes and triggers events.
func TrackProcesses(events chan<- string) {
	prevProcesses := make(map[int32]ProcessInfo)
	var mu sync.Mutex

	for {
		currentProcesses, err := GetProcesses()
		if err != nil {
			events <- fmt.Sprintf("Error retrieving processes: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		mu.Lock()
		changes := detectChanges(prevProcesses, currentProcesses)
		mu.Unlock()

		for _, event := range changes {
			events <- event
		}

		// Update the previous process list.
		mu.Lock()
		prevProcesses = make(map[int32]ProcessInfo)
		for _, proc := range currentProcesses {
			prevProcesses[proc.PID] = proc
		}
		mu.Unlock()

		time.Sleep(2 * time.Second) // Poll every 2 seconds.
	}
}

// GetProcesses retrieves a list of running processes.
func GetProcesses() ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var results []ProcessInfo
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, p := range procs {
		wg.Add(1)
		go func(proc *process.Process) {
			defer wg.Done()

			name, _ := proc.Name()
			statusSlice, _ := proc.Status()
			status := strings.Join(statusSlice, ", ")
			cpu, _ := proc.CPUPercent()
			mem, _ := proc.MemoryPercent()
			cmd, _ := proc.Cmdline()

			mu.Lock()
			results = append(results, ProcessInfo{
				PID:     proc.Pid,
				Name:    name,
				Status:  status,
				CPU:     cpu,
				Memory:  mem,
				Command: cmd,
			})
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	return results, nil
}

// detectChanges compares the old and new process lists and identifies events.
func detectChanges(prev map[int32]ProcessInfo, current []ProcessInfo) []string {
	var events []string
	currentMap := make(map[int32]ProcessInfo)

	for _, proc := range current {
		currentMap[proc.PID] = proc

		if _, exists := prev[proc.PID]; !exists {
			// Prepare process details for AI analysis.
			processData := ai.ProcessData{
				PID:     proc.PID,
				Name:    proc.Name,
				Status:  proc.Status,
				CPU:     proc.CPU,
				Memory:  proc.Memory,
				Command: proc.Command,
			}
			analysis, err := ai.AnalyzeProcess(processData)
			description := ""
			if err != nil {
				description = "AI analysis unavailable"
			} else {
				description = analysis.Description
				if analysis.Malicious {
					description += " [WARNING: Potentially Malicious]"
				} else {
					description += " [Safe]"
				}
			}
			events = append(events, fmt.Sprintf("🔵 New Process Started: %s (PID: %d). %s", proc.Name, proc.PID, description))
		} else if prev[proc.PID].Status != proc.Status {
			events = append(events, fmt.Sprintf("⚠️ Process %s (PID: %d) changed status: %s → %s",
				proc.Name, proc.PID, prev[proc.PID].Status, proc.Status))
		}
	}

	for pid, proc := range prev {
		if _, exists := currentMap[pid]; !exists {
			events = append(events, fmt.Sprintf("❌ Process Stopped: %s (PID: %d)", proc.Name, pid))
		}
	}

	return events
}
