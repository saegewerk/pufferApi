FROM golang:1.14.2


# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
#RUN go get /go/src/github.com/saegewerk/QGTodo/...
ADD . /go/src/github.com/saegewerk/pufferApi
RUN cd /go/src/github.com/saegewerk/pufferApi && go mod download && go get /go/src/github.com/saegewerk/pufferApi/...


# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/pufferServer /go/src/github.com/saegewerk/pufferApi/PufferfileDocker.yml

# Document that the service listens on port 8080.
EXPOSE 8000