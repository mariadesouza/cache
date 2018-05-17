test:
	export GOCACHE=off && go test -v `go list ./... | grep -v test-client`
