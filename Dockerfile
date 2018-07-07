FROM golang:1.10

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/torlenor/AbyleBotter

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
run go get -v ./...
RUN go install github.com/torlenor/AbyleBotter

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/AbyleBotter
# CMD ["/go/bin/AbyleBotter"]

# Document that the service listens on port 8080.
# EXPOSE 8080