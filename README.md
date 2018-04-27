# bosh-cpi-go: Library for writing BOSH CPIs in Go

See [docs/example.go](docs/example.go) for an example & [apiv1/interfaces.go](apiv1/interfaces.go) for interface details.

CPIs using this library:

- [Warden CPI](https://github.com/cppforlife/bosh-warden-cpi-release)
- [VirtualBox CPI](https://github.com/cppforlife/bosh-virtualbox-cpi-release)

### What we're attempting to do:

![](images/poddiagram.png?raw=true)

The green and orange boxes describe what we're attempting to stand up with Bosh CPIs.  We'd like to utilize Azure Files and AKS Pods to back a bosh Kubernetes CPI implementation.

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

We can use Bash (WSL / Linux):
```cmd
docs/test.sh
```

This will use the json templates to create stemcells and pods.

### Check Running Pods

```cmd
kubectl --kubeconfig kubeconfig proxy
```

Navigate to aka.ms/k8sui to check out the pods in a browser.
![](images/podbrowser.png?raw=true)

We can also use kubectl to get into a running pod.

```cmd
kubectl exec -it trustypod -- bash
```

### Pending Items

- Create Stemcell
- Fill in other parts of CPI interface