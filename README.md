# dbanon

[![Build Status](https://travis-ci.org/mpchadwick/dbanon.svg?branch=master)](https://travis-ci.org/mpchadwick/dbanon) [![codecov](https://codecov.io/gh/mpchadwick/dbanon/branch/master/graph/badge.svg)](https://codecov.io/gh/mpchadwick/dbanon)

A run-anywhere, dependency-less database anonymizer.

## Installation

Download [the latest release from GitHub](https://github.com/mpchadwick/dbanon/releases).

## Usage

`dbanon` reads from `stdin` and writes to `stdout`.

```
mysqldump mydb | dbanon -config=myconfig.yml | gzip > mydb.sql.gz
```

The `-config` flag can use bundled configurations or point to the path of a custom configuration file. 

### Configuration

Specify the path to your config file via the `-config` flag

```
mysqldump mydb | dbanon -config=myconfig.yml | gzip > mydb.sql.gz
```

See [the `etc` directory](etc/) for examples.

Columns are specified as key / value pairs. The value string winds up getting passed to [this function](https://github.com/mpchadwick/dbanon/blob/ade634a10bc282c06fecef115afbdd6661a94277/src/provider.go#L36), which gets random values from [`dmgk/faker`](https://github.com/dmgk/faker).

It is also possible to pass direct Faker function calls for [supported "raw providers"](https://github.com/mpchadwick/dbanon/blob/ade634a10bc282c06fecef115afbdd6661a94277/src/provider.go#L13-L17)

## Logging

`dbanon` records messages about anything notable (e.g. invalid configuration) to the file `dbanon.log` in the directory from which you run it.

**`-log-file`**

The `-log-file` flag can be used to have `dbanon` log to a different location.

```
mysqldump mydb | dbanon -config=myconfig.yml -log-file=var/dbanon.log
```

**`-log-level`**

The `-log-level` flag can be used to control the verbosity of logs. Supported values can be found [here](https://github.com/sirupsen/logrus/blob/d131c24e23baaa812461202af6d7cfa388e2d292/logrus.go#L25-L45).

```
mysqldump mydb | dbanon -config=myconfig.yml -log-level=debug | gzip > mydb.sql.gz
```

The default log level is `info`.

**`-silent`**

Logging can be disabled entirely by passing the `-silent` flag to `dbanon`

```
mysqldump mydb | dbanon -config=myconfig.yml -silent | gzip > mydb.sql.gz
```

## Profiling

`dbanon` will generate a CPU profile to the file `dbanon.prof` when passed the `-profile` flag.

```
mysqldump mydb | dbanon -profile -config=myconfig.yml >/dev/null
```

## Limitations

- Currently only supports MySQL

## Updating

`dbanon` will self-update when passed the `-update` flag

```
dbanon -update
```