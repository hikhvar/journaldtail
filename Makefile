.PHONY: journaldtail docker

# Set HUBUSER to build an image that you can push to a registry
HUBUSER ?= local

# build inside the Docker container, then make a runtime image
docker:
	docker build -t $(HUBUSER)/journaldtail .

journaldtail:
	go build -o journaldtail main.go
