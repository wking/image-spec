GO15VENDOREXPERIMENT=1
export GO15VENDOREXPERIMENT

DOCKER ?= $(shell which docker)
# These docs are in an order that determines how they show up in the PDF/HTML docs.
DOC_FILES := \
	README.md \
	code-of-conduct.md \
	project.md \
	media-types.md \
	manifest.md \
	serialization.md

FIGURE_FILES := \
	img/media-types.png

OUTPUT		?= output/
DOC_FILENAME	?= oci-image-spec

EPOCH_TEST_COMMIT ?= v0.2.0

default: help

help:
	@echo "Usage: make <target>"
	@echo
	@echo " * 'docs' - produce document in the $(OUTPUT) directory"
	@echo " * 'fmt' - format the json with indentation"
	@echo " * 'validate-examples' - validate the examples in the specification markdown files"
	@echo " * 'oci-image-tool' - build the oci-image-tool binary"
	@echo " * 'schema-fs' - regenerate the virtual schema http/FileSystem"
	@echo " * 'check-license' - check license headers in source files"
	@echo " * 'lint' - Execute the source code linter"
	@echo " * 'test' - Execute the unit tests"
	@echo " * 'update-deps' - Update vendored dependencies"
	@echo " * 'img/*.png' - Generate PNG from dot file"

fmt:
	for i in schema/*.json ; do jq --indent 2 -M . "$${i}" > xx && cat xx > "$${i}" && rm xx ; done

docs: $(OUTPUT)/$(DOC_FILENAME).pdf $(OUTPUT)/$(DOC_FILENAME).html

$(OUTPUT)/$(DOC_FILENAME).pdf: $(DOC_FILES) $(FIGURE_FILES)
	@mkdir -p $(OUTPUT)/ && \
	cp -ap img/ $(shell pwd)/$(OUTPUT)/&& \
	$(DOCKER) run \
	-it \
	--rm \
	-v $(shell pwd)/:/input/:ro \
	-v $(shell pwd)/$(OUTPUT)/:/$(OUTPUT)/ \
	-u $(shell id -u) \
	--workdir /input \
	vbatts/pandoc -f markdown_github -t latex -o /$(OUTPUT)/$(DOC_FILENAME).pdf $(patsubst %,/input/%,$(DOC_FILES)) && \
	ls -sh $(shell readlink -f $@)

$(OUTPUT)/$(DOC_FILENAME).html: $(DOC_FILES) $(FIGURE_FILES)
	@mkdir -p $(OUTPUT)/ && \
	cp -ap img/ $(shell pwd)/$(OUTPUT)/&& \
	$(DOCKER) run \
	-it \
	--rm \
	-v $(shell pwd)/:/input/:ro \
	-v $(shell pwd)/$(OUTPUT)/:/$(OUTPUT)/ \
	-u $(shell id -u) \
	--workdir /input \
	vbatts/pandoc -f markdown_github -t html5 -o /$(OUTPUT)/$(DOC_FILENAME).html $(patsubst %,/input/%,$(DOC_FILES)) && \
	ls -sh $(shell readlink -f $@)

code-of-conduct.md:
	curl -o $@ https://raw.githubusercontent.com/opencontainers/tob/d2f9d68c1332870e40693fe077d311e0742bc73d/code-of-conduct.md

validate-examples:
	go test -run TestValidate ./schema

oci-image-tool:
	go build ./cmd/oci-image-tool

example: cmd/example/main.go
	go build ./cmd/example

schema-fs:
	@echo "generating schema fs"
	@cd schema && printf "%s\n\n%s\n" "$$(cat ../.header)" "$$(go generate)" > fs.go

check-license:
	@echo "checking license headers"
	@./.tool/check-license

lint:
	@echo "checking lint"
	@./.tool/lint

test:
	go test -race ./...

## this uses https://github.com/Masterminds/glide and https://github.com/sgotti/glide-vc
update-deps:
	@which glide > /dev/null 2>/dev/null || (echo "ERROR: glide not found. Consider 'make install.tools' target" && false)
	glide update --strip-vcs --strip-vendor --update-vendored --delete
	glide-vc --only-code --no-tests
	# see http://sed.sourceforge.net/sed1line.txt
	find vendor -type f -exec sed -i -e :a -e '/^\n*$$/{$$d;N;ba' -e '}' "{}" \;

img/%.png: img/%.dot
	dot -Tpng $^ > $@

.PHONY: .gitvalidation

# When this is running in travis, it will only check the travis commit range
.gitvalidation:
	@which git-validation > /dev/null 2>/dev/null || (echo "ERROR: git-validation not found. Consider 'make install.tools' target" && false)
ifeq ($(TRAVIS),true)
	git-validation -q -run DCO,short-subject,dangling-whitespace
else
	git-validation -v -run DCO,short-subject,dangling-whitespace -range $(EPOCH_TEST_COMMIT)..HEAD
endif

.PHONY: install.tools

install.tools: .install.gitvalidation .install.glide .install.glide-vc

.install.gitvalidation:
	go get github.com/vbatts/git-validation

.install.glide:
	go get github.com/Masterminds/glide

.install.glide-vc:
	go get github.com/sgotti/glide-vc

clean:
	rm -rf *~ $(OUTPUT)
	rm -f oci-image-tool

.PHONY: \
	validate-examples \
	oci-image-tool \
	check-license \
	clean \
	lint \
	docs \
	test
