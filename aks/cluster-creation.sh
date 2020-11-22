CLUSTER_NAME=$1
dnsnameInt=$2

if [ -z "$CLUSTER_NAME" ]; then 
    read -r -p "cluster name(k8s-as-backend): " CLUSTER_NAME
    if [ -z "$CLUSTER_NAME" ]; then 
        CLUSTER_NAME="k8s-as-backend"
    fi
fi

if [ -z "$dnsnameInt" ]; then 
    read -r -p "dns name(demo-k8s-as-backend): " dnsnameInt
    if [ -z "$dnsnameInt" ]; then 
        dnsnameInt="demo-k8s-as-backend"
    fi
fi

export AKS_CLUSTER_NAME=$CLUSTER_NAME
export AZURE_SUBSCRIPTION_ID="d9e06499-49d3-4d60-b301-3ff03e019bb7" #vs
export AZURE_RESOURCE_GROUP="k8s-as-backend"
export AZURE_REGION="westeurope"
export DNSNAME=$dnsnameInt # FQDN will then be DNSNAME.ZONE.cloudapp.azure.com
export K8S_VERSION="1.16.15"

az login
az account set -s $AZURE_SUBSCRIPTION_ID
az group create -l $AZURE_REGION -n $AZURE_RESOURCE_GROUP
az aks create \
    --resource-group $AZURE_RESOURCE_GROUP \
    --name $AKS_CLUSTER_NAME \
    --kubernetes-version=$K8S_VERSION \
    --node-count 2 \
    --load-balancer-sku standard \
    --node-vm-size Standard_B2s \
    --vm-set-type VirtualMachineScaleSets \
    --generate-ssh-keys
    #--node-osdisk-size 30 \

az aks get-credentials --resource-group $AZURE_RESOURCE_GROUP --name $AKS_CLUSTER_NAME

CLUSTER_RG=$(az aks show --resource-group $AZURE_RESOURCE_GROUP --name $AKS_CLUSTER_NAME --query nodeResourceGroup -o tsv)
export PUBLIC_STATIC_IP=$(az network public-ip create --resource-group $CLUSTER_RG --name myAKSPublicIP --sku Standard --allocation-method static --query publicIp.ipAddress -o tsv)
echo $PUBLIC_STATIC_IP #51.138.48.67
# export PUBLIC_STATIC_IP=52.149.111.196
#envsubst < cloud-generic.yaml | cat - 
kubectl apply -f mandatory.yaml
envsubst < cloud-generic.yaml | kubectl apply -f - 

kubectl get service --namespace ingress-nginx
kubectl get pods --namespace ingress-nginx

# exit

#export DNSNAME="demo-k8s-as-backend" # FQDN will then be DNSNAME.ZONE.cloudapp.azure.com
PUBLICIPID=$(az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, '$PUBLIC_STATIC_IP')].[id]" --output tsv)

az network public-ip update --ids $PUBLICIPID --dns-name $DNSNAME

envsubst < usage.yaml | kubectl apply -f - 

## install cert-manager + letsecrypt
kubectl create ns cert-manager
#kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v0.13.0/cert-manager.yaml
kubectl apply --validate=false -f https://github.com/jetstack/cert-manager/releases/download/v1.0.4/cert-manager.yaml
# if returns the following error
# Error from server (InternalError): error when creating "le-cluster-issuer.yaml": 
# Internal error occurred: failed calling webhook "webhook.cert-manager.io": 
# Post https://cert-manager-webhook.cert-manager.svc:443/mutate?timeout=30s: 
# dial tcp 10.0.78.42:443: connect: connection refused
# try to apply it againg until you are able to get the resource
# k get clusterissuers.cert-manager.io 
# NAME          READY   AGE
# letsencrypt   True    4s

cic=0
while [ $cic -eq 0 ]
do
    kubectl apply -f le-cluster-issuer.yaml
    cic=$(kubectl get clusterissuers.cert-manager.io -o json | jq '.items | length')
    echo "clusterissuers: $cic"
done
envsubst < le-ingress.yaml | kubectl apply -f - 

echo "DONE!"


# ode=0 -- Original Error: Code="PublicIPAndLBSkuDoNotMatch" Message="Basic sku load balancer /subscriptions/300aa066-33d1-4cd8-9cac-ef9082a33e4b/resourceGroups/mc_test-cri_k8s-as-backend_weste
# urope/providers/Microsoft.Network/loadBalancers/kubernetes cannot reference Standard sku publicIP /subscriptions/300aa066-33d1-4cd8-9cac-ef9082a33e4b/resourceGroups/MC_test-cri_k8s-as-backend
# _westeurope/providers/Microsoft.Network/publicIPAddresses/myAKSPublicIP." Details=[]   