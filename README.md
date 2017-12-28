System Management - Service Deployment Management
=======================================

This provides funtionalities to deploy, update, terminate a container or containers to a certain edge device or a group of edge devices. Also, this provides APIs to create, update, and delete a group of edge devices which container(s) can be deployed at the same time.

## Prerequisites ##
- docker-ce
    - Version: 17.09
    - [How to install](https://docs.docker.com/engine/installation/linux/docker-ce/ubuntu/)
- go compiler
    - Version: 1.8 or above
    - [How to install](https://golang.org/dl/)

## How to build ##
This provides how to build sources codes to an excutable binary and dockerize it to create a Docker image.

#### 1. Executable binary ####
```shell
$ ./build.sh
```
If source codes are successfully built, you can find an output binary file, **main**, on a root of project folder.
Note that, you can find other build scripts, **build_arm.sh** and **build_arm64**, which can be used to build the codes for ARM and ARM64 machines, respectively.

#### 2. Docker Image  ####
Next, you can create it to a Docker image.
```shell
$ docker build -t system-deployment-mgmt-go-ubuntu -f Dockerfile .
```
If it succeeds, you can see the built image as follows:
```shell
$ sudo docker images
REPOSITORY                         TAG        IMAGE ID        CREATED           SIZE
system-deployment-mgmt-go-ubuntu   latest     fcbbd4c401c2    31 seconds ago    157MB
```
Note that, you can find other Dockerfiles, **Dockerfile_arm** and **Dockerfile_arm64**, which can be used to dockerize for ARM and ARM64 machines, respectively.

## How to run with Docker image ##
Required options to run Docker image
- port
    - 48099:48099
- volume
    - "host folder"/data/db:/data/db (Note that you should replace "host folder" to a desired folder on your host machine)

You can execute it with a Docker image as follows:
```shell
$ docker run -it -p 48099:48099 -v /data/db:/data/db system-deployment-mgmt-go-ubuntu
```
If it succeeds, you can see log messages on your screen as follows:
```shell
$ docker run -it -p 48099:48099 -v /data/db:/data/db system-deployment-mgmt-go-ubuntu
2017-12-28T02:40:03.777+0000 I CONTROL  [initandlisten] MongoDB starting : pid=7 port=27017 dbpath=/data/db 64-bit host=e180bb47c8c7
2017-12-28T02:40:03.777+0000 I CONTROL  [initandlisten] db version v3.4.4
2017-12-28T02:40:03.778+0000 I CONTROL  [initandlisten] git version: 888390515874a9debd1b6c5d36559ca86b44babd
2017-12-28T02:40:03.778+0000 I CONTROL  [initandlisten] OpenSSL version: LibreSSL 2.5.5
2017-12-28T02:40:03.778+0000 I CONTROL  [initandlisten] allocator: system
2017-12-28T02:40:03.778+0000 I CONTROL  [initandlisten] modules: none
2017-12-28T02:40:03.778+0000 I CONTROL  [initandlisten] build environment:
2017-12-28T02:40:03.778+0000 I CONTROL  [initandlisten]     distarch: x86_64
2017-12-28T02:40:03.778+0000 I CONTROL  [initandlisten]     target_arch: x86_64
2017-12-28T02:40:03.778+0000 I CONTROL  [initandlisten] options: { storage: { mmapv1: { smallFiles: true } } }
2017-12-28T02:40:03.783+0000 I STORAGE  [initandlisten]
2017-12-28T02:40:03.783+0000 I STORAGE  [initandlisten] ** WARNING: Using the XFS filesystem is strongly recommended with the WiredTiger storage engine
2017-12-28T02:40:03.783+0000 I STORAGE  [initandlisten] **          See http://dochub.mongodb.org/core/prodnotes-filesystem
2017-12-28T02:40:03.784+0000 I STORAGE  [initandlisten] wiredtiger_open config: create,cache_size=3453M,session_max=20000,eviction=(threads_min=4,threads_max=4),config_base=false,statistics=(fast),log=(enabled=true,archive=true,path=journal,compressor=snappy),file_manager=(close_idle_time=100000),checkpoint=(wait=60,log_size=2GB),statistics_log=(wait=0),
2017-12-28T02:40:03.972+0000 W STORAGE  [initandlisten] Detected configuration for non-active storage engine mmapv1 when current storage engine is wiredTiger
2017-12-28T02:40:03.972+0000 I CONTROL  [initandlisten]
2017-12-28T02:40:03.972+0000 I CONTROL  [initandlisten] ** WARNING: Access control is not enabled for the database.
2017-12-28T02:40:03.972+0000 I CONTROL  [initandlisten] **          Read and write access to data and configuration is unrestricted.
2017-12-28T02:40:03.972+0000 I CONTROL  [initandlisten] ** WARNING: You are running this process as the root user, which is not recommended.
2017-12-28T02:40:03.972+0000 I CONTROL  [initandlisten]
2017-12-28T02:40:03.972+0000 I CONTROL  [initandlisten]
2017-12-28T02:40:03.973+0000 I CONTROL  [initandlisten] ** WARNING: /sys/kernel/mm/transparent_hugepage/enabled is 'always'.
2017-12-28T02:40:03.973+0000 I CONTROL  [initandlisten] **        We suggest setting it to 'never'
2017-12-28T02:40:03.973+0000 I CONTROL  [initandlisten]
2017-12-28T02:40:04.021+0000 I FTDC     [initandlisten] Initializing full-time diagnostic data capture with directory '/data/db/diagnostic.data'
2017-12-28T02:40:04.093+0000 I INDEX    [initandlisten] build index on: admin.system.version properties: { v: 2, key: { version: 1 }, name: "incompatible_with_version_32", ns: "admin.system.version" }
2017-12-28T02:40:04.093+0000 I INDEX    [initandlisten] 	 building index using bulk method; build may temporarily use up to 500 megabytes of RAM
2017-12-28T02:40:04.095+0000 I INDEX    [initandlisten] build index done.  scanned 0 total records. 0 secs
2017-12-28T02:40:04.095+0000 I COMMAND  [initandlisten] setting featureCompatibilityVersion to 3.4
2017-12-28T02:40:04.095+0000 I NETWORK  [thread1] waiting for connections on port 27017
```

## (Optional) How to enable QEMU environment on your computer
QEMU could be useful if you want to test your implemetation on various CPU architectures(e.g. ARM, ARM64) but you have only Ubuntu PC. To enable QEMU on your machine, please do as follows.

Required packages for QEMU:
```shell
$ apt-get install -y qemu-user-static binfmt-support
```
For ARM 32bit:
```shell
$ echo ':arm:M::\x7fELF\x01\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00\x28\x00:\xff\xff\xff\xff\xff\xff\xff\x00\xff\xff\xff\xff\xff\xff\xff\xff\xfe\xff\xff\xff:/usr/bin/qemu-arm-static:' > /proc/sys/fs/binfmt_misc/register <br />
```
For ARM 64bit:
```shell
$ echo ':aarch64:M::\x7fELF\x02\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00\xb7:\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xfe\xff\xff:/usr/bin/qemu-aarch64-static:' > /proc/sys/fs/binfmt_misc/register <br />
```

Now, you can build your implementation and dockerize it for ARM and ARM64 on your Ubuntu PC. The below is an example for ARM build.

```shell
$ ./build.sh
$ docker build -t system-deployment-mgmt-go-arm -f Dockerfile_arm .
```
