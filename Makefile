# CONFIG_PATH=${HOME}/.ntms/

# .PHONY: init
# init:
# 	mkdir -p ${CONFIG_PATH}

# .PHONY: test
# 	go test -race ./...

# build_docker image
TAG ?= 0.0.1

build-docker:
	docker build -t docker.io/pmoth/ntms:$(TAG) .

# push_docker image to docker hub 
push-docker:
	docker push docker.io/pmoth/ntms:$(TAG)

