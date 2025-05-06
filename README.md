# Go bindings for vAccel

[Go](https://go.dev/) bindings for vAccel wrap the vAccel C API and provide a
native Go API to vAccel operations. The bindings are currently a WiP, supporting
a subset of the vAccel operations.

You can find more information about these bindings and everything vAccel in the
[Documentation](https://docs.vaccel.org).

## Usage

The Go bindings are implemented in the `vaccel` Go package.

### Requirements

- To use the `vaccel` Go package you need a valid vAccel installation. You can
  find more information on how to install vAccel in the
  [Installation](https://docs.vaccel.org/latest/getting-started/installation)
  page.

- This package requires Go 1.20 or newer. Verify your Go version with:
  ```sh
  go version
  ```
  and update Go as needed using the
  [official instructions](https://go.dev/doc/install).

### Using the `vaccel` package

You can use the package in your Go code like any other Go package with:

```go
import "github.com/nubificus/vaccel-go/vaccel"
```

### Running the examples

You can find examples in the `examples` directory. The provided examples are
similar to the C examples and you must configure vAccel in order to use them.

To run an image classification, like the C `classify`, set:

```sh
export VACCEL_PLUGINS=libvaccel-noop.so
```

and, assuming vAccel is installed at `/usr/local`, run with:

```console
$ go run github.com/nubificus/vaccel-go/examples/classify \
      /usr/local/share/vaccel/images/example.jpg
Output(1):  This is a dummy classification tag!
Output(2):  This is a dummy classification tag!
```
