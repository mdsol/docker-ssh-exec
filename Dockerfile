FROM busybox:latest

ADD pkg/docker-ssh-exec /docker-ssh-exec

ENTRYPOINT ["/docker-ssh-exec"]
