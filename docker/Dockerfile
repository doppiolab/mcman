ARG GOLANG_VERSION=1.18
FROM golang:${GOLANG_VERSION} as build

WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 go build -o /go/bin/app

FROM gcr.io/distroless/java17-debian11
COPY --from=build /go/bin/app /
COPY ./resources/server-1.19.jar /
COPY ./static /static
COPY ./templates /templates
COPY ./docker/config.yaml /default.yaml

# Minecraft Server
EXPOSE 25565
# Web UI
EXPOSE 8000

ENTRYPOINT [ "/app" ]
