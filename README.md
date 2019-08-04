# metr

`metr` gets system **metr**ics

## Usage

### metr cond

``` console
$ metr cond 'cpu > 10 or mem > 90'
```

if condition match `exit 0` else `exit 1`.

### metr get

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
