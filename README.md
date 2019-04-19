# mempig
Mempig's only job is to eat up your system's memory so you can find answers to questions like these:
* How much memory can mempig ask for before it gets killed?
* Does the system slow to a crawl or lock up completely for large memory requests? 
* If mempig is running in a container, how much memory can mempig eat before the container crashes?

## Observed behavior on Mac

On my MacBook Pro with 16G of memory, it works fine to ask for 8G:
```
$ mempig -G 8
Allocated 8GiB (8589934592 bytes)
Wallowing in the memory. Press ctrl-C to quit.
^C
```

I tried asking for 100G and that was fine too; it just got the memory by swapping out to disk.

But to my surprise, things changed when I tried running mempig inside a docker container.
Supposedly, [docker containers have their memory limit set to "unlimited" by default](https://docs.docker.com/config/containers/resource_constraints/#--kernel-memory-details),
so it's natural to assume that running mempig would have the same result there as it
would outside a container.

```
$ docker build -t mempig .
```

For 1 and 2G it was fine:
```
$ docker run -it --rm --name mp mempig -G 1
Allocated 1GiB (1073741824 bytes)
Wallowing in the memory. Press ctrl-C to quit.
^C
$ docker run -it --rm --name mp mempig -G 2
Allocated 2GiB (2147483648 bytes)
Wallowing in the memory. Press ctrl-C to quit.
^C^C^C^C^C^C
```
But for 4G it took a long time allocating the memory and then crashed with no message:
```
$ time docker run -it --rm --name mp mempig -G 4

real    2m12.856s
user    0m0.042s
sys     0m0.019s
```

Does it get better with an explicit setting for --memory?
```
$ time docker run -it --rm --name mp --memory=8G mempig -G 4

real    2m15.208s
user    0m0.040s
sys     0m0.021s
```
Nope.
