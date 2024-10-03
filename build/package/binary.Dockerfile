FROM alpine:3.20.1

RUN apk add --no-cache bash tzdata dumb-init

WORKDIR /opt/app

# Copy files to docker image
COPY configs/.env configs/.env
COPY backend .

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]

CMD [ "./backend" ]