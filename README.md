# go-with-docker
Sample repo for CI with docker using multistage builds. 
The Dockerfile example shows how to execute unit tests for different packages in parallel with the [BuildKit](https://docs.docker.com/develop/develop-images/build_enhancements/) engine.

## Preconditions

Installations:

- [docker](https://docs.docker.com/get-docker/)

## Details

Execution steps are:
 - Stage: `build` - execute build to ensure that binary can be prodiced from sources
 - Stage: `test-prepare` - prepares the tests environment by fetching all required go libraries in separate layer
 - Stage: `test-main`, `test-pckg01`, `test-pckg02` - execute the tests for each one of the packages
 - Stage: `join` - joins the build results by copying a dummy files that were created as result of each test stage
 - Stage: `release` - prepare image from scratch that contain only binary file that was produced from `build` stage.

To start the tests with BuildKit enabled:
```
DOCKER_BUILDKIT=1 docker build . --no-cache
```
Result shows the steps:
```
[+] Building 17.0s (22/22) FINISHED                                                                                 
 => [internal] load .dockerignore                                                                              0.0s
 => => transferring context: 2B                                                                                0.0s
 => [internal] load build definition from Dockerfile                                                           0.0s
 => => transferring dockerfile: 37B                                                                            0.0s
 => [internal] load metadata for docker.io/library/golang:alpine                                               0.0s
 => [internal] load build context                                                                              0.0s
 => => transferring context: 18.79kB                                                                           0.0s
 => CACHED [release 1/2] WORKDIR /bin                                                                          0.0s
 => [build 1/4] FROM docker.io/library/golang:alpine                                                           0.0s
 => CACHED [build 2/4] WORKDIR /go/src/github.com/dzahariev/go-with-prow/                                      0.0s
 => [build 3/4] COPY . ./                                                                                      0.0s
 => [build 4/4] RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "$(build/ld  2.8s
 => [test-prepare 1/1] RUN CGO_ENABLED=0 go test -i ./...                                                      5.8s
 => [test-main 1/2] RUN CGO_ENABLED=0 go test . -ginkgo.noColor -v                                             7.3s
 => [test-pckg02 1/2] RUN CGO_ENABLED=0 go test ./pkg02 -ginkgo.noColor -v                                     7.6s
 => [test-pckg01 1/2] RUN CGO_ENABLED=0 go test ./pkg01 -ginkgo.noColor -v                                     7.6s
 => [test-main 2/2] RUN touch finished.test-main                                                               0.3s 
 => [test-pckg02 2/2] RUN touch finished.test-pckg02                                                           0.4s 
 => [join 1/4] COPY --from=test-main /go/src/github.com/dzahariev/go-with-prow/finished.test-main /test-resul  0.0s 
 => [test-pckg01 2/2] RUN touch finished.test-pckg01                                                           0.4s
 => [join 2/4] COPY --from=test-pckg01 /go/src/github.com/dzahariev/go-with-prow/finished.test-pckg01 /test-r  0.0s 
 => [join 3/4] COPY --from=test-pckg02 /go/src/github.com/dzahariev/go-with-prow/finished.test-pckg02 /test-r  0.1s 
 => [join 4/4] COPY --from=build /app /app                                                                     0.0s 
 => [release 2/2] COPY --from=join /app /bin/.                                                                 0.0s
 => exporting to image                                                                                         0.0s
 => => exporting layers                                                                                        0.0s
 => => writing image sha256:6e3cdfbce50ac4dfd3da27e6713a1db951ac30589bc48041014fc47795fef6f6                   0.0s
```

In case tets are failing, the release image is not prodiced and build result look like this:
```
[+] Building 16.6s (14/21)                                                                                          
 => [internal] load build definition from Dockerfile                                                           0.0s
 => => transferring dockerfile: 37B                                                                            0.0s
 => [internal] load .dockerignore                                                                              0.0s
 => => transferring context: 2B                                                                                0.0s
 => [internal] load metadata for docker.io/library/golang:alpine                                               0.0s
 => [internal] load build context                                                                              0.0s
 => => transferring context: 17.90kB                                                                           0.0s
 => CACHED [release 1/2] WORKDIR /bin                                                                          0.0s
 => [build 1/4] FROM docker.io/library/golang:alpine                                                           0.0s
 => CACHED [build 2/4] WORKDIR /go/src/github.com/dzahariev/go-with-prow/                                      0.0s
 => [build 3/4] COPY . ./                                                                                      0.0s
 => [build 4/4] RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags "$(build/ld  2.7s
 => [test-prepare 1/1] RUN CGO_ENABLED=0 go test -i ./...                                                      5.9s
 => CANCELED [test-main 1/2] RUN CGO_ENABLED=0 go test . -ginkgo.noColor -v                                    7.6s
 => [test-pckg02 1/2] RUN CGO_ENABLED=0 go test ./pkg02 -ginkgo.noColor -v                                     7.6s
 => ERROR [test-pckg01 1/2] RUN CGO_ENABLED=0 go test ./pkg01 -ginkgo.noColor -v                               7.6s
 => [test-pckg02 2/2] RUN touch finished.test-pckg02                                                           0.2s 
------                                                                                                              
 > [test-pckg01 1/2] RUN CGO_ENABLED=0 go test ./pkg01 -ginkgo.noColor -v:                                          
#14 7.462 === RUN   TestPkg01Suite                                                                                  
#14 7.462 Running Suite: Pkg01 Suite                                                                                
#14 7.462 ==========================                                                                                
#14 7.462 Random Seed: 1592409908
#14 7.463 Will run 1 of 1 specs
#14 7.463 
#14 7.465 • Failure [0.001 seconds]
#14 7.465 helloPkg01 function
#14 7.466 /go/src/github.com/dzahariev/go-with-prow/pkg01/pkg01_test.go:20
#14 7.466   output
#14 7.466   /go/src/github.com/dzahariev/go-with-prow/pkg01/pkg01_test.go:21
#14 7.466     should be Hello from Pkg02! [It]
#14 7.466     /go/src/github.com/dzahariev/go-with-prow/pkg01/pkg01_test.go:22
#14 7.466 
#14 7.466     Expected
#14 7.466         <string>: Hello from Pkg01!
#14 7.466     to be identical to
#14 7.466         <string>: Hello from Pkg02!
#14 7.466 
#14 7.466     /go/src/github.com/dzahariev/go-with-prow/pkg01/pkg01_test.go:23
#14 7.466 ------------------------------
#14 7.466 
#14 7.466 
#14 7.467 Summarizing 1 Failure:
#14 7.467 
#14 7.467 [Fail] helloPkg01 function output [It] should be Hello from Pkg02! 
#14 7.467 /go/src/github.com/dzahariev/go-with-prow/pkg01/pkg01_test.go:23
#14 7.467 
#14 7.467 Ran 1 of 1 Specs in 0.004 seconds
#14 7.467 FAIL! -- 0 Passed | 1 Failed | 0 Pending | 0 Skipped
#14 7.467 --- FAIL: TestPkg01Suite (0.01s)
#14 7.467 FAIL
#14 7.469 FAIL	dzahariev/go-with-prow/pkg01	0.010s
#14 7.471 FAIL
------
failed to solve with frontend dockerfile.v0: failed to build LLB: executor failed running [/bin/sh -c CGO_ENABLED=0 go test ./pkg01 -ginkgo.noColor -v]: runc did not terminate sucessfully
```
