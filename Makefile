TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(TAG:v%=%)

ifneq ($(COMMIT), $(TAG_COMMIT))
    VERSION := $(VERSION)-$(COMMIT)
endif

ifeq ($(VERSION),)
	VERSION := $(COMMIT)-$(DATA)
endif

ifneq ($(shell git status --porcelain),)
    VERSION := $(VERSION)-dirty
endif


run: build
	docker run --env-file .env -it --rm cryptovoxels/compressor:latest

build:
	docker build -t cryptovoxels/compressor:latest -t cryptovoxels/compressor:$(VERSION) .

push:
	docker push cryptovoxels/compressor:latest
	docker push cryptovoxels/compressor:$(VERSION)
