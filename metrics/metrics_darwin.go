// +build darwin

package metrics

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func AvailableMetrics() []Metric {
	return []Metric{
		{"cpu", "Percentage of cpu used.", "%f", "%"},
		{"mem", "Percentage of RAM used.", "%f", "%"},
		{"swap", "Amount of memory that has been swapped out to disk (bytes).", "%d", "bytes"},

		{"user", "Percentage of CPU utilization that occurred while executing at the user level.", "%f", "%"},
		{"system", "Percentage of CPU utilization that occurred while executing at the system level.", "%f", "%"},
		{"idle", "Percentage of time that CPUs were idle and the system did not have an outstanding disk I/O request.", "%f", "%"},
		{"nice", "Percentage of CPU utilization that occurred while executing at the user level with nice priority.", "%f", "%"},

		{"load1", "Load avarage for 1 minute.", "%f", ""},
		{"load5", "Load avarage for 5 minutes.", "%f", ""},
		{"load15", "Load avarage for 15 minutes.", "%f", ""},
	}
}

func (m *Metrics) Collect() (*Metrics, error) {
	wg := &sync.WaitGroup{}

	// 2 = goroutine count
	errChan := make(chan error, 2)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		cpuPercent, err := cpu.Percent(m.interval, false)
		if err != nil {
			errChan <- err
			return
		}
		m.Store("cpu", cpuPercent[0])
		wg.Done()
	}(wg)

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	m.Store("mem", vm.UsedPercent)
	sw, err := mem.SwapMemory()
	if err != nil {
		return nil, err
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

		if m.interval == 0 {
			m.Store("user", before[0].User/beforeTotal*100)
			m.Store("system", before[0].System/beforeTotal*100)
			m.Store("idle", before[0].Idle/beforeTotal*100)
			m.Store("nice", before[0].Nice/beforeTotal*100)
			return
		}

		time.Sleep(m.interval)

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
		return nil, err
	}
	m.Store("load1", l.Load1)
	m.Store("load5", l.Load5)
	m.Store("load15", l.Load15)

	wg.Wait()

	select {
	case err := <-errChan:
		return nil, err
	default:
	}

	return m, nil
}
