# CONFIG_PATH=${HOME}/.ntms/

# .PHONY: init
# init:
# 	mkdir -p ${CONFIG_PATH}

# .PHONY: test
# 	go test -race ./...

# build_docker image
TAG ?= 0.1.0

.PHONY: uninstall-helm
uninstall-helm:
	helm uninstall ntms --debug

.PHONY: install-helm
install-helm:
	helm install ntms ./helm/ntms --set image.tag=$(TAG) --debug

.PHONY: build-go
build-go:
	go build -o ./cmd/ntms/ntms -ldflags "-X main.version=$(TAG)" ./cmd/ntms/main.go

.PHONY: run-go
run-go:
	go run -ldflags "-X main.version=$(TAG)" ./cmd/ntms/main.go

.PHONY: build-docker
build-docker:
	docker build -t github.com/pmoth/ntms:$(TAG) .
	docker tag github.com/pmoth/ntms:$(TAG) docker.io/pmoth/ntms:$(TAG)

.PHONY: push-docker
push-docker:
	docker push docker.io/pmoth/ntms:$(TAG)

.PHONY: run-docker-httpd
run-docker-httpd:
 # This works because the Dockerfile has already copied the config.yaml into the container
	docker run -it --rm \
	    -p 8401:8401 \
		-p 8400:8400 \
		-p 8080:8080 github.com/pmoth/ntms:$(TAG)

.PHONY: run-docker
run-docker:
 # This works because the Dockerfile has already copied the config.yaml into the container
	docker run -it --rm \
	    -p 8401:8401 \
		-p 8400:8400 \
		-p 8080:8080 github.com/pmoth/ntms:$(TAG)

PHONY: clean-docker
clean-docker:
	docker rmi -f github.com/pmoth/ntms:$(TAG)
	docker rmi -f docker.io/pmoth/ntms:$(TAG)

PHONY: clean-hub
clean-hub:
	docker rmi -f docker.io/pmoth/ntms:$(TAG)