# BRUM

usage

```shell
while true; do echo "$(mpstat 1 1 | awk 'END{print 100-$NF}')"; sleep 0.01; done | go run brum/main.go 
```
