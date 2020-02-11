.PHONY: deploy linux

deploy: linux
	rsync -Paz GeoLite2-Country.mmdb nodes.yaml build/lihaiguo.bin $$HOST:~

linux:
	mkdir -p build
	GOOS=linux go build -o build/lihaiguo.bin github.com/hayeah/lihaiguo