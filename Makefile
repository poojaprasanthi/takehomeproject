service:
	go get github.com/go-redis/redis
	go build readwrite.go
	go run readwrite.go $(path)
test:
	go test

