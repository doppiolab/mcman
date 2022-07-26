# Contributing

## Issue reporting

If you have any issues with this, [you can submit the issue first](https://github.com/doppiolab/mcman/issues). The more detailed and clear description is provided, the faster it will be resolved. We expect as follows.

* OS information
* mcman version or git hash that you used
* Minecraft Java Edition version
* Server logs

## Development

### Make

We use `Makefile` for the project commands.

```sh
$ make help
Usage:

 build-linux   build the project for linux
 lint          lint the source code
 test          run all tests
 cover         run all tests with checking code coverages
 gen           run go generate
 help          prints this help message
$ make lint # run lint
...
$ make test # run test
...
```

### Build

You can use this project without Docker, but using Docker is strongly suggested.

To support building multi-CPU architecture, we build project binaries locally and then copy them to Docker.
If you want to build Docker image, run the commands as follows.

```sh
CPU_ARCH=amd64 # change this arm64 if you want
make build-linux # this will produce bin/mcman.linux.(arm64|amd64)
docker build \
    --platform ${CPU_ARCH} \
    --build-arg ARCH=${CPU_ARCH} \
    --file docker/Dockerfile \s
    --tag ${IMAGE_NAME}:${IMAGE_TAG}-${CPU_ARCH} \
    .
```
