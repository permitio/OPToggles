FROM golang:1.16 AS build

WORKDIR /src
COPY src/go.mod ./
COPY src/go.sum ./
RUN go mod download

COPY src/ ./
RUN go build -o /optoggles


FROM debian:10

WORKDIR /
COPY --from=build /optoggles ./optoggles
RUN groupadd -r nonroot && useradd -r -g nonroot nonroot
USER nonroot:nonroot

ENTRYPOINT ["/optoggles"]
