#
# to build image locally run the following from root of bizzle:
# docker build -f build/package/user/Dockerfile -t gcr.io/$(gcloud config get-value project)/user .
#
# to run locally
# docker run -p 8080:8080 gcr.io/$(gcloud config get-value project)/user
#
# to push:
# docker push gcr.io/$(gcloud config get-value project)/user
#
# ensure that docker is is authorised to push to gcloud. set up with:
# gcloud auth configure-docker

# Base build image
FROM golang:latest AS build_base

# set working directory inside the container
WORKDIR /user

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
RUN go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o user cmd/user/main.go

# this last stage produces the final lean production image
# start from a fresh Alpine image to reduce the image size
# (i.e. not ship the Go compiler in our production artifacts)
FROM alpine
RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /user/user /user

# Copy config file to the production image from the builder stage
COPY --from=builder /user/configs/user/config.toml /config.toml

# Run the web service on container startup
CMD ["/user"]