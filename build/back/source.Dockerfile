FROM golang:1.24.0 AS builder

WORKDIR /opt/app

# Copy files to docker image
COPY cmd/ cmd/
COPY configs/ configs/
COPY docs/ docs/
COPY internal/ internal/
COPY pkg/ pkg/
COPY go.mod .
COPY go.sum .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o backend cmd/go-api_test/go-api_test.go


FROM alpine:3.21.3

RUN apk add --no-cache bash tzdata dumb-init

WORKDIR /opt/app

# Copy files to docker image
COPY --from=builder /opt/app/configs/.env configs/.env
COPY --from=builder /opt/app/backend .

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]

CMD [ "./backend" ]