# Base container for compile service
FROM golang:alpine AS builder

# Install dependencies
RUN apk add make

# Go to builder workdir
WORKDIR /go/src/github.com/cagodoy/tenpo-restaurants-api/

# Copy go modules files
COPY go.mod .
COPY go.sum .

# Install dependencies
RUN go mod download

# Copy all source code
COPY . .

# Compile service
RUN make linux

#####################################################################
#####################################################################

# Base container for run service
FROM alpine

# Go to workdir
WORKDIR /src/tenpo-restaurants-api

# Install dependencies
RUN apk add --update ca-certificates wget

# Copy binaries
COPY --from=builder /go/src/github.com/cagodoy/tenpo-restaurants-api/bin/tenpo-restaurants-api /usr/bin/tenpo-restaurants-api

# Expose service port
EXPOSE 5030

# Run service
CMD ["/bin/sh", "-l", "-c", "tenpo-restaurants-api"]