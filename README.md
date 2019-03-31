# dbanon

## Usage

```
mysqldump --complete-insert mydb | dbanon | gzip > mydb.sql.gz
```

## Limitations

Currently requires `mysqldump` be run with `--complete-insert` flag.