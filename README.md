# metr [![Build Status](https://travis-ci.org/k1LoW/metr.svg?branch=master)](https://travis-ci.org/k1LoW/metr) [![GitHub release](https://img.shields.io/github/release/k1LoW/metr.svg)](https://github.com/k1LoW/metr/releases)

`metr` provides an easy way to use host/process **metr**ics for shell script/monitoring tool.

## Usage

### `metr cond` ( alias: `metr test` )

``` console
$ metr cond 'cpu > 10 or mem > 90'
```

if condition match `exit 0` else `exit 1`, like `test` command.

#### Available Operators

`+`, `-`, `*`, `/`, `==`, `!=`, `<`, `>`, `<=`, `>=`, `not`, `and`, `or`, `!`, `&&`, `||`

### `metr check`

``` console
$ metr check -w 'cpu > 10 or mem > 50' -c 'cpu > 50 and mem > 90'
METR WARNING: w(cpu > 10 or mem > 50) c(cpu > 50 and mem > 90)
```

`metr check` is compatible with

- Nagios plugin
- Mackerel check plugin `command`
- Consul check `command`
- Sensu checks

| Exit status code | Meaning  |
| ---------------- | -------- |
| 0	               | OK       |
| 1                | WARNING  |
| 2                | CRITICAL |
| 3                | UNKNOWN  |

### `metr get`

``` console
$ metr get all
cpu:4.239401
mem:64.362717
swap:61603840
user:2.525253
system:1.515152
idle:95.959596
nice:0.000000
load1:1.310000
load5:1.410000
load15:1.550000
numcpu:8
$ metr get all -p `pgrep -n docker`
proc_cpu:0.079323
proc_mem:0.230384
proc_rss:39579648
proc_vms:4903047168
proc_swap:0
proc_connections:0
cpu:0.000000
mem:59.643674
swap:781451264
user:3.719606
system:2.253192
idle:94.027202
nice:0.000000
load1:1.560000
load5:1.720000
load15:1.510000
numcpu:8
$ metr get cpu
3.241895
```

### `metr list`

``` console
$ metr list
cpu (now:33.084577 %): Percentage of cpu used.
mem (now:66.468358 %): Percentage of RAM used.
swap (now:875823104 bytes): Amount of memory that has been swapped out to disk (bytes).
user (now:18.610422 %): Percentage of CPU utilization that occurred while executing at the user level.
system (now:14.143921 %): Percentage of CPU utilization that occurred while executing at the system level.
idle (now:67.245658 %): Percentage of time that CPUs were idle and the system did not have an outstanding disk I/O request.
nice (now:0.000000 %): Percentage of CPU utilization that occurred while executing at the user level with nice priority.
load1 (now:3.640000 ): Load avarage for 1 minute.
load5 (now:4.210000 ): Load avarage for 5 minutes.
load15 (now:4.600000 ): Load avarage for 15 minutes.
numcpu (now:8 ): Number of logical CPUs.
(metric measurement interval: 500 ms)
$ metr list -p `pgrep -n docker`
proc_cpu (now:1.820857 %): Percentage of the CPU time the process uses.
proc_mem (now:1.264739 %): Percentage of the total RAM the process uses.
proc_rss (now:217280512 bytes): Non-swapped physical memory the process uses (bytes).
proc_vms (now:7010299904 bytes): Amount of virtual memory the process uses (bytes).
proc_swap (now:0 bytes): Amount of memory that has been swapped out to disk the process uses (bytes).
proc_connections (now:0 ): Amount of connections(TCP, UDP or UNIX) the process uses.
cpu (now:22.000000 %): Percentage of cpu used.
mem (now:59.768772 %): Percentage of RAM used.
swap (now:781451264 bytes): Amount of memory that has been swapped out to disk (bytes).
user (now:14.925373 %): Percentage of CPU utilization that occurred while executing at the user level.
system (now:6.467662 %): Percentage of CPU utilization that occurred while executing at the system level.
idle (now:78.606965 %): Percentage of time that CPUs were idle and the system did not have an outstanding disk I/O request.
nice (now:0.000000 %): Percentage of CPU utilization that occurred while executing at the user level with nice priority.
load1 (now:1.360000 ): Load avarage for 1 minute.
load5 (now:1.610000 ): Load avarage for 5 minutes.
load15 (now:1.490000 ): Load avarage for 15 minutes.
numcpu (now:8 ): Number of logical CPUs.
(metric measurement interval: 500 ms)
```

## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/metr
```

**manually:**

Download binany/deb/rpm from [releases page](https://github.com/k1LoW/metr/releases)

**go get:**

```console
$ go get github.com/k1LoW/metr
```
