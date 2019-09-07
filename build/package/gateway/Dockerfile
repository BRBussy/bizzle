# Base build image
FROM golang:latest AS build_base

# set working directory inside the container
WORKDIR /gateway

# copy in go dependency management files
COPY go.mod go.sum ./

# populate the module cache based on the go.{mod,sum} files
RUN go mod download

# the binary is built in this image
FROM build_base AS builder

# copy the source from the current directory to the container's working directory
COPY . .

# build the Go app binary
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o gateway cmd/gateway/main.go

# this last stage produces the final lean production image
# start from a fresh Alpine image to reduce the image size
# (i.e. not ship the Go compiler in our production artifacts)
FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /gateway/gateway /gateway

# Run the web service on container startup
CMD ["/gateway"]