FROM golang:1.16 AS build

WORKDIR /src
COPY src/go.mod ./
COPY src/go.sum ./
RUN go mod download

COPY src/ ./
RUN go build -o /optoggles


FROM debian:10

WORKDIR /
RUN apt-get update && apt-get install -y ca-certificates && apt-get clean
RUN groupadd -r nonroot && useradd -r -g nonroot nonroot
COPY --from=build /optoggles ./optoggles
USER nonroot:nonroot

ENTRYPOINT ["/optoggles"]
