.PHONY: build
build:
	go build -o ./dist/gmg .

.PHONY: install
install:
	cp dist/gmg /usr/local/bin/gmg && chmod 755 /usr/local/bin/gmg

.PHONY: all
all: build install

.PHONY: clean
clean:
	-rm -rfi dist
