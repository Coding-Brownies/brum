# BRUM

usage

```shell
while true; do echo "$(mpstat 1 1 | awk 'END{print 100-$NF}')"; sleep 0.01; done | go run brum/main.go 
```


to have fun:

```shell
# in a shell
while true; do echo "$(mpstat 1 1 | awk 'END{print 100-$NF}')"; sleep 0.01; done | go run brum/main.go 

# in another 
hashcat -a 0 -m 3200 hash.txt ~/Downloads/rockyou.txt --force
```