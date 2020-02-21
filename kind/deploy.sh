#!/bin/sh

KIND_CLUSTER_NAME="k8s-as-backend"
BUILD="N"
#informer
sh ../informer/deploy.sh $KIND_CLUSTER_NAME $BUILD

#TodoApp
sh ../TodoApp/deploy.sh $KIND_CLUSTER_NAME $BUILD

#webhook-server
sh ../webhook-server/deploy.sh $KIND_CLUSTER_NAME $BUILD

# Test
kubectl get todos

STATUS_CODE=""
while [ "$STATUS_CODE" != "200" ]
do
    STATUS_CODE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost/todo-app/api/todo)
    echo "$STATUS_CODE"
    sleep 10
done




