FROM golang:1.12-alpine AS builder

RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/password && \
    echo 'nobody:x:65534:' > /user/group
RUN apk add 

WORKDIR /src

COPY go.mod .
COPY go.sum .
COPY . .

# TODO: Consider how this docker image can be build for non-x86 targets
RUN CGO_ENABLED=0 GOOS=linux GOARCH=linux GOFLAGS=-mod=vendor go build -a -ldflags '-s' -installsuffix cgo -o /cheapskate

#############

FROM scratch AS final

COPY --from=builder /user/group /user/password /etc
COPY --from=builder /etc/ssl/certs/ca-certificate.crt /etc/ssl/certs
COPY --from=builder /cheapskate /usr/local/bin/cheapskate

USER nobody:nobody
EXPOSE 8080

ENTRYPOINT [ "/usr/local/bin/cheapskate" ]
