build-arm:
	GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o powermon ./main.go

send: build-arm
	scp powermon pi@powermon.edjusted.com: