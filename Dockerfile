FROM golang:1.13 AS builder

WORKDIR /src

RUN git clone  https://github.com/synyx/imagescaler

WORKDIR /src/imagescaler
RUN GOPROXY=https://proxy.golang.org CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o /imagescaler .

FROM scratch
COPY --from=builder /imagescaler ./
ENTRYPOINT ["./imagescaler"]
