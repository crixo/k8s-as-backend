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

if [ $BUILD = 'y' ]; then 
    go mod vendor
    docker build -t crixo/k8s-as-backend-webhook-server:v.0.0.0 .
fi

kind load docker-image crixo/k8s-as-backend-webhook-server:v.0.0.0 --name $KIND_CLUSTER_NAME --nodes="k8s-as-backend-worker,k8s-as-backend-worker2"

sh ./artifacts/webhook-created-signed-cert.sh
cd artifacts
sh ./webhook-patch-ca-bundle.sh
cd ..

kubectl apply -f artifacts/deployment.yaml,artifacts/service.yaml,artifacts/webhook-registration.yaml  
