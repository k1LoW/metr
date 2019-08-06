# metr [![Build Status](https://travis-ci.org/k1LoW/metr.svg?branch=master)](https://travis-ci.org/k1LoW/metr) [![GitHub release](https://img.shields.io/github/release/k1LoW/metr.svg)](https://github.com/k1LoW/metr/releases)

`metr` provides an easy way to use system **metr**ics for shell script/monitoring tool.

## Usage

### `metr cond`

``` console
$ metr cond 'cpu > 10 or mem > 90'
```

if condition match `exit 0` else `exit 1`, like `test` command.

#### Available Operators

`+`, `-`, `*`, `/`, `==`, `!=`, `<`, `>`, `<=`, `>=`, `not`, `and`, `or`, `!`, `&&`, `||`

### `metr check`

``` console
$ metr check -w 'cpu > 10 or mem > 50' -c 'cpu > 50 and mem > 90'
```

`metr check` is compatible with

- Nagios plugin
- Mackerel check plugin `command`
- Sensu checks

| Exit status code | Meaning  |
| ---------------- | -------- |
| 0	               | OK       |
| 1                | WARNING  |
| 2                | CRITICAL |
| 3                | UNKNOWN  |

### `metr get`

``` console
$ metr get all -i 500
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
$ metr get cpu -i 500
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
(metric measurement interval: 500 ms)
```

## Install

**homebrew tap:**

```console
$ brew install k1LoW/tap/metr
```

**manually:**

Download binany from [releases page](https://github.com/k1LoW/metr/releases)

**go get:**

```console
$ go get github.com/k1LoW/metr
```
