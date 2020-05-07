#
# Makefile
# @author Hans-Peter Schadler <hps@abyle.org>
# Initial concept for Makefile stolen from https://github.com/yyyar/gobetween/tree/master/dist (thanks!)
#

.PHONY: update clean build build-all test authors dist vendor build-container-latest \
build-container-tagged build-container-gitcommit release-container release-container-gitcommit

NAME := abylebotter
BINARIES := botterinstance bottercontrol botter
VERSION := $(shell cat VERSION)
COMPTIME := $(shell date -Is)
LDFLAGS := -X main.version=${VERSION} -X main.compTime=${COMPTIME}
SRCPATH := .
DOCKERBASETAG := torlenor/abylebotter
CURRENTGITCOMMIT := $(shell git log -1 --format=%h)
CURRENTGITUNTRACKED := $(shell git diff-index --quiet HEAD -- || echo "_untracked")
ENVFLAGS := GONOSUMDB=git.abyle.org/redseligg GONOPROXY=git.abyle.org/redseligg

default: build

clean:
	@echo Cleaning up...
	@rm bin/* -rf
	@rm dist/* -rf
	@echo Done.

build:
	@echo Building...
	@mkdir -p ./bin
	@for cmd in ${BINARIES}; \
	do \
		echo "\t$${cmd}" ;\
		${ENVFLAGS} go build -o ./bin/$${cmd} -ldflags '${LDFLAGS}' ./cmd/$${cmd}/ ;\
	done
	@echo Done.

race:
	@echo Building...
	mkdir -p ./bin
	@for cmd in ${BINARIES}; \
	do \
		echo "\t$${cmd}" ;\
		${ENVFLAGS} go build -o ./bin/$${cmd} -ldflags '${LDFLAGS}' -race ./cmd/$${cmd}/ ;\
	done
	@echo Done.

build-static:
	@echo Building...
	mkdir -p ./bin
	@for cmd in ${BINARIES}; \
	do \
		echo "\t$${cmd}" ;\
		CGO_ENABLED=0 ${ENVFLAGS} go build -o ./bin/$${cmd} -ldflags '-s -w --extldflags "-static" ${LDFLAGS}' ./cmd/$${cmd}/ ;\
	done
	@echo Done.

test:
	@echo "Running unit tests"
	@${ENVFLAGS} go test -covermode=count -coverprofile=coverage.out ./...

test-verbose:
	@echo "Running unit tests"
	@${ENVFLAGS} go test -v -covermode=count -coverprofile=coverage.out ./...

install: build
	install -d ${DESTDIR}/usr/local/bin/
	@for cmd in ${BINARIES}; \
	do \
		echo "\t$${cmd}" ;\
		install -m 755 ./bin/$${cmd} ${DESTDIR}/usr/local/bin/$${cmd} ;\
	done

uninstall:
	@for cmd in ${BINARIES}; \
	do \
		echo "\t$${cmd}" ;\
		rm -f ${DESTDIR}/usr/local/bin/$${cmd} ;\
	done

deps:
	${ENVFLAGS} go get -v ./...

clean-dist:
	rm -rf ./dist/${VERSION}

dist:
	@# For linux 386 when building on linux amd64 you'll need 'libc6-dev-i386' package
	@echo Building dist
	# we need this for Windows
	GOOS=windows GOARCH=386 ${ENVFLAGS} go get -v github.com/konsorten/go-windows-terminal-sequences

	@#             os    arch  cgo ext
	@for arch in "linux   386  1      "  "linux   amd64 1      "  \
				 "windows 386  0 .exe "  "windows amd64 0 .exe "  \
				 "darwin  386  0      "  "darwin  amd64 0      "; \
	do \
		set -- $$arch ; \
		echo "******************* $$1_$$2 ********************" ;\
		distpath="./dist/${VERSION}/$$1_$$2" ;\
		mkdir -p $$distpath ; \
		for cmd in ${BINARIES}; \
		do \
			echo "\t$${cmd}" ;\
			CGO_ENABLED=$$3 GOOS=$$1 GOARCH=$$2 ${ENVFLAGS} go build -o $$distpath/$${cmd}$$4 -ldflags '-s -w --extldflags "-static" ${LDFLAGS}' ./cmd/$${cmd}/ ;\
		done ;\
		cp "README.md" "LICENSE" "CHANGELOG.md" "AUTHORS" $$distpath ;\
		cp "cfg/bots.toml" $$distpath/bots.toml ;\
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
	@echo Building docker image ${DOCKERBASETAG}:${VERSION}-${CURRENTGITCOMMIT}${CURRENTGITUNTRACKED}
	docker build -t ${DOCKERBASETAG}:${VERSION}-${CURRENTGITCOMMIT}${CURRENTGITUNTRACKED} .

release-container: build-container-tagged
	@echo Pushing docker image ${DOCKERBASETAG}:${VERSION}
	docker tag ${DOCKERBASETAG}:${VERSION} ${DOCKERBASETAG}:latest
	docker push ${DOCKERBASETAG}:${VERSION}

release-container-tagged: build-container-tagged
	@echo Pushing docker image ${DOCKERBASETAG}:${VERSION}
	docker push ${DOCKERBASETAG}:${VERSION}
	docker tag ${DOCKERBASETAG}:${VERSION} ${DOCKERBASETAG}:latest
	docker push ${DOCKERBASETAG}:latest

release-container-gitcommit: build-container-gitcommit
	@echo Pushing docker image ${DOCKERBASETAG}:${VERSION}-${CURRENTGITCOMMIT}${CURRENTGITUNTRACKED}
	docker push ${DOCKERBASETAG}:${VERSION}-${CURRENTGITCOMMIT}${CURRENTGITUNTRACKED}

