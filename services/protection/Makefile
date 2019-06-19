.PHONY: pb

pb:
	for f in proto/*.proto; do \
		protoc --proto_path=${GOPATH}/src:. --go_out=plugins=grpc:. $$f; \
		echo compiled: $$f; \
	done
