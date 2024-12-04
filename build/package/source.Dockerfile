FROM golang:1.23.3 AS builder

WORKDIR /app

# Copy files to docker image
COPY cmd/ cmd/
COPY configs/ configs/
COPY docs/ docs/
COPY internal/ internal/
COPY pkg/ pkg/
COPY go.mod .
COPY go.sum .
COPY Makefile .

RUN make go-build


FROM alpine:3.20.3

RUN apk add --no-cache bash tzdata dumb-init

WORKDIR /opt/app

# Copy files to docker image
COPY --from=builder /app/configs/.env configs/.env
COPY --from=builder /app/backend .

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]

CMD [ "./backend" ]