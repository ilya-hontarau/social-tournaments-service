FROM golang:1.12.6 as build
ADD . /go/src/github.com/illfate/social-tournaments-service/
WORKDIR /go/src/github.com/illfate/social-tournaments-service/
RUN make dep
RUN make

FROM scratch
COPY --from=build /go/src/github.com/illfate/social-tournaments-service/bin/sts .
ENV PORT 8080
CMD ["./sts"]
