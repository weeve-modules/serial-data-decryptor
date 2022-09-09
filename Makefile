MODULE=weevenetwork/serial-data-decryptor
VERSION_NAME=v1.0.0

create_image:
	docker build -t ${MODULE}:${VERSION_NAME} . -f docker/Dockerfile
.PHONY: create_image

run_image:
	docker run -p 80:80 --rm --env-file=./.env ${MODULE}:${VERSION_NAME}
.PHONY: run_image

debug_image:
	docker run -p 80:80 --rm --env-file=./.env --entrypoint /bin/bash -it ${MODULE}:${VERSION_NAME}
.PHONY: debug_image

create_and_push_multi_platform:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v6,linux/arm/v7 -t ${MODULE}:${VERSION_NAME} --push . -f docker/Dockerfile
.PHONY: create_and_push_multi_platform
