read -r -p "kind cluster name(k8s-as-backend): " KIND_CLUSTER_NAME
if [ -z "$KIND_CLUSTER_NAME" ]; then 
    # echo "KIND_CLUSTER_NAME is mandatory"
    # exit
    KIND_CLUSTER_NAME="k8s-as-backend"
fi

go mod vendor
docker build -t crixo/k8s-as-backend-informer:v.0.0.0 .
kind load docker-image crixo/k8s-as-backend-informer:v.0.0.0 --name $KIND_CLUSTER_NAME
kubectl apply -f artifacts/crd.yaml,artifacts/app.yaml
sleep 10
kubectl apply -f artifacts/todo.yaml  