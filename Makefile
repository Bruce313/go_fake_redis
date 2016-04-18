svr: *.go
	go build
cli: cli/*.go
	cd cli; go build
clean:
	rm cli/cli; rm go_fake_redis
