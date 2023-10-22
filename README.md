# BRUM

install

```shell
git clone https://github.com/Coding-Brownies/brum
cd brum
go install cmd/brum.go
```

usage

```shell
while true; do echo `mpstat 1 1 | awk 'END{print 100-$NF}'`; done | brum
```


to have fun:

```shell
# in a shell -----
while true; do echo `mpstat 1 1 | awk 'END{print 100-$NF}'`; done | brum

# in another -----
curl -o rock-you.txt https://github.com/brannondorsey/naive-hashcat/releases/download/data/rockyou.txt
echo '$2b$10$iIp9E/euW0r6ErekR4y8suX65xW9nkO8y7XOVso9EZxlb5gCgUUUu' > hash.txt
hashcat -a 0 -m 3200 hash.txt rock-you.txt --force
```
