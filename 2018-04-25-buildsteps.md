
### build steps on Windows


```cmd
go get k8s.io/apimachinery/pkg/api/errors
go get k8s.io/client-go/kubernetes
go get github.com/cppforlife/bosh-cpi-go/apiv1
go get k8s.io/client-go/tools/clientcmd

cd %USERPROFILE%
mkdir source
mkdir source\repos
cd source\repos
git clone https://github.com/plooploops/bosh-cpi-go
cd apiv1
go build
cd ..\rpc
cd ..\docs
``` 


go get -u github.com/golang/dep/cmd/dep

