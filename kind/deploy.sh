#!/bin/sh

if [ "$#" -ne 2 ]; then

    read -r -p "kind cluster name(k8s-as-backend): " CLUSTER_NAME
    if [ -z "$CLUSTER_NAME" ]; then 
        # echo "KIND_CLUSTER_NAME is mandatory"
        # exit
        CLUSTER_NAME="k8s-as-backend"
    fi

    read -r -p "build image(y|N):  " BUILD
    if [[ ! $BUILD =~ ^(y)$ ]]; then 
        BUILD="N"
    fi
else
    CLUSTER_NAME=$1
    BUILD=$2
fi

# use default ns due to hardcoded ns in webhook-server/artifacts/webhook-created-signed-cert.sh
# NS=myns
# kubectl create ns $NS
# kubectl config set-context --current --namespace=$NS


#informer
echo "deploying informer"
cd ../informer
sh deploy.sh $CLUSTER_NAME $BUILD
cd ../kind

#TodoApp
echo "deploying TodoApp"
cd ../TodoApp
sh deploy.sh $CLUSTER_NAME $BUILD
cd ../kind

#webhook-server
echo "deploying webhook"
cd ../webhook-server
sh deploy.sh $CLUSTER_NAME $BUILD
cd ../kind

# Test
echo "getting todos "
kubectl get todos

STATUS_CODE=""
BASE_URL="http://localhost/todo-app"
URL="$BASE_URL/api/todo"
echo "url: $URL"
echo "swagger ui: $BASE_URL/swagger-ui/index.html"
while [ "$STATUS_CODE" != "200" ]
do
    STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" $URL)
    echo "$STATUS_CODE"
    sleep 10
done




