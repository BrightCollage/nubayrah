# syntax=docker/dockerfile:1
FROM golang:1.23-alpine as build

WORKDIR /src

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/nubayrah ./cmd/api/main.go

FROM scratch

EXPOSE 5050
COPY --from=build /bin/nubayrah /
# Run
CMD ["/nubayrah"]