web:
	env GOOS=js GOARCH=wasm go build -o game.wasm src/main.go
	cp `go env GOROOT`/misc/wasm/wasm_exec.js .
	zip -u SpaceDeflation.zip game.wasm index.html wasm_exec.js 

windows:
	GOOS=windows go build src/main.go

serve:
	python -m http.server 8080
