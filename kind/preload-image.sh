#!/bin/sh
read -r -p "cluster name(k8s-as-backend): " CLUSTER_NAME
if [ -z "$CLUSTER_NAME" ]; then 
    # echo "CLUSTER_NAME is mandatory"
    # exit
    CLUSTER_NAME="k8s-as-backend"
fi

declare -a arr=("crixo/k8s-as-backend-todo-app:v0.0.0" 
                "crixo/k8s-as-backend-informer:v.0.0.0"
                "crixo/k8s-as-backend-webhook-server:v.0.0.0"
                "bitnami/kubectl:1.16"
                )

for image in "${arr[@]}"
do
   echo "IMAGE: $image"
   docker pull $image
   #command=(kind load docker-image $image --name $CLUSTER_NAME --nodes='k8s-as-backend-worker,k8s-as-backend-worker2')
   #"${command[@]}"
   kind load docker-image $image --name $CLUSTER_NAME --nodes='k8s-as-backend-worker,k8s-as-backend-worker2'
done
