.PHONY: app debug

app:
	go build -o pequod ./cmd/main.go

gif: pequod
	vhs < demo.tape

debug:
	go build -o ./debug/pequod -gcflags "all=-N -l" ./cmd/main.go


