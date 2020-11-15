export OPERATOR_IMAGE_NAME="crixo/k8s-as-backend-operator:v0.1.2"
export INGRESS_HOST="demo-k8s-as-backend.westeurope.cloudapp.azure.com"


envsubst < ../operator/deploy/operator-envsubst.yaml | kubectl apply -f - 