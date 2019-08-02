# metr

`metr` get system metrics

## Usage

### metr cond

``` console
$ metr cond 'cpu > 10 or mem 90'
```

if condition match `exit 0` else `exit 1`.

### [WIP] metr get

``` console
$ metr get all
irq:0
load15:2.55
mem:71.0200548171997
idle:87.23796606135554
iowait:0
softirq:0
steal:0
guest:0
nice:0
stolen:0
load1:2.21
load5:2.73
cpu:0
swap:7270301696
user:7.6438811331984535
system:5.118152805446011
guest_nice:0
```
