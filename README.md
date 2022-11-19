[![Lint and Test](https://github.com/doppiolab/mcman/actions/workflows/lint-and-test.yml/badge.svg)](https://github.com/doppiolab/mcman/actions/workflows/lint-and-test.yml)
[![Create and publish a Docker image](https://github.com/doppiolab/mcman/actions/workflows/publish-docker.yml/badge.svg)](https://github.com/doppiolab/mcman/actions/workflows/publish-docker.yml)
[![codecov](https://codecov.io/gh/doppiolab/mcman/branch/main/graph/badge.svg?token=BAYOLNV6XI)](https://codecov.io/gh/doppiolab/mcman)

# mcman (**M**ine**C**raft server **MAN**ager)

<p align='center'>
    <img src='https://raw.githubusercontent.com/doppiolab/mcman/main/static/favicon.ico'>
</p>
<p align='center'>
    Yet another Minecraft Server Manager for players.
</p>

## Introduction

This mcman is a small toy project to provide the Web UI to manage a dedicated server for Minecraft Java Edition.
It provides the features as follows:

* Shows near-realtime top-view map in the Web UI with player positions.
* Sends server command and prints the logs from Minecraft server in the Web UI.
* Sends Minecraft server logs via webhook to Discord and Slack.

Enjoy!

## Usage

You can use this project with Docker.
If you want to use it without Docker, you can build it for yourself.

### With Docker

You can simply run mcman server packed with Minecraft dedicated server as follows.

```sh
docker run \
    -v `pwd`/server-data:/mcman/data \
    -v `pwd`/config.yaml:/config.yaml \
    -p SOME_PORT:8000 \
    ghcr.io/doppiolab/mcman:{IMAGE_TAG}-amd64 \
    -config /config.yaml
```

* check image tags in [GitHub Packages](https://github.com/doppiolab/mcman/pkgs/container/mcman)
* `arm64` and `amd64` are supported.
* If you want to persist your server data, try to mount `/mcman/data` path. (or the value of `minecraft.working_dir` option)

### Configuration file

It is recommended to write your own configuration file from [`./docker/config.yaml`](./docker/config.yaml).

```yaml
# Server configuration
server:
    # Host to listen. default value is ":8000"
    host: 0.0.0.0:8000
    # Adjust log level to debug if set true.
    debug: false
    # path to static directory.
    # recommend not to change the value in docker env.
    static_path: /static
    # path to template directory./
    # recommend not to change the value in docker env.
    template_path: /templates
    # temp path.
    # this is used for caching the map images.
    temporary_path: /tmp
    # Environment variable key for initial password.
    # If password_env_key is not set or no env var named password_env_key, mcman generate random password and print it in the server log.
    password_env_key:

# Minecraft dedicated server configuration
minecraft:
    # java command or java binary path
    java_command: /usr/bin/java
    # minecraft server jar path.
    # It is recommended to use this path as absolute path,
    # since this is relative to working_dir value.
    jar_path: "/server-1.19.2.jar"
    # Working directory for minecraft server
    # Miencraft server will save the all server data into this directory.
    working_dir: "/mcman/data"
    # Java option.
    # If you want to add options like -Xms1024M, use this option.
    java_options:
        # - option1
        # - option2
    # Minecraft jar args.
    args:
        - nogui
    # Config used for development.
    # Do not use this for real usage.
    skip_start_for_debug:
    # Log webhook option.
    # If you want to stream the minecraft server log to your messenger, use this option.
    log_webhook:
        # Slack incoming webhook url
        slack: ""
        # Discord webhook url
        discord: ""
        # mcman will debounce the server log to avoid rate limiting.
        # mcman sends logs if no new logs occur during debounce_threshold.
        # unit: millisecond
        debounce_threshold: 100
```

## Contributes

Check [CONTRIBUTING.md](./CONTRIBUTING.md)
