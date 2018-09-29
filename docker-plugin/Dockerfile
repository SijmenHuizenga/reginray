#FROM golang:1.10 as builder
#
#RUN mkdir /src
#ADD . /src
#WORKDIR /src
#RUN go get -d -v
#RUN CGO_ENABLED=0 GOOS=linux go build -a -o /go/bin/aspicio
#
#
#FROM alpine:3.8
#RUN apk --no-cache add ca-certificates
#COPY --from=builder /go/bin/aspicio /usr/bin/aspicio-log-driver
#ENTRYPOINT ["/usr/bin/aspicio-log-driver"]

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
ADD aspicio /usr/bin/aspicio-log-driver
ENTRYPOINT ["/usr/bin/aspicio-log-driver"]