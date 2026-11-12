.PHONY: all build test web-build clean demo

all: web-build build test

web-build:
	cd web && npm install && npm run build

build:
	go build -v -o bin/uxswarm.exe ./cmd/uxswarm

test:
	go test -v -race ./...

demo:
	./bin/uxswarm.exe demo --port 3000

run-test:
	./bin/uxswarm.exe test --demo --config examples/personas_b2b_saas.yaml --report usability_report.html

clean:
	rm -rf bin/ dist/ web/dist/ web/node_modules/ *.db *.html *.json
