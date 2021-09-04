install:
	go install ./cmd/shellbridge

demo: install
	go run ./_examples/demo/main.go