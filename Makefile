#
# Makefile
# @author Hans-Peter Schadler <hps@abyle.org>
# Initial concept for Makefile stolen from https://github.com/yyyar/gobetween/tree/master/dist (thanks!)
#

.PHONY: update clean build build-all test authors dist vendor build-container-latest \
build-container-tagged build-container-gitcommit release-container release-container-gitcommit

NAME := abylebotter
VERSION := $(shell cat VERSION)
COMPTIME := $(shell date -Is)
LDFLAGS := -X main.version=${VERSION} -X main.compTime=${COMPTIME}
SRCPATH := .
DOCKERBASETAG := hpsch/abylebotter
CURRENTGITCOMMIT := $(shell git log -1 --format=%h)
CURRENTGITUNTRACKED := $(shell git diff-index --quiet HEAD -- || echo "_untracked")

default: build

clean:
	@echo Cleaning up...
	@rm bin/* -rf
	@rm dist/* -rf
	@echo Done.

build:
	@echo Building...
	go build -o ./bin/$(NAME) -ldflags '${LDFLAGS}'
	@echo Done.

race:
	@echo Building...
	go build -o ./bin/$(NAME) -ldflags '${LDFLAGS}' -race
	@echo Done.

build-static:
	@echo Building...
	CGO_ENABLED=0 go build -o ./bin/$(NAME) -ldflags '-s -w --extldflags "-static" ${LDFLAGS}'
	@echo Done.

test:
	@echo "Running unit tests"
	@go test -covermode=count -coverprofile=coverage.out ./...

test-verbose:
	@echo "Running unit tests"
	@go test -v -covermode=count -coverprofile=coverage.out ./...

install: build
	install -d ${DESTDIR}/usr/local/bin/
	install -m 755 ./bin/${NAME} ${DESTDIR}/usr/local/bin/${NAME}

uninstall:
	rm -f ${DESTDIR}/usr/local/bin/${NAME}

authors:
	@git log --format='%aN <%aE>' | LC_ALL=C.UTF-8 sort | uniq -c -i | sort -nr | sed "s/^ *[0-9]* //g" > AUTHORS
	@cat AUTHORS

deps:
	go get -v ./...

clean-dist:
	rm -rf ./dist/${VERSION}

dist:
	@# For linux 386 when building on linux amd64 you'll need 'libc6-dev-i386' package
	@echo Building dist
	# we need this for Windows
	GOOS=windows GOARCH=386 go get -v github.com/konsorten/go-windows-terminal-sequences

	@#             os    arch  cgo ext
	@for arch in "linux   386  1      "  "linux   amd64 1      "  \
				 "windows 386  0 .exe "  "windows amd64 0 .exe "  \
				 "darwin  386  0      "  "darwin  amd64 0      "; \
	do \
	  set -- $$arch ; \
	  echo "******************* $$1_$$2 ********************" ;\
	  distpath="./dist/${VERSION}/$$1_$$2" ;\
	  mkdir -p $$distpath ; \
	  CGO_ENABLED=$$3 GOOS=$$1 GOARCH=$$2 go build -v -o $$distpath/$(NAME)$$4 -ldflags '-s -w --extldflags "-static" ${LDFLAGS}' ${SRCPATH}/*.go ;\
	  cp "README.md" "LICENSE" "CHANGELOG.md" "AUTHORS" $$distpath ;\
	  cp "cfg/config.toml" $$distpath/config_example.toml ;\
	  if [ "$$1" = "linux" ]; then \
		  cd $$distpath && tar -zcvf ../../${NAME}_${VERSION}_$$1_$$2.tar.gz * && cd - ;\
	  else \
		  cd $$distpath && zip -r ../../${NAME}_${VERSION}_$$1_$$2.zip . && cd - ;\
	  fi \
	done

build-container-latest: build-static
	@echo Building docker image ${DOCKERBASETAG}:latest
	docker build -t ${DOCKERBASETAG}:latest .

build-container-tagged: build-static
	@echo Building docker image ${DOCKERBASETAG}:${VERSION}
	docker build -t ${DOCKERBASETAG}:${VERSION} .

build-container-gitcommit: build-static
	@echo Building docker image ${DOCKERBASETAG}:${CURRENTGITCOMMIT}${CURRENTGITUNTRACKED}
	docker build -t ${DOCKERBASETAG}:${VERSION}-${CURRENTGITCOMMIT}${CURRENTGITUNTRACKED} .

release-container: build-container-tagged
	@echo Pushing docker image ${DOCKERBASETAG}:${VERSION}
	docker push ${DOCKERBASETAG}:${VERSION}

release-container-gitcommit: build-container-tagged
	@echo Pushing docker image ${DOCKERBASETAG}:${VERSION}
	docker push ${DOCKERBASETAG}:${VERSION}
