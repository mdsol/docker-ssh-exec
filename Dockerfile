FROM golang:1.9 as builder
ADD . /go/src/docker-ssh-exec
ENV CGO_ENABLED=0 GOOS=linux
WORKDIR /go/src/docker-ssh-exec
RUN go build -ldflags '-w -s' -a -installsuffix cgo -o /docker-ssh-exec

FROM scratch as runtime
COPY --from=builder /docker-ssh-exec /docker-ssh-exec
EXPOSE 80
ENTRYPOINT ["/docker-ssh-exec"]
