# Setup AKS Cluster on Azure
```
az group create --name myResourceGroup --location eastus
az aks create --resource-group myResourceGroup --name myAKSCluster --node-count 3 --generate-ssh-keys

az aks get-credentials --resource-group myResourceGroup --name myAKSCluster
```

# Setup Azure File Storage and Kubernetes Service Class
Reference: https://docs.microsoft.com/en-us/azure/aks/azure-files-dynamic-pv

Inside the same resource group as your kubernetes nodes, create the Azure File storage.

```
az storage account create --resource-group MC_myAKSCluster_myAKSCluster_eastus --name mystorageaccount
```
Create the storage class using the .yaml file.
```
kubectl create -f azure-file-sc.yaml
```

