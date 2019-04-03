# dbanon

![](https://travis-ci.org/mpchadwick/dbanon.svg?branch=master)

## Usage

```
mysqldump --complete-insert mydb | dbanon | gzip > mydb.sql.gz
```

## Limitations

Currently requires `mysqldump` be run with `--complete-insert` flag.