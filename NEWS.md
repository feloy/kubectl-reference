# News

## New in v1.19

- new common option `field-manager`
* new option `port` in `create/deployment`
* new option `replicas` in `create/deployment`
* new option `privileged` in `run`
- new option `dry-run` in `scale`
- new option `selector` in `diff`
* new option `list` in `annotate`
* new option `copy-to` in `alpha/debug`
* new option `replace` in `alpha/debug`
* new option `same-node` in `alpha/debug`
* new option `share-processes` in `alpha/debug`
- new arguments `COMMAND args...` to `create deployment`
- removed option `generator` from `create/clusterrolebinding`
- removed option `export` from `get`
- removed option `heapster-namespace` from `top/pod`
- removed option `heapster-port` from from `top/pod`
- removed option `heapster-scheme` from `top/pod`
- removed option `heapster-service` from `top/pod`
- removed option `server-dry-run` from `apply`

## New in v1.18

- `dry-run` options accept `server`/`client`/`none` instead of `true`/`false`
- new option `dry-run` in `delete`
- new option `dry-run` in `replace`
- new option `dry-run` in `taint`
- new option `filename` in `exec`
- new option `disable-eviction` in `drain`
- new option `skip-wait-for-delete-timeout` in `drain`
- new command `alpha debug`
