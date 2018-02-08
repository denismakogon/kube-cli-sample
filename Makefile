all: deps build

deps:
	glide install -v

build:
	go build -o kube-configmaps

first-time: deps
	$(MAKE) install

.PHONY: deps build
