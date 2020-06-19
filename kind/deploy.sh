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

#informer
sh ../informer/deploy.sh $CLUSTER_NAME $BUILD

#TodoApp
sh ../TodoApp/deploy.sh $CLUSTER_NAME $BUILD

#webhook-server
sh ../webhook-server/deploy.sh $CLUSTER_NAME $BUILD

# Test
kubectl get todos

STATUS_CODE=""
while [ "$STATUS_CODE" != "200" ]
do
    STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/todo-app/api/todo)
    echo "$STATUS_CODE"
    sleep 10
done




