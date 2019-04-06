# dbanon

[![Build Status](https://travis-ci.org/mpchadwick/dbanon.svg?branch=master)](https://travis-ci.org/mpchadwick/dbanon)

A run-anywhere, dependency-less database anonymizer.

## Usage

`dbanon` reads from `stdin` and writes to `stdout`.

```
mysqldump --complete-insert mydb | dbanon -config=myconfig.yml | gzip > mydb.sql.gz
```

The `-config` flag can use bundled configurations (magento1, magento2) or point to the path of a custom configuration file. See [the `etc` directory](etc/) for examples

### Magento EAV

In order to anonymize a Magento database and allow a user friendly config, `dbanon` must first map EAV attribute codes to their respective attribute IDs. This can be done with the `map-eav` subcommand:

```
mysqldump --complete-insert mydb eav_entity_type eav_attribute | dbanon -config=magento2 map-eav > ~/magento2.yml
```

`map-eav` will replace the attribute codes in the config file with attribute ids. Next, run `dbanon` with the config file generated via `map-eav`

```
mysqldump --complete-insert mydb | dbanon -config=~/magento2.yml | gzip > mydb.sql.gz
```


## Limitations

Currently requires `mysqldump` be run with `--complete-insert` flag.