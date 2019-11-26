GOC=go
GO111MODULE=on

DIRECTORIES=$(sort $(dir $(wildcard ./*/*/)))
MOCKS=$(foreach x, $(DIRECTORIES), mocks/$(x))


.PHONY: build test test_race lint vet install-deps coverage mocks clean-mocks


test:
	go test ./...

lint:
	golint $(go list ./... | grep -v mocks)

vet:
	go vet $(go list ./... | grep -v mocks)

clean-mocks:
	rm -rf mocks

mocks: $(MOCKS)
	
$(MOCKS): mocks/% : %
	mockery -output=$@ -dir=$^ -all
	
#install-mock:
#	go get github.com/golang/mock/gomock
#	go install github.com/golang/mock/mockgen

