export NS=$1
export DNSNAME=$2

#echo "NS: $NS - DNSNAME: $DNSNAME"

if [ -z "$NS" ]; then 
    read -r -p "namespace(operator-in-cluster): " NS
    if [ -z "$NS" ]; then 
        NS="operator-in-cluster"
    fi
fi

if [ -z "$DNSNAME" ]; then 
    read -r -p "dns name(demo-k8s-as-backend): " DNSNAME
    if [ -z "$DNSNAME" ]; then 
        DNSNAME="demo-k8s-as-backend"
    fi
fi

#echo "NS: $NS - DNSNAME: $DNSNAME"
#exit

cd ../operator

# Deploy operator CRD 
kubectl apply -f deploy/crds/k8s-as-backend.example.com_k8sasbackends_crd.yaml

# Create and configure namespace for this demo
kubectl create namespace $NS
kubectl config set-context --current --namespace $NS

# Deploy RBAC resources for the operator app running in cluster
kubectl apply -f deploy/service_account.yaml
#TODO: replace hardcoded ns
envsubst < deploy/cluster_role_binding_cluster_admin-with-var.yaml | kubectl apply -f - 

# Deploy operator
export OPERATOR_IMAGE_NAME="crixo/k8s-as-backend-operator:v0.0.0"
export INGRESS_HOST="$DNSNAME.westeurope.cloudapp.azure.com"
envsubst < deploy/operator-envsubst.yaml | kubectl apply -f - 

# Deploy the operator CR IOW your solution/app instance
kubectl apply -f deploy/crds/kab01.yaml
