
ios:
	go get -d golang.org/x/mobile/cmd/gomobile && \
	gomobile bind -target=ios	

android:
	go get -d golang.org/x/mobile/cmd/gomobile && \
	gomobile bind -target=android

js:
	GOOS=js GOARCH=wasm go build -ldflags="-w -s" -o ./static/main.wasm ./cmd/jsFunc.go
