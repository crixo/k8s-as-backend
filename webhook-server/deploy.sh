read -r -p "kind cluster name: " KIND_CLUSTER_NAME
if [ -z "$KIND_CLUSTER_NAME" ]; then 
    echo "KIND_CLUSTER_NAME is mandatory"
    exit
fi

echo "KIND_CLUSTER_NAME: $KIND_CLUSTER_NAME"
#exit

go mod vendor

docker build -t crixo/k8s-as-backend-webhook-server:v.0.0.0 .

kind load docker-image crixo/k8s-as-backend-webhook-server:v.0.0.0 --name $KIND_CLUSTER_NAME

sh ./artifacts/webhook-created-signed-cert.sh
cd artifacts
sh ./webhook-patch-ca-bundle.sh
cd ..

kubectl apply -f artifacts/deployment.yaml,artifacts/service.yaml,artifacts/webhook-registration.yaml  
