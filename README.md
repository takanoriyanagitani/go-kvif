# go-kvif
Key/Value interface with official implementations

[![Go Reference](https://pkg.go.dev/badge/github.com/takanoriyanagitani/go-kvif.svg)](https://pkg.go.dev/github.com/takanoriyanagitani/go-kvif)
[![Go Report Card](https://goreportcard.com/badge/github.com/takanoriyanagitani/go-kvif)](https://goreportcard.com/report/github.com/takanoriyanagitani/go-kvif)
[![codecov](https://codecov.io/gh/takanoriyanagitani/go-kvif/branch/main/graph/badge.svg?token=T4X9DP2KSH)](https://codecov.io/gh/takanoriyanagitani/go-kvif)

## Overview

#### Archive Type

###### Zip

- bucket: path/to/sample1/a.zip
  - path/inside/archive/file1.dat
  - path/inside/archive/deep/file2.gz
  - ...
- bucket: path/to/sample2/date/2022-09-24/device/cafef00d-dead-beaf-face-864299792458/time/06-04.zip
  - 06-04-00.txt
  - 06-04-01.txt
  - 06-04-02.txt
  - ...
  - 06-04-59.txt
  - 06-04-60.txt

#### Sql Type

###### PostgreSQL

- bucket: postgres.public.sample1
  - key1
  - key2
  - ...
- bucket: postgres.public.date_2022_09_24_device_cafef00ddeadbeafface864299792458
  - 06_04_00
  - 06_04_01
  - 06_04_02
  - ...
  - 06_04_60
