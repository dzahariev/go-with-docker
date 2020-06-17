FROM golang:alpine AS build
WORKDIR /go/src/github.com/dzahariev/go-with-docker/
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "$(build/ldflags)" -o /app main.go

FROM build AS test-prepare
RUN CGO_ENABLED=0 go test -i ./... 

FROM test-prepare AS test-main
RUN CGO_ENABLED=0 go test . -ginkgo.noColor -v
RUN touch finished.test-main

FROM test-prepare AS test-pckg01
RUN CGO_ENABLED=0 go test ./pkg01 -ginkgo.noColor -v
RUN touch finished.test-pckg01

FROM test-prepare AS test-pckg02
RUN CGO_ENABLED=0 go test ./pkg02 -ginkgo.noColor -v
RUN touch finished.test-pckg02

FROM scratch AS join
COPY --from=test-main /go/src/github.com/dzahariev/go-with-docker/finished.test-main /test-results/.
COPY --from=test-pckg01 /go/src/github.com/dzahariev/go-with-docker/finished.test-pckg01 /test-results/.
COPY --from=test-pckg02 /go/src/github.com/dzahariev/go-with-docker/finished.test-pckg02 /test-results/.
COPY --from=build /app /app

FROM scratch AS release
WORKDIR /bin
COPY --from=join /app /bin/.
ENTRYPOINT [ "./bin/app" ]
