ARG GO_VERSION=1.22.2
ARG GIT_BRANCH=main
ARG GIT_COMMIT
ARG BUILD_DATE

FROM golang:${GO_VERSION} AS build-env
ARG GO_VERSION=1.22.2
ARG GIT_BRANCH=main
ARG GIT_COMMIT
ARG BUILD_DATE
WORKDIR /src

COPY ./go.mod .
COPY ./go.sum .
COPY cmd/server/config.json .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build  -ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.GitBranch=${GIT_BRANCH} -X main.GoVersion=${GO_VERSION} -X main.BuildDate=${BUILD_DATE}" -o main cmd/server/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=build-env  /src/main /app/main
COPY --from=build-env  /src/config.json /app/config.json

EXPOSE 8088

ENTRYPOINT ["./main"]