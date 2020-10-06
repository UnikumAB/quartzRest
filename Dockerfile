FROM golang:alpine AS builder

RUN apk --no-cache add make gcc libc-dev

WORKDIR /build

COPY go.* /build/

RUN go mod download

COPY Makefile .
COPY cmd /build/cmd
COPY pkg /build/pkg
COPY .golangci.yml /build/

RUN make build

FROM alpine AS runner

COPY --from=builder /build/quartzRestServer /
CMD ["/quartzRestServer"]

# Add Tini
RUN apk --no-cache add tini
ENTRYPOINT ["/sbin/tini", "--"]
