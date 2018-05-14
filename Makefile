
test:
	go test -v `go list ./... | grep -v test-client`
