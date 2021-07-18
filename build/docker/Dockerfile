# Stage 1: Building
FROM golang:1.16-buster AS builder

# Environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Mark the "main" directory used by the builder stage
WORKDIR /build

# Install make which is later used to build
RUN apt-get update && apt-get install make

# Ensure go has the nessecary information to download the modules
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download

# Copy everything not ignored by the .dockerignore file in the root dir
COPY . .

# invokes the make target with the call the go compiler
# which in turns builds the binary executable
RUN make golang_build

# Make workdir to the distribution
WORKDIR /dist

# Copy the go binary executable to the current workdir (./dist)
RUN cp ../build/dist/golang/bin/main ./

# Stage 2: Executing
# scratch is a no-op (Does not add a layer) image used to create
# minimal images which often only contains a single binary.
# We use sratch because the golang compiler can build a single executable.
FROM scratch

# We want to copy the distribution build for the container to the
# this stage because the binary is the one to be executed.
COPY --from=builder ./dist/main /

# Expose the port on which the REST API HTTP server listens to
EXPOSE 8080/tcp

# Execute the binary executable
ENTRYPOINT ["/main"]
