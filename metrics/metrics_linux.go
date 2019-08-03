// +build linux

package metrics

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func getMetrics(i time.Duration) (map[string]interface{}, error) {
	mutex := new(sync.Mutex)
	wg := &sync.WaitGroup{}

	m := map[string]interface{}{
		"cpu":        0,
		"mem":        0,
		"swap":       0,
		"user":       0,
		"system":     0,
		"idle":       0,
		"nice":       0,
		"iowait":     0,
		"irq":        0,
		"softirq":    0,
		"steal":      0,
		"guest":      0,
		"guest_nice": 0,
		"stolen":     0,
		"load1":      0,
		"load5":      0,
		"load15":     0,
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		hostCpuPercent, err := cpu.Percent(i, false)
		if err != nil {
			// TODO
		}
		mutex.Lock()
		m["cpu"] = hostCpuPercent[0]
		mutex.Unlock()
		wg.Done()
	}(wg)

	vm, err := mem.VirtualMemory()
	if err != nil {
		return map[string]interface{}{}, err
	}
	mutex.Lock()
	m["mem"] = vm.UsedPercent
	mutex.Unlock()

	sw, err := mem.SwapMemory()
	if err != nil {
		return map[string]interface{}{}, err
	}
	mutex.Lock()
	m["swap"] = sw.Used
	mutex.Unlock()

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		ts1, err := cpu.Times(false)
		if err != nil {
			// TODO
		}
		total1 := ts1[0].Total()

		if i == 0 {
			mutex.Lock()
			m["user"] = ts1[0].User / total1 * 100
			m["system"] = ts1[0].System / total1 * 100
			m["idle"] = ts1[0].Idle / total1 * 100
			m["nice"] = ts1[0].Nice / total1 * 100
			m["iowait"] = ts1[0].Iowait / total1 * 100
			m["irq"] = ts1[0].Irq / total1 * 100
			m["softirq"] = ts1[0].Softirq / total1 * 100
			m["steal"] = ts1[0].Steal / total1 * 100
			m["guest"] = ts1[0].Guest / total1 * 100
			m["guest_nice"] = ts1[0].GuestNice / total1 * 100
			m["stolen"] = ts1[0].Stolen / total1 * 100
			mutex.Unlock()
			return
		}

		time.Sleep(i)

		ts2, err := cpu.Times(false)
		if err != nil {
			// TODO
		}
		total2 := ts2[0].Total()

		total := total2 - total1
		mutex.Lock()
		m["user"] = (ts2[0].User - ts1[0].User) / total * 100
		m["system"] = (ts2[0].System - ts1[0].System) / total * 100
		m["idle"] = (ts2[0].Idle - ts1[0].Idle) / total * 100
		m["nice"] = (ts2[0].Nice - ts1[0].Nice) / total * 100
		m["iowait"] = (ts2[0].Iowait - ts1[0].Iowait) / total * 100
		m["irq"] = (ts2[0].Irq - ts1[0].Irq) / total * 100
		m["softirq"] = (ts2[0].Softirq - ts1[0].Softirq) / total * 100
		m["steal"] = (ts2[0].Steal - ts1[0].Steal) / total * 100
		m["guest"] = (ts2[0].Guest - ts1[0].Guest) / total * 100
		m["guest_nice"] = (ts2[0].GuestNice - ts1[0].GuestNice) / total * 100
		m["stolen"] = (ts2[0].Stolen - ts1[0].Stolen) / total * 100
		mutex.Unlock()
	}(wg)

	l, err := load.Avg()
	if err != nil {
		return map[string]interface{}{}, err
	}
	mutex.Lock()
	m["load1"] = l.Load1
	m["load5"] = l.Load5
	m["load15"] = l.Load15
	mutex.Unlock()

	wg.Wait()

	return m, nil
}
