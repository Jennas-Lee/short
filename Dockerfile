FROM public.ecr.aws/docker/library/golang:1.24.1-alpine3.21 AS builder

WORKDIR /tmp/tiny-golang-image

COPY . ./

RUN go build ./main.go

FROM busybox:1.37.0-glibc

COPY --from=builder /tmp/tiny-golang-image/main /main

ENTRYPOINT ["/main"]