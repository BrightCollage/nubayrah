# syntax=docker/dockerfile:1

# Builds the react project
FROM node:18-alpine as build-client

WORKDIR /app
COPY /client .
# Add -- --mode docker to ensure that default build outDir is used.
RUN npm install && npm run build -- --mode docker


# Builds the go modules as a binary
FROM golang AS build-api

WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/nubayrah --tags=docker -a -ldflags '-linkmode external -extldflags "-static"' ./cmd/nubayrah/main.go


# Final image to host application
FROM alpine

COPY --from=build-api /bin/nubayrah /nubayrah
COPY --from=build-client /app/static /static
# Create the base data and library directories.
RUN mkdir /data /library


EXPOSE 8090 5050


# Run
CMD ["/nubayrah"]