# Two-stage build:
#    first  FROM prepares a binary file in full environment ~780MB
#    second FROM takes only binary file ~10MB

FROM golang:1.12 AS builder

RUN go version

COPY . .
WORKDIR "/api"

RUN make
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main


####### second stage to obtain a very small image
FROM scratch

COPY --from=builder /api .

ENV DATABASE_CONNECTION 'mongodb://127.0.0.1:27017/?gssapiServiceName=mongodb'
ENV DATABASE_NAME 'sas'

EXPOSE 8080

CMD ["/main"]

#
#FROM golang:1.11.1-alpine3.8 as build-env
## All these steps will be cached
#RUN mkdir /hello
#WORKDIR /hello
#COPY go.mod . # <- COPY go.mod and go.sum files to the workspace
#COPY go.sum .
#
## Get dependancies - will also be cached if we won't change mod/sum
#RUN go mod download
## COPY the source code as the last step
#COPY . .
#
## Build the binary
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/hello
#FROM scratch # <- Second step to build minimal image
#COPY --from=build-env /go/bin/hello /go/bin/hello
#ENTRYPOINT ["/go/bin/hello"]

