ARG ARCH=amd64
FROM gcr.io/distroless/java17-debian11:latest-${ARCH}

ARG ARCH
COPY ./bin/mcman.linux.${ARCH} /app
COPY ./resources/server-1.20.4.jar /
COPY ./static /static
COPY ./templates /templates
COPY ./docker/config.yaml /default.yaml

# Minecraft Server
EXPOSE 25565
# Web UI
EXPOSE 8000

ENTRYPOINT [ "/app" ]
