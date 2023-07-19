# [Container image good practices](https://sysdig.com/blog/dockerfile-best-practices/)

# Image reduction tutorial source
https://marcolenzo.eu/reduce-the-size-of-docker-images-with-few-steps/


## Build sample containers

CD docker_image_reduction_build_files/
```
cd sample-1-golang
docker build -t sample-original -f Dockerfile.original .
docker build -t sample-debian -f Dockerfile.multistage.debian .
docker build -t sample-alpine -f Dockerfile.multistage.alpine .
docker build -t sample-distroless -f Dockerfile.multistage.distroless .
```
check sizes 
```
docker image ls | grep sample
```
Build origional and optimized images
```
cd sample-3-layers
docker build -t layers-original -f Dockerfile.original .
docker build -t layers-optimized -f Dockerfile.optimized .
docker image ls | grep layers
```
check sizes

## Run container on port 8090
```
docker run -d -p8090:8090 --name c1 sample-distroless
```
view response
```
curl localhost:8090
```
try to run shell or other commands
```
docker exec -ti c1 sh
docker exec -ti c1 bash
```
# Use dive tool to analyse images
https://github.com/wagoodman/dive

Install via brew
or winget on windows
```
dive layers-original
```
use arrows to  review layers
```
dive layers-optimized
```
Check the efficency difference