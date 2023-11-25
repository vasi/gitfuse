# gitfuse
Demo of poking at git in a FUSE filesystem. Exposes the author of each commit hash.

# Usage

```
$ mkdir mnt
$ go run main.go . mnt
# In another terminal
$ cat mnt/mnt/a5417c87d3db74901a10eae3f74af230811d2886
Dave Vasilevsky
```
