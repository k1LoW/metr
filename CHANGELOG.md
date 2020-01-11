# Changelog

## [v0.10.0](https://github.com/k1LoW/metr/compare/v0.9.1...v0.10.0) (2020-01-11)

* Remove completion script when remove metr [#26](https://github.com/k1LoW/metr/pull/26) ([k1LoW](https://github.com/k1LoW))
* Install completion script when install metr via deb/rpm [#25](https://github.com/k1LoW/metr/pull/25) ([k1LoW](https://github.com/k1LoW))
* Add `metr completion` [#24](https://github.com/k1LoW/metr/pull/24) ([k1LoW](https://github.com/k1LoW))

## [v0.9.1](https://github.com/k1LoW/metr/compare/v0.9.0...v0.9.1) (2019-10-27)

* Update spf13/pflag to v1.0.5 [#23](https://github.com/k1LoW/metr/pull/23) ([k1LoW](https://github.com/k1LoW))

## [v0.9.0](https://github.com/k1LoW/metr/compare/v0.8.0...v0.9.0) (2019-09-18)

* Add `--log-dir` option [#22](https://github.com/k1LoW/metr/pull/22) ([k1LoW](https://github.com/k1LoW))
* Fix: metr cond does not have --name option [#21](https://github.com/k1LoW/metr/pull/21) ([k1LoW](https://github.com/k1LoW))

## [v0.8.0](https://github.com/k1LoW/metr/compare/v0.7.0...v0.8.0) (2019-09-02)

* [BREAKING] Remove `proc_connections` / Fix `proc_open_files` counting logic [#20](https://github.com/k1LoW/metr/pull/20) ([k1LoW](https://github.com/k1LoW))

## [v0.7.0](https://github.com/k1LoW/metr/compare/v0.6.0...v0.7.0) (2019-08-26)

* [BREAKING] Show metrics `proc_*` even if no process exists. [#19](https://github.com/k1LoW/metr/pull/19) ([k1LoW](https://github.com/k1LoW))

## [v0.6.0](https://github.com/k1LoW/metr/compare/v0.5.1...v0.6.0) (2019-08-24)

* Add `--name` option for process metrics [#18](https://github.com/k1LoW/metr/pull/18) ([k1LoW](https://github.com/k1LoW))
* metr/metrics support collecting multi processes metrics (use total) [#17](https://github.com/k1LoW/metr/pull/17) ([k1LoW](https://github.com/k1LoW))

## [v0.5.1](https://github.com/k1LoW/metr/compare/v0.5.0...v0.5.1) (2019-08-10)

* Add `metr test` command that alias for `metr cond` [#16](https://github.com/k1LoW/metr/pull/16) ([k1LoW](https://github.com/k1LoW))

## [v0.5.0](https://github.com/k1LoW/metr/compare/v0.4.0...v0.5.0) (2019-08-08)

* Add `--pid` option for process metrics [#15](https://github.com/k1LoW/metr/pull/15) ([k1LoW](https://github.com/k1LoW))

## [v0.4.0](https://github.com/k1LoW/metr/compare/v0.3.0...v0.4.0) (2019-08-07)

* Add metric `numcpu` [#14](https://github.com/k1LoW/metr/pull/14) ([k1LoW](https://github.com/k1LoW))

## [v0.3.0](https://github.com/k1LoW/metr/compare/v0.2.2...v0.3.0) (2019-08-06)

* Handle error [#13](https://github.com/k1LoW/metr/pull/13) ([k1LoW](https://github.com/k1LoW))

## [v0.2.2](https://github.com/k1LoW/metr/compare/v0.2.1...v0.2.2) (2019-08-06)

* Add test [#12](https://github.com/k1LoW/metr/pull/12) ([k1LoW](https://github.com/k1LoW))

## [v0.2.1](https://github.com/k1LoW/metr/compare/v0.2.0...v0.2.1) (2019-08-05)

* Change bindir (deb, rpm) `/usr/local/bin` -> `/usr/bin` [#11](https://github.com/k1LoW/metr/pull/11) ([k1LoW](https://github.com/k1LoW))

## [v0.2.0](https://github.com/k1LoW/metr/compare/v0.1.1...v0.2.0) (2019-08-05)

* Add `metr check` [#10](https://github.com/k1LoW/metr/pull/10) ([k1LoW](https://github.com/k1LoW))

## [v0.1.1](https://github.com/k1LoW/metr/compare/v0.1.0...v0.1.1) (2019-08-05)

* Fix broken `metr cond` [#9](https://github.com/k1LoW/metr/pull/9) ([k1LoW](https://github.com/k1LoW))

## [v0.1.0](https://github.com/k1LoW/metr/compare/v0.0.4...v0.1.0) (2019-08-04)

* Set default `--interval` to 500ms [#8](https://github.com/k1LoW/metr/pull/8) ([k1LoW](https://github.com/k1LoW))
* Set available metrics for each OS [#7](https://github.com/k1LoW/metr/pull/7) ([k1LoW](https://github.com/k1LoW))
* Change `metr get` output [#6](https://github.com/k1LoW/metr/pull/6) ([k1LoW](https://github.com/k1LoW))
* Add `metr list` [#5](https://github.com/k1LoW/metr/pull/5) ([k1LoW](https://github.com/k1LoW))
* Use sync.Map with ordered keys [#4](https://github.com/k1LoW/metr/pull/4) ([k1LoW](https://github.com/k1LoW))

## [v0.0.4](https://github.com/k1LoW/metr/compare/v0.0.3...v0.0.4) (2019-08-03)

* Fix NaN when the interval value is small [#3](https://github.com/k1LoW/metr/pull/3) ([k1LoW](https://github.com/k1LoW))

## [v0.0.3](https://github.com/k1LoW/metr/compare/v0.0.2...v0.0.3) (2019-08-03)

* Fix map violate via goroutine [#2](https://github.com/k1LoW/metr/pull/2) ([k1LoW](https://github.com/k1LoW))

## [v0.0.2](https://github.com/k1LoW/metr/compare/v0.0.1...v0.0.2) (2019-08-02)

* Add --interval option [#1](https://github.com/k1LoW/metr/pull/1) ([k1LoW](https://github.com/k1LoW))

## [v0.0.1](https://github.com/k1LoW/metr/compare/4eeada302c57...v0.0.1) (2019-08-02)

