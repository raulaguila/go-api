FROM alpine:3.21.3

RUN apk add --no-cache bash tzdata dumb-init

WORKDIR /opt/api

# Copy files to docker image
COPY configs/.env configs/.env
COPY binbackend .

RUN chmod +x binbackend

ENTRYPOINT [ "/usr/bin/dumb-init", "--" ]

CMD [ "./binbackend" ]