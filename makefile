build:
	go build -o bin/gomodoro
	cp bin/gomodoro /usr/local/bin

build_linux:
	echo "[INFO] building [Go]modoro for linux"
	GOOS=linux GOARCH=amd64 go build -o bin/gomodoro_linux

structure:
	mkdir -p ${HOME}/gomodoro
	cp config/config.yaml ${HOME}/gomodoro	

	mkdir -p /etc/gomodoro/
	cp assets/tone.mp3 /etc/gomodoro
	

install: structure build