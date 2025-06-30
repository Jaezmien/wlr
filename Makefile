FILENAME=wlr

all:
	make linux arm7

arm7:
	mkdir -p "dist/arm7"
	GOOS=linux GOARCH=arm64 GOARM=7 \
		CGO_ENABLED=0 \
		go build -o "dist/arm7/$(FILENAME)" .
	
linux:
	GOOS=linux OARCH=amd64 \
		GCGO_ENABLED=0 \
		go build -o "dist/linux/$(FILENAME)" .
