read -r -p "kind cluster name(k8s-as-backend): " KIND_CLUSTER_NAME
if [ -z "$KIND_CLUSTER_NAME" ]; then 
    # echo "KIND_CLUSTER_NAME is mandatory"
    # exit
    KIND_CLUSTER_NAME="k8s-as-backend"
fi

read -r -p "build image(y|N):  " BUILD
if [[ ! $BUILD =~ ^(y)$ ]]; then 
    BUILD="n"
fi


if [ $BUILD = 'y' ]; then 
    go mod vendor
    docker build -t crixo/k8s-as-backend-informer:v.0.0.0 .
fi
kind load docker-image crixo/k8s-as-backend-informer:v.0.0.0 --name $KIND_CLUSTER_NAME --nodes="k8s-as-backend-worker,k8s-as-backend-worker2"

kubectl apply -f artifacts/crd.yaml
kubectl apply -f artifacts/app.yaml
sleep 10
kubectl apply -f artifacts/todo.yaml  