#!/bin/sh
read -r -p "build cluster(y | default N): " BUILD_CLUSTER
if [ -z "$BUILD_CLUSTER" ]; then 
    # echo "KIND_CLUSTER_NAME is mandatory"
    # exit
    BUILD_CLUSTER="N"
fi

NS=operator-running
# CR=example-k8sasbackend
# TARGET_URL="http://localhost/$NS/$CR/todo-app/swagger-ui/index.html"
# TODO_LIST_URL="http://localhost/$NS/$CR/todo-app/api/Todo"

if [[ ${BUILD_CLUSTER} == 'y' ]]; then
cd ../kind
sh create-cluster.sh
sh preload-image.sh

cd ../operator
fi

kubectl apply -f deploy/crds/k8s-as-backend.example.com_k8sasbackends_crd.yaml

kubectl create namespace $NS

echo "Open a new terminal"

echo "deploy cr"
echo "kubectl apply -n $NS -f deploy/crds/kab01.yaml"

echo "browse at http://localhost/$NS/kab01/todo-app/swagger-ui/index.html"

operator-sdk run --local --namespace=$NS

