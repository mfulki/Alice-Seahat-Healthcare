start:
	go run .
dev:
	nodemon -e go --exec go run . --signal SIGTERM
mock:
	mockery --all