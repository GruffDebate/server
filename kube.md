kubectl apply -f $PATH/kube/
kubectl port-forward $(kubectl get pod --selector="app=server" --output jsonpath='{.items[0].metadata.name}') 8080:8080