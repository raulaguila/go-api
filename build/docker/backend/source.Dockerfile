FROM golang:1.24.4 AS builder

WORKDIR /opt/api

COPY go.mod .
COPY go.sum .

RUN go mod download

# Copy files to docker image
COPY cmd/ cmd/
COPY configs/ configs/
COPY docs/ docs/
COPY internal/ internal/
COPY pkg/ pkg/
COPY Makefile Makefile

RUN make build


FROM alpine:3.21.3

RUN apk add --no-cache bash tzdata dumb-init

WORKDIR /opt/api

# Copy files to docker image
COPY --from=builder /opt/api/configs/.env configs/.env
COPY --from=builder /opt/api/binbackend .

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]

CMD [ "./binbackend" ]