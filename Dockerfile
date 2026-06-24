FROM golang:1.24-alpine3.23 AS builder
RUN apk update && apk upgrade
RUN apk add git
COPY go.* /circ/
COPY ./vendor /circ/vendor
COPY *.go /circ/
RUN cd /circ && CGO_ENABLED=0 go build -mod=vendor -trimpath -ldflags="-s -w"

FROM gcr.io/distroless/static-debian13:nonroot
WORKDIR /circ
COPY --from=builder /circ/circ /circ/circ
USER 65532:65532
ENTRYPOINT ["/circ/circ"]
