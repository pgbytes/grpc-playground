SERVICE_TARGET=echo
PROTO_SOURCE_DIR=proto/chat
PROTO_GENERATED_DIR=proto/chat

SERVICE_GO_SOURCES:=$(shell find ./service -iname "*.go")

PROTO_SOURCES:=$(shell find ${PROTO_SOURCE_DIR} -iname "*.proto")
PROTO_GENERATED:=$(subst $(PROTO_SOURCE_DIR),$(PROTO_GENERATED_DIR),$(subst .proto,.pb.go,$(PROTO_SOURCES)))

.PHONY: all
all: $(SERVICE_TARGET)

$(SERVICE_TARGET): $(SERVICE_GO_SOURCES) $(PROTO_GENERATED)
	@echo "Compiling $@..."
	cd service; go build; go vet ./...

$(PROTO_GENERATED): $(PROTO_SOURCES)
	@echo "Generating ${PROTO_SOURCE_DIR} sources..."
	protoc --proto_path=${PROTO_SOURCE_DIR} ${PROTO_SOURCE_DIR}/*.proto --go_out=plugins=grpc:${PROTO_SOURCE_DIR}

proto-gen: $(PROTO_GENERATED)

# Build docker image, tag it, and push it to the docker registry.
# This assumes that you have run "make install"
# and moved the executable. For instance:
# make install
# mv $GOPATH/bin/service _images/service/
docker-build-and-push: _images/service/service
	cd _images/service/; ./build.sh

# Install binaries at $GOPATH/bin.	# This will pick up all Go packages in the project and build and
# install them.  Everything that is not a Go package will be ignored.
install: ${SERVICE_TARGET}
	go install bitbucket.org/egym-com/adonis-gateway/service

vet: $(SERVICE_TARGET)
	go vet ./service/...

clean:
	@echo "Cleaning.."
	rm -f $(PROTO_GENERATED_D\IR)/*.go
	cd service; go clean