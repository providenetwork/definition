GOC=go
GO111MODULE=on

DIRECTORIES=$(sort $(dir $(wildcard command/*/) ))
DIRECTORIES += $(sort $(dir $(wildcard schema/*/) $(wildcard validator/*/)))

INTERNAL_DIRECTORIES=$(sort $(dir $(wildcard internal/*/)))
MOCKS=$(foreach x, $(DIRECTORIES), mocks/$(x))
INTERNAL_MOCKS=$(foreach x, $(INTERNAL_DIRECTORIES), internal/mocks/$(x))

.PHONY: build test test_race lint vet install-deps coverage mocks clean-mocks


test:
	go test ./...

lint:
	golint $(go list ./... | grep -v mocks)

vet:
	go vet $(go list ./... | grep -v mocks)

clean-mocks:
	rm -rf mocks
	rm -rf internal/mocks

mocks: $(MOCKS) $(INTERNAL_MOCKS)
	
$(MOCKS): mocks/% : %
	mockery -output=$@ -dir=$^ -all

$(INTERNAL_MOCKS):  internal/mocks/% : %
	mockery -output=$@ -dir=$^ -all