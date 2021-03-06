.PHONY: test build cli job web test generate install

install:
	go install github.com/golang/mock/mockgen@v1.6.0

build:
	cli
	web
	job

cli:
	go build -o ./cli-conv ./cmd/cli

web:
	go build -o ./web ./cmd/web

job:
	go build -o ./job ./cmd/job

test-build:
	docker-compose -f docker-compose-test.yml build currency_conversion_test 

test:
	docker-compose -f docker-compose-test.yml run  currency_conversion_test go test ./... -cover -timeout 10s 

test-cleanup:
	docker-compose -f docker-compose-test.yml  down

test-down:
	docker-compose -f docker-compose-test.yml  rm -sf  currency_conversion_test
	make testapp-cleanup

generate:
	mockgen -package=middlewares -source=store/store.go  > web/middlewares/store_mocks_test.go
	mockgen -package=handlers -source=store/store.go  > web/handlers/store_mocks_test.go
	mockgen -package=job -source=store/store.go  > job/store_mocks_test.go
	mockgen -package=job -source=converter/converter.go  > job/converter_mocks_test.go
