// +build linux

package metrics

import (
	"runtime"
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

		{"iowait", "Percentage of time that CPUs were idle during which the system had an outstanding disk I/O request.", "%f", "%"},
		{"irq", "Percentage of time spent by CPUs to service hardware interrupts.", "%f", "%"},
		{"softirq", "Percentage of time spent by CPUs to service software interrupts.", "%f", "%"},
		{"steal", "Percentage of time spent in involuntary wait by the virtual CPUs while the hypervisor was servicing another virtual processor.", "%f", "%"},
		{"guest", "Percentage of time spent by CPUs to run a virtual processor.", "%f", "%"},
		{"guest_nice", "Percentage of time spent by CPUs to run a virtual processor with nice priority.", "%f", "%"},

		{"load1", "Load avarage for 1 minute.", "%f", ""},
		{"load5", "Load avarage for 5 minutes.", "%f", ""},
		{"load15", "Load avarage for 15 minutes.", "%f", ""},

		{"numcpu", "Number of logical CPUs.", "%d", ""},
	}
}

func (m *Metrics) Collect() error {
	wg := &sync.WaitGroup{}

	// 2 = goroutine count
	errChan := make(chan error, 2)

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
