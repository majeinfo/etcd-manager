FROM golang:1.24-bookworm as builder
RUN apt-get update && apt-get install git -y
WORKDIR /
RUN git clone -b main https://github.com/majeinfo/etcd-manager.git
WORKDIR /etcd-manager/backend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -o etcd-manager . 

FROM scratch
COPY --from=builder --chmod=755 /etcd-manager/backend/etcd-manager /
ENTRYPOINT ["/etcd-manager"]
