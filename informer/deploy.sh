if [ "$#" -ne 2 ]; then

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
else
    KIND_CLUSTER_NAME=$1
    BUILD=$2
fi

# echo $KIND_CLUSTER_NAME
# echo $BUILD
# echo $0
# SCRIPT_PATH=$(dirname $0)
# cat "$SCRIPT_PATH/artifacts/crd.yaml"
# #$(cd $(dirname $0)/../../; pwd)
# exit

SCRIPT_PATH=$(dirname $0)
echo $SCRIPT_PATH

if [ $BUILD = 'y' ]; then 
    go mod vendor
    docker build -t crixo/k8s-as-backend-informer:v0.0.0 .
fi
kind load docker-image crixo/k8s-as-backend-informer:v0.0.0 --name $KIND_CLUSTER_NAME --nodes="k8s-as-backend-worker,k8s-as-backend-worker2"

kubectl apply -f "$SCRIPT_PATH/artifacts/crd.yaml"
kubectl apply -f "$SCRIPT_PATH/artifacts/app.yaml"
sleep 10
kubectl apply -f "$SCRIPT_PATH/artifacts/todo.yaml" 

echo "deploy informer completed"