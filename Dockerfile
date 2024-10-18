# syntax=docker/dockerfile:1
FROM golang as build

WORKDIR /src

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Required for sqlite to work
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/nubayrah ./cmd/api/main.go

# Must use image that contains dynamic linking in order to copy files to /bin
FROM debian:bookworm

EXPOSE 5050
COPY --from=build /bin/nubayrah /
# Run
CMD ["/nubayrah"]