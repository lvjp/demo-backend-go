FROM golang:1.24.5@sha256:ef5b4be1f94b36c90385abd9b6b4f201723ae28e71acacb76d00687333c17282 AS build

WORKDIR /go/src/app

RUN go env
RUN [ $(go env GOCACHE) = /root/.cache/go-build ]
RUN [ $(go env GOMODCACHE) = /go/pkg/mod ]

RUN --mount=type=bind,target=.,readonly \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go mod download -x && CGO_ENABLED=0 go build -o /go/bin/app .

FROM gcr.io/distroless/static-debian12:latest@sha256:b7b9a6953e7bed6baaf37329331051d7bdc1b99c885f6dbeb72d75b1baad54f9

COPY --from=build /go/bin/app /
CMD ["/app", "server"]
