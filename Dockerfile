FROM golang:1.10

# Copy the local package files to the container's workspace
ADD . /go/src/github.com/torlenor/AbyleBotter

# Fetch dependencies and build AbyleBotter inside the container
RUN go get -v ./...
RUN go install github.com/torlenor/AbyleBotter

# Run the command by default when the container starts
ENTRYPOINT /go/bin/AbyleBotter

# Document that the service listens on port 8080
# Not necessary, yet
# EXPOSE 8080