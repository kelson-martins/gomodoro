build:
	go build -o bin/gomodoro

build_linux:
	echo "[INFO] building [Go]modoro for linux"
	GOOS=linux GOARCH=amd64 go build -o bin/gomodoro_linux