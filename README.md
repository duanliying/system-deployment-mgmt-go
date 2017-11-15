Docs @ System Management - Service Deployment Management
=======================================

This provides funtionalities to deploy, update, terminate a container or containers to a certain edge device or a group of edge devices. Also, this provides APIs to create, update, and delete a group of edge devices which container(s) can be deployed at the same time.

### How to build source codes

> ./build.sh

### How to build a Docker image

> docker build --tag "image_name":"tag" -f Dockerfile .

### How to get a Docker image of service deployment management (ONLY for Samsung internal)

> docker pull docker.sec.samsung.net:5000/edge/servicedeployment/servicedeploymentagentmanager/ubuntu_x86_64

### How to run a container of service deployment management

> docker run -it -p 48099:48099 -v "host folder"/data/db:/data/db "image_name":"tag"

> Note that you can also using **"docker-compose"**. <br />
> **docker-compose -f ./docker-compose_ubuntu.yml up**

### How to enable QEMU environment on your computer (i.e. Ubuntu machine)

> apt-get install -y qemu-user-static binfmt-support

> (For ARM 32bit) <br />
> echo ':arm:M::\x7fELF\x01\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00\x28\x00:\xff\xff\xff\xff\xff\xff\xff\x00\xff\xff\xff\xff\xff\xff\xff\xff\xfe\xff\xff\xff:/usr/bin/qemu-arm-static:' > /proc/sys/fs/binfmt_misc/register <br />

> (For ARM 64bit) <br />
> echo ':aarch64:M::\x7fELF\x02\x01\x01\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02\x00\xb7:\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xfe\xff\xff:/usr/bin/qemu-aarch64-static:' > sudo /proc/sys/fs/binfmt_misc/register <br />

## Reference

##### Golang install
> https://github.com/golang/go/wiki/Ubuntu

##### Dockerfile
> for ubuntu_x86_64 : Dockerfile <br />
> for raspberry pi3 : Dockerfile_RPI3 <br />

