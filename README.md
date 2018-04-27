# bosh-cpi-go: Library for writing BOSH CPIs in Go

See [docs/example.go](docs/example.go) for an example & [apiv1/interfaces.go](apiv1/interfaces.go) for interface details.

CPIs using this library:

- [Warden CPI](https://github.com/cppforlife/bosh-warden-cpi-release)
- [VirtualBox CPI](https://github.com/cppforlife/bosh-virtualbox-cpi-release)


### build steps

You need `dep` for pulling in all Golang dependencies:

```cmd
go get -u github.com/golang/dep/cmd/dep
```

And compile

```cmd
git clone https://github.com/plooploops/bosh-cpi-go && cd bosh-cpi-go
dep ensure
go build docs/kubernetes-cpi.go
```

### Testing Steps

We will need to run docs/test.sh (we can use Bash or Windows Subsystem).  This will use the json templates to create stemcells and pods.

### Pending Items

Create Stemcell

