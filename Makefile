.PHONY: all
all: noop classify exec nonser tf tflite exec_genop torch

prepare:
	@go mod tidy
	@mkdir -p bin/

noop: prepare
	go build -o bin/noop examples/noop/main.go

classify: prepare
	go build -o bin/classify examples/classify/main.go

exec: prepare
	go build -o bin/exec examples/exec/main.go

nonser: prepare 
	go build -o bin/nonser examples/nonser/main.go

tf: prepare
	go build -o bin/tf examples/tf/main.go

tflite: prepare
	go build -o bin/tflite examples/tflite/main.go

exec_genop: prepare
	go build -o bin/exec_genop examples/exec_genop/main.go

torch:
	go build -o bin/torch examples/torch/main.go

clean:
	rm -rf bin
