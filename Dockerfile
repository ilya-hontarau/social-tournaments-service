FROM golang:1.12.6 as build
ADD . /social-tournaments-service/
WORKDIR /social-tournaments-service/
RUN make

FROM scratch
COPY --from=build /social-tournaments-service/bin/sts .
COPY tournament.graphql .
COPY user.graphql .
ENV PORT 8080
CMD ["./sts"]
