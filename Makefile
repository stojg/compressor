run: build
	docker run --env-file .env -it --rm cryptovoxels/compressor:latest

build:
	docker build -t cryptovoxels/compressor:latest cryptovoxels/compressor:$(shell git rev-parse --verify HEAD) .

push:
	docker push cryptovoxels/compressor:latest
	docker push cryptovoxels/compressor:$(shell git rev-parse --verify HEAD)
