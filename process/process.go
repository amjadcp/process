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
	PID         int32   `json:"pid"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	CPU         float64 `json:"cpu"`
	Memory      float32 `json:"memory"`
	Command     string  `json:"command"`
	Description string  `json:"description"`
}

// TrackProcesses monitors process changes and triggers events.
// It uses a polling interval and sends events to the provided channel.
func TrackProcesses(events chan<- ProcessInfo, pollInterval time.Duration) {
	prevProcesses := make(map[int32]ProcessInfo)

	for {
		currentProcesses, err := GetProcesses()
		if err != nil {
			events <- ProcessInfo{
				PID:         -1,
				Name:        "Error",
				Status:      "Error",
				CPU:         0.0,
				Memory:      0.0,
				Command:     "",
				Description: fmt.Sprintf("Error retrieving processes: %v", err),
			}
			time.Sleep(pollInterval)
			continue
		}

		// Detect changes between previous and current process states.
		newProcs, changeEvents := detectChanges(prevProcesses, currentProcesses)

		// Send change events (status changes, stopped processes).
		for _, event := range changeEvents {
			events <- event
		}

		// For each new process, analyze it asynchronously.
		for _, proc := range newProcs {
			go analyzeNewProcess(proc, events)
		}

		// Update previous processes.
		prevProcesses = make(map[int32]ProcessInfo)
		for _, proc := range currentProcesses {
			prevProcesses[proc.PID] = proc
		}

		time.Sleep(pollInterval)
	}
}

// GetProcesses retrieves a list of running processes using a worker pool to limit concurrency.
func GetProcesses() (map[int32]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(map[int32]ProcessInfo)
	concurrencyLimit := 10
	sem := make(chan struct{}, concurrencyLimit)

	for _, p := range procs {
		sem <- struct{}{}
		wg.Add(1)
		go func(proc *process.Process) {
			defer wg.Done()
			defer func() { <-sem }()

			name, err := proc.Name()
			if err != nil {
				name = "Unknown"
			}
			statusSlice, err := proc.Status()
			status := "Unknown"
			if err == nil {
				status = strings.Join(statusSlice, ", ")
			}
			cpu, err := proc.CPUPercent()
			if err != nil {
				cpu = 0.0
			}
			mem, err := proc.MemoryPercent()
			if err != nil {
				mem = 0.0
			}
			cmd, err := proc.Cmdline()
			if err != nil {
				cmd = ""
			}

			mu.Lock()
			results[proc.Pid] = ProcessInfo{
				PID:     proc.Pid,
				Name:    name,
				Status:  status,
				CPU:     cpu,
				Memory:  mem,
				Command: cmd,
			}
			mu.Unlock()
		}(p)
	}

	wg.Wait()
	return results, nil
}

// detectChanges compares old and new process lists and identifies new processes and other events.
func detectChanges(prev map[int32]ProcessInfo, current map[int32]ProcessInfo) (newProcs map[int32]ProcessInfo, events []ProcessInfo) {
	newProcs = make(map[int32]ProcessInfo)

	for pid, proc := range current {
		if _, isExisted := prev[pid]; !isExisted {
			newProcs[pid] = proc
		} else if prev[pid].Status != proc.Status {
			events = append(events, ProcessInfo{
				PID:     proc.PID,
				Name:    proc.Name,
				Status:  proc.Status,
				CPU:     proc.CPU,
				Memory:  proc.Memory,
				Command: proc.Command,
				Description: fmt.Sprintf("⚠️ Process %s (PID: %d) changed status: %s → %s",
					proc.Name, proc.PID, prev[pid].Status, proc.Status),
			})
		}
	}

	for pid, proc := range prev {
		if _, isExisted := current[pid]; !isExisted {
			events = append(events, ProcessInfo{
				PID:         proc.PID,
				Name:        proc.Name,
				Status:      proc.Status,
				CPU:         proc.CPU,
				Memory:      proc.Memory,
				Command:     proc.Command,
				Description: fmt.Sprintf("❌ Process Stopped: %s (PID: %d)", proc.Name, pid),
			})
		}
	}

	return newProcs, events
}

// analyzeNewProcess performs AI analysis for a new process and sends an event message.
func analyzeNewProcess(proc ProcessInfo, events chan<- ProcessInfo) {
	// Prepare process data for AI analysis.
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
		fmt.Printf("Getting Error from the AI Service: %s \n", err.Error())
		description = "AI analysis unavailable"
	} else {
		description = analysis.Description
		if analysis.Malicious {
			description += " [WARNING: Potentially Malicious]"
		} else {
			description += " [Safe]"
		}
	}
	events <- ProcessInfo{
		PID:         proc.PID,
		Name:        proc.Name,
		Status:      proc.Status,
		CPU:         proc.CPU,
		Memory:      proc.Memory,
		Command:     proc.Command,
		Description: description,
	}
}
