# VERSION=$(shell git describe --tags)
GIT_COMMIT=${shell git rev-parse --short HEAD}
GO_VERSION=1.22
GIT_BRANCH=${shell git rev-parse --abbrev-ref HEAD}
BUILD_DATE=$(shell git log -n1 --pretty='format:%cd' --date=format:'%Y%m%d')
VU ?= 100
DURATION ?= 1m

.PHONY: run-server
run-server:
	cd cmd/server ; go run main.go

.PHONY: build
build: 
	go build  -v -ldflags="-X 'main.GitCommit=${GIT_COMMIT}' -X 'main.GoVersion=${GO_VERSION}'  -X 'main.BuildDate=${BUILD_DATE}'  -X 'main.GitBranch=${GIT_BRANCH}'" -o main ./cmd/server/main.go

.PHONY: swag-init
swag-init:
	swag init --parseDependency --parseInternal --parseDepth 1 -o api/docs -g cmd/server/main.go

.PHONY: docker-build
docker-build:
	docker build \
	--build-arg "GIT_COMMIT=${GIT_COMMIT}" \
	--build-arg "BUILD_DATE=$(BUILD_DATE)" \
	--build-arg "GO_VERSION=$(GO_VERSION)" \
	--build-arg "GIT_BRANCH=$(GIT_BRANCH)" \
	-t coupon_rush_server -f cmd/server/Dockerfile --no-cache . ;

.PHONY: k6-run
k6-run:
	k6 run -vu $(VU) --duration $(DURATION) api/docs/swagger-k6/script.js

.PHONY: go-test
go-test:
	.PHONY: go-test
go-test:
	go test ./... -coverprofile=coverage.out
	( \
		go tool cover -func=coverage.out & \
		go tool cover -html=coverage.out & \
		wait \
	)