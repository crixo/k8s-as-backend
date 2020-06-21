if [ "$#" -ne 2 ]; then

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
else
    KIND_CLUSTER_NAME=$1
    BUILD=$2
fi

SCRIPT_PATH=$(dirname $0)

if [ $BUILD = 'y' ]; then 
    go mod vendor
    docker build -t crixo/k8s-as-backend-webhook-server:v0.0.0 .
fi

kind load docker-image crixo/k8s-as-backend-webhook-server:v0.0.0 --name $KIND_CLUSTER_NAME --nodes="k8s-as-backend-worker,k8s-as-backend-worker2"

sh "$SCRIPT_PATH/artifacts/webhook-created-signed-cert.sh"
#sh ./artifacts/webhook-patch-ca-bundle.sh
export CA_BUNDLE=$(kubectl config view --raw -o json | jq -r '.clusters[] | select(.name == "'$(kubectl config current-context)'") | .cluster."certificate-authority-data"')
kubectl apply -f "$SCRIPT_PATH/artifacts/deployment.yaml"
kubectl apply -f "$SCRIPT_PATH/artifacts/service.yaml"
envsubst < "$SCRIPT_PATH/artifacts/webhook-registration-template.yaml" | kubectl apply -f -  

echo "deploy webhook completed"