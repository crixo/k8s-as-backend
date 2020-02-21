read -r -p "cluster name(k8s-as-backend): " CLUSTER_NAME
if [ -z "$CLUSTER_NAME" ]; then 
    # echo "KIND_CLUSTER_NAME is mandatory"
    # exit
    CLUSTER_NAME="k8s-as-backend"
fi

read -r -p "resource group(test-cri): " AZURE_RESOURCE_GROUP
if [ -z "$AZURE_RESOURCE_GROUP" ]; then 
    # echo "KIND_CLUSTER_NAME is mandatory"
    # exit
    AZURE_RESOURCE_GROUP="test-cri"
fi

kubectl config unset 'users.clusterUser_'$AZURE_RESOURCE_GROUP'_'$CLUSTER_NAME
kubectl config unset "contexts.$CLUSTER_NAME"
kubectl config unset "clusters.$CLUSTER_NAME"
