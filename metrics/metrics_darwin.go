// +build darwin

package metrics

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const calcInterval = 0

func getMetrics() (map[string]interface{}, error) {
	hostCpuPercent, err := cpu.Percent(calcInterval, false)
	if err != nil {
		return map[string]interface{}{}, err
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return map[string]interface{}{}, err
	}

	sw, err := mem.SwapMemory()
	if err != nil {
		return map[string]interface{}{}, err
	}

	ts, err := cpu.Times(false)
	if err != nil {
		return map[string]interface{}{}, err
	}
	total := ts[0].Total()

	l, err := load.Avg()
	if err != nil {
		return map[string]interface{}{}, err
	}

	m := map[string]interface{}{
		"cpu":        hostCpuPercent[0],
		"mem":        vm.UsedPercent,
		"swap":       sw.Used,
		"user":       ts[0].User / total * 100,
		"system":     ts[0].System / total * 100,
		"idle":       ts[0].Idle / total * 100,
		"nice":       ts[0].Nice / total * 100,
		"iowait":     ts[0].Iowait / total * 100,
		"irq":        ts[0].Irq / total * 100,
		"softirq":    ts[0].Softirq / total * 100,
		"steal":      ts[0].Steal / total * 100,
		"guest":      ts[0].Guest / total * 100,
		"guest_nice": ts[0].GuestNice / total * 100,
		"stolen":     ts[0].Stolen / total * 100,
		"load1":      l.Load1,
		"load5":      l.Load5,
		"load15":     l.Load15,
	}

	return m, nil
}
