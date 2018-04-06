# Contributing

## Setup your machine

`download-rds-logs` is written in [Go](https://golang.org/).

Prerequisites are:

* Build:
  * `make`
  * [Go 1.8+](http://golang.org/doc/install)

Clone `download-rds-logs` from source into `$GOPATH`:

```sh
$ go get github.com/niranjan94/download-rds-logs
$ cd $GOPATH/src/github.com/niranjan94/download-rds-logs
```

Install the build dependencies:

``` sh
$ make setup
```

## Test your change

You can create a branch for your changes and try to build from the source as you go:

``` sh
$ go build
```

## Submit a pull request

Push your branch to your `download-rds-logs` fork and open a pull request against the
master branch.
