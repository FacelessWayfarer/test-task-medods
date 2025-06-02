
lint:
	golangci-lint run --config .golangci.yml ./...
mock:
	find ./internal -type d -name "*mock*" |  xargs rm -dfR
	find ./pkg -type d -name "*mock*" | xargs rm -dfR
	@ go generate ./...
testing:
	go test ./... -count=1
swagger:
	swag init -d cmd --pdl 3