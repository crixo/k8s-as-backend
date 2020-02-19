kind create cluster --config 3nodes-ingress-controller.yaml --name k8s-as-backend
kubectl apply -f mandatory.yaml
kubectl apply -f service-nodeport.yaml
kubectl patch deployments -n ingress-nginx nginx-ingress-controller --patch "$(cat nginx-ingress-controller-deployment-patch.yaml)"