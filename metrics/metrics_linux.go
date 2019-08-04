// +build linux

package metrics

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func availableMetrics() []Metric {
	return []Metric{
		{"cpu", "Percentage of cpu used.", "%f"},
		{"mem", "Percentage of RAM used.", "%f"},
		{"swap", "Amount of memory that has been swapped out to disk (bytes).", "%d"},

		{"user", "Percentage of CPU utilization that occurred while executing at the user level.", "%f"},
		{"system", "Percentage of CPU utilization that occurred while executing at the system level.", "%f"},
		{"idle", "Percentage of time that CPUs were idle and the system did not have an outstanding disk I/O request.", "%f"},
		{"nice", "Percentage of CPU utilization that occurred while executing at the user level with nice priority.", "%f"},

		{"iowait", "Percentage of time that CPUs were idle during which the system had an outstanding disk I/O request.", "%f"},
		{"irq", "Percentage of time spent by CPUs to service hardware interrupts.", "%f"},
		{"softirq", "Percentage of time spent by CPUs to service software interrupts.", "%f"},
		{"steal", "Percentage of time spent in involuntary wait by the virtual CPUs while the hypervisor was servicing another virtual processor.", "%f"},
		{"guest", "Percentage of time spent by CPUs to run a virtual processor.", "%f"},
		{"guest_nice", "Percentage of time spent by CPUs to run a virtual processor with nice priority.", "%f"},

		{"load1", "Load avarage for 1 minute.", "%f"},
		{"load5", "Load avarage for 5 minutes.", "%f"},
		{"load15", "Load avarage for 15 minutes.", "%f"},
	}
}

func getMetrics(i time.Duration) (*Metrics, error) {
	wg := &sync.WaitGroup{}

	m := NewMetrics()

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		hostCpuPercent, err := cpu.Percent(i, false)
		if err != nil {
			// TODO
		}
		m.Store("cpu", hostCpuPercent[0])
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
		ts1, err := cpu.Times(false)
		if err != nil {
			// TODO
		}
		total1 := ts1[0].Total()

		if i == 0 {
			m.Store("user", ts1[0].User/total1*100)
			m.Store("system", ts1[0].System/total1*100)
			m.Store("idle", ts1[0].Idle/total1*100)
			m.Store("nice", ts1[0].Nice/total1*100)
			m.Store("iowait", ts1[0].Iowait/total1*100)
			m.Store("irq", ts1[0].Irq/total1*100)
			m.Store("softirq", ts1[0].Softirq/total1*100)
			m.Store("steal", ts1[0].Steal/total1*100)
			m.Store("guest", ts1[0].Guest/total1*100)
			m.Store("guest_nice", ts1[0].GuestNice/total1*100)
			return
		}

		time.Sleep(i)

		ts2, err := cpu.Times(false)
		if err != nil {
			// TODO
		}
		total2 := ts2[0].Total()

		total := total2 - total1
		if total == 0 {
			m.Store("user", ts2[0].User/total2*100)
			m.Store("system", ts2[0].System/total2*100)
			m.Store("idle", ts2[0].Idle/total2*100)
			m.Store("nice", ts2[0].Nice/total2*100)
			m.Store("iowait", ts2[0].Iowait/total2*100)
			m.Store("irq", ts2[0].Irq/total2*100)
			m.Store("softirq", ts2[0].Softirq/total2*100)
			m.Store("steal", ts2[0].Steal/total2*100)
			m.Store("guest", ts2[0].Guest/total2*100)
			m.Store("guest_nice", ts2[0].GuestNice/total2*100)
			return
		}

		m.Store("user", (ts2[0].User-ts1[0].User)/total*100)
		m.Store("system", (ts2[0].System-ts1[0].System)/total*100)
		m.Store("idle", (ts2[0].Idle-ts1[0].Idle)/total*100)
		m.Store("nice", (ts2[0].Nice-ts1[0].Nice)/total*100)
		m.Store("iowait", (ts2[0].Iowait-ts1[0].Iowait)/total*100)
		m.Store("irq", (ts2[0].Irq-ts1[0].Irq)/total*100)
		m.Store("softirq", (ts2[0].Softirq-ts1[0].Softirq)/total*100)
		m.Store("steal", (ts2[0].Steal-ts1[0].Steal)/total*100)
		m.Store("guest", (ts2[0].Guest-ts1[0].Guest)/total*100)
		m.Store("guest_nice", (ts2[0].GuestNice-ts1[0].GuestNice)/total*100)
	}(wg)

	l, err := load.Avg()
	if err != nil {
		return nil, err
	}
	m.Store("load1", l.Load1)
	m.Store("load5", l.Load5)
	m.Store("load15", l.Load15)

	wg.Wait()

	return m, nil
}
