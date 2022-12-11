.PHONY: app

app:
	go build -o pequod ./cmd/main.go

gif: pequod
	vhs < demo.tape


