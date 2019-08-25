// +build darwin

package metrics

import (
	"runtime"
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
		{"proc_connections", "Amount of connections(TCP, UDP or UNIX) the process uses.", "%d", "", int(0)},
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
			return
		}

		m.Store("user", (after[0].User-before[0].User)/total*100)
		m.Store("system", (after[0].System-before[0].System)/total*100)
		m.Store("idle", (after[0].Idle-before[0].Idle)/total*100)
		m.Store("nice", (after[0].Nice-before[0].Nice)/total*100)
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
	connectionTotal := 0
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
			connections, err := p.Connections()
			if err != nil {
				return
			}
			mutex.Lock()

			cpuPercentTotal = cpuPercentTotal + cpuPercent
			memPercentTotal = memPercentTotal + memPercent
			memRSSTotal = memRSSTotal + memInfo.RSS
			memVMSTotal = memVMSTotal + memInfo.VMS
			memSwapTotal = memSwapTotal + memInfo.Swap
			connectionTotal = connectionTotal + len(connections)
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
	m.Store("proc_connections", connectionTotal)
	m.Store("proc_count", processCount)
}
