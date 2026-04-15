# specify base image and use multi-stage build
FROM golang:1.24.3-alpine AS build
WORKDIR /go/src/proglog
COPY . .
# fully static compilation and generate the program to target directory
RUN CGO_ENABLED=0 go build -o /go/bin/proglog ./cmd/proglog

# start a Docker image from an empty base(scratch)
FROM scratch
# --from enables you to copy file from other image
COPY --from=build /go/bin/proglog /bin/proglog
# makes the container run /bin/proglog on start
ENTRYPOINT ["/bin/proglog"]

