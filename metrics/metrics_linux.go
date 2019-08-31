// +build linux

package metrics

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

func AvailableMetrics() []Metric {
	return []Metric{
		{"cpu", "Percentage of cpu used.", "%f", "%", float64(0)},
		{"mem", "Percentage of RAM used.", "%f", "%", float64(0)},
		{"swap", "Amount of memory that has been swapped out to disk (bytes).", "%d", "bytes", uint64(0)},

		{"user", "Percentage of CPU utilization that occurred while executing at the user level.", "%f", "%", float64(0)},
		{"system", "Percentage of CPU utilization that occurred while executing at the system level.", "%f", "%", float64(0)},
		{"idle", "Percentage of time that CPUs were idle and the system did not have an outstanding disk I/O request.", "%f", "%", float64(0)},
		{"nice", "Percentage of CPU utilization that occurred while executing at the user level with nice priority.", "%f", "%", float64(0)},

		{"iowait", "Percentage of time that CPUs were idle during which the system had an outstanding disk I/O request.", "%f", "%", float64(0)},
		{"irq", "Percentage of time spent by CPUs to service hardware interrupts.", "%f", "%", float64(0)},
		{"softirq", "Percentage of time spent by CPUs to service software interrupts.", "%f", "%", float64(0)},
		{"steal", "Percentage of time spent in involuntary wait by the virtual CPUs while the hypervisor was servicing another virtual processor.", "%f", "%", float64(0)},
		{"guest", "Percentage of time spent by CPUs to run a virtual processor.", "%f", "%", float64(0)},
		{"guest_nice", "Percentage of time spent by CPUs to run a virtual processor with nice priority.", "%f", "%", float64(0)},

		{"load1", "Load avarage for 1 minute.", "%f", "", float64(0)},
		{"load5", "Load avarage for 5 minutes.", "%f", "", float64(0)},
		{"load15", "Load avarage for 15 minutes.", "%f", "", float64(0)},

		{"numcpu", "Number of logical CPUs.", "%d", "", int(0)},
	}
}

func AvailableProcMetrics() []Metric {
	return []Metric{
		{"proc_cpu", "Percentage of the CPU time the process uses.", "%f", "%", float64(0)},
		{"proc_mem", "Percentage of the total RAM the process uses.", "%f", "%", float32(0)},
		{"proc_rss", "Non-swapped physical memory the process uses (bytes).", "%d", "bytes", uint64(0)},
		{"proc_vms", "Amount of virtual memory the process uses (bytes).", "%d", "bytes", uint64(0)},
		{"proc_swap", "Amount of memory that has been swapped out to disk the process uses (bytes).", "%d", "bytes", uint64(0)},
		{"proc_open_files", "Amount of files and file discripters opend by the process.", "%d", "", int(0)},

		{"proc_count", "Number of the processes.", "%d", "", int(0)},
	}
}

func (m *Metrics) Collect() error {
	wg := &sync.WaitGroup{}

	// 3 = goroutine count
	errChan := make(chan error, 3)

	if len(m.procPIDs) > 0 {
		wg.Add(1)
		go m.collectProc(wg)
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		cpuPercent, err := cpu.Percent(m.collectInterval, false)
		if err != nil {
			errChan <- err
			return
		}
		m.Store("cpu", cpuPercent[0])
		wg.Done()
	}(wg)

	vm, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	m.Store("mem", vm.UsedPercent)
	sw, err := mem.SwapMemory()
	if err != nil {
		return err
	}
	m.Store("swap", sw.Used)
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		before, err := cpu.Times(false)
		if err != nil {
			errChan <- err
			return
		}
		beforeTotal := before[0].Total()

		if m.collectInterval == 0 {
			m.Store("user", before[0].User/beforeTotal*100)
			m.Store("system", before[0].System/beforeTotal*100)
			m.Store("idle", before[0].Idle/beforeTotal*100)
			m.Store("nice", before[0].Nice/beforeTotal*100)
			m.Store("iowait", before[0].Iowait/beforeTotal*100)
			m.Store("irq", before[0].Irq/beforeTotal*100)
			m.Store("softirq", before[0].Softirq/beforeTotal*100)
			m.Store("steal", before[0].Steal/beforeTotal*100)
			m.Store("guest", before[0].Guest/beforeTotal*100)
			m.Store("guest_nice", before[0].GuestNice/beforeTotal*100)
			return
		}

		time.Sleep(m.collectInterval)

		after, err := cpu.Times(false)
		if err != nil {
			errChan <- err
			return
		}
		afterTotal := after[0].Total()

		total := afterTotal - beforeTotal
		if total == 0 {
			m.Store("user", after[0].User/afterTotal*100)
			m.Store("system", after[0].System/afterTotal*100)
			m.Store("idle", after[0].Idle/afterTotal*100)
			m.Store("nice", after[0].Nice/afterTotal*100)
			m.Store("iowait", after[0].Iowait/afterTotal*100)
			m.Store("irq", after[0].Irq/afterTotal*100)
			m.Store("softirq", after[0].Softirq/afterTotal*100)
			m.Store("steal", after[0].Steal/afterTotal*100)
			m.Store("guest", after[0].Guest/afterTotal*100)
			m.Store("guest_nice", after[0].GuestNice/afterTotal*100)
			return
		}

		m.Store("user", (after[0].User-before[0].User)/total*100)
		m.Store("system", (after[0].System-before[0].System)/total*100)
		m.Store("idle", (after[0].Idle-before[0].Idle)/total*100)
		m.Store("nice", (after[0].Nice-before[0].Nice)/total*100)
		m.Store("iowait", (after[0].Iowait-before[0].Iowait)/total*100)
		m.Store("irq", (after[0].Irq-before[0].Irq)/total*100)
		m.Store("softirq", (after[0].Softirq-before[0].Softirq)/total*100)
		m.Store("steal", (after[0].Steal-before[0].Steal)/total*100)
		m.Store("guest", (after[0].Guest-before[0].Guest)/total*100)
		m.Store("guest_nice", (after[0].GuestNice-before[0].GuestNice)/total*100)
	}(wg)

	l, err := load.Avg()
	if err != nil {
		return err
	}
	m.Store("load1", l.Load1)
	m.Store("load5", l.Load5)
	m.Store("load15", l.Load15)

	m.Store("numcpu", runtime.NumCPU())

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
	}

	return nil
}

func (m *Metrics) collectProc(wg *sync.WaitGroup) {
	defer wg.Done()

	procWg := &sync.WaitGroup{}
	mutex := &sync.Mutex{}

	cpuPercentTotal := float64(0)
	memPercentTotal := float32(0)
	memRSSTotal := uint64(0)
	memVMSTotal := uint64(0)
	memSwapTotal := uint64(0)

	openFilesAll := []string{}
	processCount := 0

	for _, pid := range m.procPIDs {
		procWg.Add(1)
		go func(pid int32, procWg *sync.WaitGroup, mutex *sync.Mutex) {
			defer procWg.Done()
			p, err := process.NewProcess(pid)
			if err != nil {
				return
			}
			cpuPercent, err := p.CPUPercent()
			if err != nil {
				return
			}
			memPercent, err := p.MemoryPercent()
			if err != nil {
				return
			}
			memInfo, err := p.MemoryInfo()
			if err != nil {
				return
			}
			files, err := openFiles(pid)
			if err != nil {
				return
			}
			mutex.Lock()

			cpuPercentTotal = cpuPercentTotal + cpuPercent
			memPercentTotal = memPercentTotal + memPercent
			memRSSTotal = memRSSTotal + memInfo.RSS
			memVMSTotal = memVMSTotal + memInfo.VMS
			memSwapTotal = memSwapTotal + memInfo.Swap
			openFilesAll = append(openFilesAll, files...)
			processCount = processCount + 1

			mutex.Unlock()
		}(pid, procWg, mutex)
	}
	procWg.Wait()

	m.Store("proc_cpu", cpuPercentTotal)
	m.Store("proc_mem", memPercentTotal)
	m.Store("proc_rss", memRSSTotal)
	m.Store("proc_vms", memVMSTotal)
	m.Store("proc_swap", memSwapTotal)
	m.Store("proc_open_files", len(uniqueSlice(openFilesAll)))

	m.Store("proc_count", processCount)
}

// openFiles
// reference https://github.com/shirou/gopsutil/blob/2f74cb51781d173ae257a61e955594064ab6f7b0/process/process_linux.go#L781-L788
func openFiles(pid int32) ([]string, error) {
	fdPath := hostProc(strconv.Itoa(int(pid)), "fd")
	d, err := os.Open(fdPath)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	fnames, err := d.Readdirnames(-1)
	return fnames, err
}

// copy from gopsutil/internal/common
func hostProc(combineWith ...string) string {
	value := os.Getenv("HOST_PROC")
	if value == "" {
		value = "/proc"
	}
	switch len(combineWith) {
	case 0:
		return value
	case 1:
		return filepath.Join(value, combineWith[0])
	default:
		all := make([]string, len(combineWith)+1)
		all[0] = value
		copy(all[1:], combineWith)
		return filepath.Join(all...)
	}
}

func uniqueSlice(in []string) []string {
	results := make([]string, 0, len(in))
	encountered := map[string]bool{}
	for i := 0; i < len(in); i++ {
		if !encountered[in[i]] {
			encountered[in[i]] = true
			results = append(results, in[i])
		}
	}
	return results
}
