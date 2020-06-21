if [ "$#" -ne 2 ]; then

    read -r -p "kind cluster name(k8s-as-backend): " CLUSTER_NAME
    if [ -z "$CLUSTER_NAME" ]; then 
        # echo "KIND_CLUSTER_NAME is mandatory"
        # exit
        CLUSTER_NAME="k8s-as-backend"
    fi

    read -r -p "build image(y|N):  " BUILD
    if [[ ! $BUILD =~ ^(y)$ ]]; then 
        BUILD="n"
    fi
else
    CLUSTER_NAME=$1
    BUILD=$2
fi

SCRIPT_PATH=$(dirname $0)

if [ $BUILD = 'y' ]; then 
    docker build -t crixo/k8s-as-backend-todo-app:v0.0.0 .
fi
kind load docker-image crixo/k8s-as-backend-todo-app:v0.0.0 --name $CLUSTER_NAME --nodes="k8s-as-backend-worker,k8s-as-backend-worker2"
kubectl apply -f "$SCRIPT_PATH/artifacts/app.yaml"

echo "deploy TodoApp completed"