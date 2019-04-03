# dbanon

[![Build Status](https://travis-ci.org/mpchadwick/dbanon.svg?branch=master)](https://travis-ci.org/mpchadwick/dbanon)

A run-anywhere, dependency-less database anonymizer.

## Usage

`dbanon` reads from `stdin` and writes to `stdout`.

```
mysqldump --complete-insert mydb | dbanon -config=magento2 | gzip > mydb.sql.gz
```

The `-config` flag can use bundled configurations (magento1, magento2) or point to the path of a custom configuration file. See [the `etc` directory](etc/) for examples

## Limitations

Currently requires `mysqldump` be run with `--complete-insert` flag.