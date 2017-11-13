Docs @ System Management - Service Deployment Management
=======================================

This provides funtionalities to deploy, update, terminate a container or containers to a certain edge device or a group of edge devices. Also, this provides APIs to create, update, and delete a group of edge devices which container(s) can be deployed at the same time.

### How to build source codes

> ./build.sh

### How to build Docker image

> docker build --tag "image_name":"tag" -f Dockerfile .

### How to get an Docker image of service deployment management (ONLY for Samsung internal)

> dockr pull docker.sec.samsung.net:5000/edge/servicedeployment/servicedeploymentagentmanager/ubuntu_x86_64

### How to run a container of service deployment management

> docker run -it -p 48099:48099 -v "host folder"/data/db:/data/db "image_name":"tag"

> Note that you can also using **"docker-compose"**. <br />
> **docker-compose -f ./docker-compose_ubuntu.yml up**

## Reference

##### Golang install
> https://github.com/golang/go/wiki/Ubuntu

##### Dockerfile
> for ubuntu_x86_64 : Dockerfile <br />
> for raspberry pi3 : Dockerfile_RPI3 <br />

## Copyright and license

 Code and documentation copyright Team Sharerience.

## Contact

 Jihun Ha <jihun.ha@samsung.com>
