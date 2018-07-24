all: build

build:
	vgo build -o bin/cheapskate

run: build
	./bin/cheapskate
