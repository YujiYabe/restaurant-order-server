
GO_TOOLCHAIN := go1.27.1

# ----------------------------
.PHONY: stop
stop:
	docker-compose stop


# ----------------------------
.PHONY: removeall
removeall: stop
	docker stop $(docker ps -q)
	docker system prune -f


# ----------------------------
.PHONY: build
build:
	cd backend && go build -o app ./cmd

# ----------------------------
.PHONY: sqlite-init
sqlite-init:
	@test -n "$(STORE_ID)" || (echo "STORE_IDを指定してください: STORE_ID=0 make sqlite-init"; exit 1)
	cd backend/cmd/db && STORE_ID="$(STORE_ID)" go run .

# ----------------------------
.PHONY: doctor
doctor:
	./backend/script/doctor.sh

# ----------------------------
.PHONY: prepare
prepare:
	./backend/script/doctor.sh --prepare

# ----------------------------
.PHONY: start
start:
	./start.sh


# ----------------------------
.PHONY: up
up: redis-up firestore-up

# ----------------------------
.PHONY: redis-up
redis-up:
	@if [ -f backend/.env ]; then set -a; . ./backend/.env; set +a; fi; docker compose up -d --wait redis

# ----------------------------
.PHONY: redis-down
redis-down:
	docker compose stop redis

# ----------------------------
.PHONY: firestore-up
firestore-up:
	docker compose up -d --build --wait firestore

# ----------------------------
.PHONY: firestore-down
firestore-down:
	docker compose down


# ----------------------------
.PHONY: docker-build
docker-build:
	docker build -t backend -f ./backend/build/Dockerfile .


# ----------------------------
.PHONY: docker-up
docker-up:
	# docker run -itd --restart=always --name=backend -v $(pwd)/backend:/go/src/backend -p 6789:6789 -t backend
	# docker run --restart=always --name=backend -v $(pwd)/backend:/go/src/backend -p 6789:6789 -t backend /go/src/backend/script/entrypoint.sh





# ----------------------------
.PHONY: reup
reup: build up



# ----------------------------
.PHONY: restart
restart: stop up

# ----------------------------
.PHONY: resetall
# resetall: removeall build up
resetall: removeall up



# ----------------------------
.PHONY: lint
lint:
	cd backend && GOFLAGS=-mod=mod ./bin/golangci-lint run --config=.golangci.yml ./cmd/... ./internal/...

# ----------------------------
.PHONY: test-e2e
test-e2e:
	cd backend && go test -tags=e2e ./test/e2e

# ----------------------------
.PHONY: setup-git
setup-git:
	git config core.hooksPath .githooks

# ----------------------------
.PHONY: install-tools
install-tools:
	# Create bin directory if it doesn't exist
	mkdir -p backend/bin
	
	# Install protoc compiler
	# GOBIN=$(PWD)/backend/bin go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	# GOBIN=$(PWD)/backend/bin go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	
	# Install other tools
	GOTOOLCHAIN=$(GO_TOOLCHAIN) GOBIN=$(PWD)/backend/bin go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0
	GOTOOLCHAIN=$(GO_TOOLCHAIN) GOBIN=$(PWD)/backend/bin go install github.com/air-verse/air@latest
	GOTOOLCHAIN=$(GO_TOOLCHAIN) GOBIN=$(PWD)/backend/bin go install github.com/swaggo/swag/v2/cmd/swag@latest
