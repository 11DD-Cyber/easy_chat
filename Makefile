.PHONY: run-user run-dev

# Run user RPC + API services
run-user:
	go run ./apps/user/rpc/user.go
	go run ./apps/user/api/user.go

# Start user / task / social / im services with one command
run-dev:
	go run ./cmd/devrun
