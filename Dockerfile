#https://hub.docker.com/r/supinf/go-swagger/dockerfile

#FROM golang:1.12-alpine AS build_base
FROM golang:1.13.1-alpine3.10 AS build 

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

ENV SWAGGER_VERSION=v0.20.1

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

#ENV SwaggerREPO="github.com/go-swagger/go-swagger/cmd/swagger"
RUN go get -u github.com/go-swagger/go-swagger/cmd/swagger

WORKDIR /go/src/github.com/go-swagger/go-swagger


CMD [ "ls" "/go/src/github.com/go-swagger/go-swagger" ]

# RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

# COPY swagger.exe /app

# WORKDIR /app

#RUN go mod download

# # Build the Go app
#RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./out . #/app/out

# # Start fresh from a smaller image
# FROM alpine:3.9



# COPY --from=build_base /app/out /apic.exe

# WORKDIR /

#CMD ["/apic.exe"]