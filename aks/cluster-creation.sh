read -r -p "cluster name(k8s-as-backend): " CLUSTER_NAME
if [ -z "$CLUSTER_NAME" ]; then 
    # echo "KIND_CLUSTER_NAME is mandatory"
    # exit
    CLUSTER_NAME="k8s-as-backend"
fi

export AKS_CLUSTER_NAME=$CLUSTER_NAME
export AZURE_SUBSCRIPTION_ID="d9e06499-49d3-4d60-b301-3ff03e019bb7" #vs
export AZURE_RESOURCE_GROUP="test-cri"
export AZURE_REGION="westeurope"
export DNSNAME="demo-k8s-as-backend" # FQDN will then be DNSNAME.ZONE.cloudapp.azure.com

az login
az account set -s $AZURE_SUBSCRIPTION_ID
az group create -l $AZURE_REGION -n $AZURE_RESOURCE_GROUP
az aks create \
    --resource-group $AZURE_RESOURCE_GROUP \
    --name $AKS_CLUSTER_NAME \
    --node-count 2 \
    --load-balancer-sku standard \
    --node-vm-size Standard_B2s \
    --vm-set-type VirtualMachineScaleSets \
    #--node-osdisk-size 30 \
    --generate-ssh-keys
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

exit

export DNSNAME="demo-k8s-as-backend" # FQDN will then be DNSNAME.ZONE.cloudapp.azure.com
PUBLICIPID=$(az network public-ip list --query "[?ipAddress!=null]|[?contains(ipAddress, '$PUBLIC_STATIC_IP')].[id]" --output tsv)

az network public-ip update --ids $PUBLICIPID --dns-name $DNSNAME

envsubst < usage.yaml | kubectl apply -f - 

echo "DONE!"

# kubectl config unset "users.clusterUser_$AZURE_RESOURCE_GROUP_$CLUSTER_NAME"
# kubectl config unset "contexts.$CLUSTER_NAME"
# kubectl config unset "clusters.$CLUSTER_NAME"

# ode=0 -- Original Error: Code="PublicIPAndLBSkuDoNotMatch" Message="Basic sku load balancer /subscriptions/300aa066-33d1-4cd8-9cac-ef9082a33e4b/resourceGroups/mc_test-cri_k8s-as-backend_weste
# urope/providers/Microsoft.Network/loadBalancers/kubernetes cannot reference Standard sku publicIP /subscriptions/300aa066-33d1-4cd8-9cac-ef9082a33e4b/resourceGroups/MC_test-cri_k8s-as-backend
# _westeurope/providers/Microsoft.Network/publicIPAddresses/myAKSPublicIP." Details=[]   