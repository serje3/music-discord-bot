# Use an official Golang runtime as the base image
FROM golang:1.16-alpine

# Set the working directory in the container to /app
WORKDIR /app


# Copy the local package files to the container's workspace
COPY ./src/ .
COPY ./config.cfg .

ENV CGO_ENABLED=0

RUN apk add pkgconfig opus-dev opusfile-dev


RUN go clean -modcache

RUN go mod download

# Build the Go app
RUN go build -o captain-cocker ./

# Run the compiled binary program
CMD ["./captain-cocker"]