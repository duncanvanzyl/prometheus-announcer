FROM golang as builder
RUN mkdir -p /build/
WORKDIR /build/
ADD . /build/
RUN export VERSION=$( cat ./VERSION ) \
  BUILDDATE=$( date "+%Y-%m-%dT%H:%M:%S") \
  GITCOMMIT=$( git rev-parse --short HEAD ) \
  CGO_ENABLED=0 GOOS=linux \
  && go build -tags timetzdata -a -installsuffix cgo -ldflags "-s -w -extldflags '-static' -X 'main.version=$VERSION' -X 'main.buildTime=$BUILDDATE' -X 'main.gitCommit=$GITCOMMIT'" -o /build/main ./cmd/server/


FROM scratch
ENV LOG_LEVEL info
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["/app/main"]
