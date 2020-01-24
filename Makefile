GOC=go
GO111MODULE=on

DIRECTORIES=$(sort $(dir $(wildcard command/*/) ))
DIRECTORIES += $(sort $(dir $(wildcard schema/*/) $(wildcard validator/*/)))
MOCKS=$(foreach x, $(DIRECTORIES), mocks/$(x))

INTERNAL_DIRECTORIES=$(shell find ./pkg -type d | sed -e 's/\.\/pkg\///g' | grep -v pkg)
INTERNAL_MOCKS=$(foreach x, $(INTERNAL_DIRECTORIES), pkg/mocks/$(x)/)

.PHONY: build test test_race lint vet install-deps coverage mocks clean-mocks

test:
	go test ./...

lint:
	golint $(go list ./... | grep -v mocks)

vet:
	go vet $(go list ./... | grep -v mocks)

clean-mocks:
	rm -rf mocks
	rm -rf pkg/mocks

mocks: $(MOCKS) $(INTERNAL_MOCKS)

$(MOCKS): mocks/% : %
	mockery -output=$@ -dir=$^ -all

$(INTERNAL_MOCKS):  pkg/mocks/% : pkg/%
	mockery -output=$@ -dir=$^ -all
