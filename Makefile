SHELl:=/bin/zsh

clean:
	rm -f ./dns-proxy

build: clean
	go build -o dns-proxy

run: build
	source <(cat $(PWD)/.env | awk '{print "export "$$1}') && sudo -E ./dns-proxy



