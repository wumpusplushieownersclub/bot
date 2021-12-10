start:
	go run ./src
dev:
	reflex -r '\.go' -s -- sh -c "go run ./src"