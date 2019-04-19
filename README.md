# mempig
Mempig's only job is to eat up your system's memory so you can find answers to questions like these:
* How much memory can mempig ask for before it gets killed?
* Does the system slow to a crawl or lock up completely for large memory requests?
* If mempig is running in a container, how much memory can mempig eat before the container crashes?

## Mac

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

I'm running Docker Desktop, and it has a little icon at the top right of my screen. Clicking that and then clicking the Preferences menu item that appears, a window came up to configure it. I clicked the Advanced tab and a slider appeared for the memory. Sure enough it was set at 2GiB. I moved it to 8GiB, leaving the Swap slider at 1GiB, and clicked Apply & Restart. Did it help?
```
[ ~/src/github.com/ijt/mempig ] time docker run -it --rm --name mp mempig -G 7
Allocated 7GiB (7516192768 bytes)
Wallowing in the memory. Press ctrl-C to quit.
^C
```
Yes it did.

What happens for 8G? Mempig slows to a crawl after allocating 93.49% of the 8G it is supposed to allocate. I ran out of patience waiting for it. It also takes a long for the docker run command to respond to ctrl-C at this point.
```
[ ~/src/github.com/ijt/mempig ] time docker run -it --rm --name mp mempig -G 8
^C^Ccated 8035699713 of 8589934592 bytes (93.55%)
real    5m52.719s
user    0m4.759s
sys     0m9.610s
```

## Ubuntu

On an Ubuntu instance on GCE with about 3G of memory, this happens:

```
issactrotts@memtest:~$ cat /proc/meminfo | head
MemTotal:        3792984 kB
MemFree:         3360260 kB
MemAvailable:    3466828 kB
Buffers:           47144 kB
Cached:           238012 kB
SwapCached:            0 kB
Active:           187120 kB
Inactive:         143568 kB
Active(anon):      45732 kB
Inactive(anon):     6320 kB
issactrotts@memtest:~$ mempig -G 3
Allocated 3GiB (3221225472 bytes)
Wallowing in the memory. Press ctrl-C to quit.
^C
```
So allocating 3GiB is fine.
However, asking for 4G makes mempig crash instead of swapping out on this system:
```
issactrotts@memtest:~$ mempig -G 4
fatal error: runtime: out of memory

runtime stack:
runtime.throw(0x4b53f9, 0x16)
        /usr/lib/go-1.7/src/runtime/panic.go:566 +0x95
runtime.sysMap(0xc420100000, 0x100000000, 0x7f0a44b68e00, 0x5291b8)
        /usr/lib/go-1.7/src/runtime/mem_linux.go:219 +0x1d0
runtime.(*mheap).sysAlloc(0x5109a0, 0x100000000, 0x7f0a00000001)
        /usr/lib/go-1.7/src/runtime/malloc.go:407 +0x37a
runtime.(*mheap).grow(0x5109a0, 0x80000, 0x0)
        /usr/lib/go-1.7/src/runtime/mheap.go:726 +0x62
runtime.(*mheap).allocSpanLocked(0x5109a0, 0x80000, 0x2000)
        /usr/lib/go-1.7/src/runtime/mheap.go:630 +0x4f2
runtime.(*mheap).alloc_m(0x5109a0, 0x80000, 0x7f0100000000, 0x40c679)
        /usr/lib/go-1.7/src/runtime/mheap.go:515 +0xe0
runtime.(*mheap).alloc.func1()
        /usr/lib/go-1.7/src/runtime/mheap.go:579 +0x4b
runtime.systemstack(0x7ffdb13c86a0)
        /usr/lib/go-1.7/src/runtime/asm_amd64.s:314 +0xab
runtime.(*mheap).alloc(0x5109a0, 0x80000, 0x10100000000, 0x7f0a44b08078)
        /usr/lib/go-1.7/src/runtime/mheap.go:580 +0x73
runtime.largeAlloc(0x100000000, 0x529201, 0x7f0a44b08078)
        /usr/lib/go-1.7/src/runtime/malloc.go:774 +0x93
runtime.mallocgc.func1()
        /usr/lib/go-1.7/src/runtime/malloc.go:669 +0x3e
runtime.systemstack(0x50db00)
        /usr/lib/go-1.7/src/runtime/asm_amd64.s:298 +0x79
runtime.mstart()
        /usr/lib/go-1.7/src/runtime/proc.go:1079

goroutine 1 [running]:
runtime.systemstack_switch()
        /usr/lib/go-1.7/src/runtime/asm_amd64.s:252 fp=0xc420037db0 sp=0xc420037da8
runtime.mallocgc(0x100000000, 0x497280, 0x456601, 0xc4200160c0)
        /usr/lib/go-1.7/src/runtime/malloc.go:670 +0x903 fp=0xc420037e50 sp=0xc420037db0
runtime.makeslice(0x497280, 0x100000000, 0x100000000, 0x4b6e58, 0x1f, 0xc420012238)
        /usr/lib/go-1.7/src/runtime/slice.go:57 +0x7b fp=0xc420037ea8 sp=0xc420037e50
main.main()
        /home/issactrotts/src/github.com/ijt/mempig/main.go:12 +0x9b fp=0xc420037f48 sp=0xc420037ea8
runtime.main()
        /usr/lib/go-1.7/src/runtime/proc.go:183 +0x1f4 fp=0xc420037fa0 sp=0xc420037f48
runtime.goexit()
        /usr/lib/go-1.7/src/runtime/asm_amd64.s:2086 +0x1 fp=0xc420037fa8 sp=0xc420037fa0
```

Now let's see what happens in a docker container on Ubuntu, using a GCE image that supports docker.
```
$ gcloud compute instances create memtest2 --image=cos-73-11647-121-0 --image-project=cos-cloud --zone=us-east1-d --machine-type n1-standard-1
rCreated [https://www.googleapis.com/compute/v1/projects/sourcegraph-server/zones/us-east1-d/instances/memtest2].
NAME      ZONE        MACHINE\_TYPE   PREEMPTIBLE  INTERNAL\_IP  EXTERNAL\_IP   STATUS
memtest2  us-east1-d  n1-standard-1               10.142.0.2   34.74.75.205  RUNNING

$ gcloud compute ssh memtest2
No zone specified. Using zone [us-east1-d] for instance: [memtest2].
Warning: Permanently added 'compute.3342904289974829875' (ED25519) to the list of known hosts.

issactrotts@memtest2 ~ $ git clone https://github.com/ijt/mempig
Cloning into 'mempig'...
remote: Enumerating objects: 41, done.
remote: Counting objects: 100% (41/41), done.
remote: Compressing objects: 100% (33/33), done.
remote: Total 41 (delta 13), reused 33 (delta 7), pack-reused 0
Unpacking objects: 100% (41/41), done.
issactrotts@memtest2 ~ $ cd mempig/

issactrotts@memtest2 ~/mempig $ docker build -t mempig .
(yadda yadda)
```
If we ask for 3G, it works fine:
```
issactrotts@memtest2 ~/mempig $ docker run -it --rm --name mp mempig -G 3
Allocated 3GiB (3221225472 bytes)
Wallowing in the memory. Press ctrl-C to quit.
^C
```
However, if we ask for 4G, more than is available on the system, the container gets killed after a few seconds:
```
issactrotts@memtest2 ~/mempig $ time docker run -it --rm --name mp mempig -G 4

real    0m7.183s
user    0m0.025s
sys     0m0.033s
```

