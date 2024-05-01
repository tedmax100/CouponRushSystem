# VERSION=$(shell git describe --tags)
GIT_COMMIT=${shell git rev-parse --short HEAD}
GO_VERSION=1.22
GIT_BRANCH=${shell git rev-parse --abbrev-ref HEAD}
BUILD_DATE=$(shell git log -n1 --pretty='format:%cd' --date=format:'%Y%m%d')
VU ?= 300
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

.PHONY: add_new_coupon_active
add_new_coupon_active:
	docker cp build/add_new_coupon_active.sql db:/tmp/add_new_coupon_active.sql
	docker exec -it db psql -U userabc -d coupon -f /tmp/add_new_coupon_active.sql

.PHONY: docker-run-server
docker-run-server:
	docker compose up --build --force-recreate --remove-orphans --detach


.PHONY: k6-run-reserve
k6-run-reserve:
	k6 run -vu $(VU) --duration $(DURATION) api/docs/swagger-k6/reserve.js

.PHONY: k6-run-breaking-test-reserve
k6-run-breaking-test-reserve:
	k6 run api/docs/swagger-k6/breaking-test-reserve.js


.PHONY: k6-run-purchase
k6-run-purchase:
	k6 run -vu $(VU) --duration 1m api/docs/swagger-k6/purchase.js

.PHONY: go-test
go-test:
	go test ./... -coverprofile=coverage.out
	( \
		go tool cover -func=coverage.out & \
		go tool cover -html=coverage.out & \
		wait \
	)

