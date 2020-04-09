#!/bin/sh
read -r -p "build cluster(y | default N): " BUILD_CLUSTER
if [ -z "$BUILD_CLUSTER" ]; then 
    # echo "KIND_CLUSTER_NAME is mandatory"
    # exit
    BUILD_CLUSTER="N"
fi

NS=operator-test
CR=example-k8sasbackend
TARGET_URL="http://localhost/$NS/$CR/todo-app/swagger-ui/index.html"
TODO_LIST_URL="http://localhost/$NS/$CR/todo-app/api/Todo"

if [[ ${BUILD_CLUSTER} == 'y' ]]; then
cd ../kind
sh create-cluster.sh
sh sh preload-image.sh

cd ../operator
fi
kubectl apply -f deploy/crds/k8s-as-backend.example.com_k8sasbackends_crd.yaml

kubectl create namespace $NS

operator-sdk test local ./test/e2e \
--namespace $NS \
--up-local --go-test-flags "-skipcleanup=true" 
#--global-manifest "./deploy/crds/k8s-as-backend.example.com_k8sasbackends_crd.yaml"

echo "Curling $TARGET_URL"
STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/operator-test/example-k8sasbackend/todo-app/swagger-ui/index.html)
while [ "$STATUS_CODE" != "200" ]
do
    echo $STATUS_CODE
    sleep 10
done

TODOS_COUNT=$(curl -s $TODO_LIST_URL | jq '. | length')
echo "TODOS_COUNT: $TODOS_COUNT\n"

curl -X POST "http://localhost/operator-test/example-k8sasbackend/todo-app/api/Todo" -H "accept: text/plain" -H "Content-Type: application/json" -d "{\"id\":\"3fa85f64-5717-4562-b3fc-2c963f66afa6\",\"code\":\"string\",\"when\":\"2020-04-08T22:35:34.956Z\",\"message\":\"string3\"}"
echo "\nTODO CREATED\n"

TODOS_COUNT=$(curl -s $TODO_LIST_URL | jq '. | length')
echo "TODOS_COUNT: $TODOS_COUNT\n"