FROM golang:1.19.0-alpine AS build
ARG TARGETOS
ARG TARGETARCH
WORKDIR /src
COPY . /src/
RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o app ./api

FROM alpine
WORKDIR /api
COPY --from=build /src /api
ENTRYPOINT ["./app"]