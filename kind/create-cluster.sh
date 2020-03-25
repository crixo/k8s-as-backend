#!/bin/sh
read -r -p "cluster name(k8s-as-backend): " CLUSTER_NAME
if [ -z "$CLUSTER_NAME" ]; then 
    # echo "KIND_CLUSTER_NAME is mandatory"
    # exit
    CLUSTER_NAME="k8s-as-backend"
fi
kind create cluster --config 3nodes-ingress-controller.yaml --name $CLUSTER_NAME

nginxctrlimage='quay.io/kubernetes-ingress-controller/nginx-ingress-controller:master'
docker pull $nginxctrlimage
kind load docker-image $nginxctrlimage --name $CLUSTER_NAME --nodes='k8s-as-backend-control-plane'

kubectl apply -f mandatory.yaml
kubectl apply -f service-nodeport.yaml
kubectl patch deployments -n ingress-nginx nginx-ingress-controller --patch "$(cat nginx-ingress-controller-deployment-patch.yaml)"
kubectl apply -f usage.yaml

NGINX_POD_STATUS=""
while [ "$NGINX_POD_STATUS" != "Running" ]
do
    NGINX_POD_STATUS=$(kubectl get po --all-namespaces -l app.kubernetes.io/name=ingress-nginx -o jsonpath="{.items[*].status.phase}")
    echo "$NGINX_POD_STATUS"
    sleep 10
done

curl http://localhost/foo