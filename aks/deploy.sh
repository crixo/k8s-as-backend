
#informer
docker rmi crixo/k8s-as-backend-informer
dokcer tag crixo/k8s-as-backend-informer:v.0.0.0 crixo/k8s-as-backend-informer
docker push crixo/k8s-as-backend-informer

kubectl apply -f ../informer/artifacts/crd.yaml
sleep 10
kubectl apply -f ../informer/artifacts/todo.yaml
kubectl apply -f ../informer/artifacts/app-aks.yaml

#TodoApp
docker rmi crixo/k8s-as-backend-informer
dokcer tag crixo/k8s-as-backend-informer:v.0.0.0 crixo/k8s-as-backend-informer
docker push crixo/k8s-as-backend-informer

kubectl apply -f ../TodoApp/artifacts/app-aks.yaml

#webhook-server
docker rmi crixo/k8s-as-backend-webhook-server
dokcer tag crixo/k8s-as-backend-webhook-server:v.0.0.0 crixo/k8s-as-backend-webhook-server
docker push crixo/k8s-as-backend-webhook-server

sh ../webhook-server/artifacts/webhook-created-signed-cert.sh
#sh ./artifacts/webhook-patch-ca-bundle.sh
export CA_BUNDLE=$(kubectl config view --raw -o json | jq -r '.clusters[] | select(.name == "'$(kubectl config current-context)'") | .cluster."certificate-authority-data"')
kubectl apply -f ../webhook-server/artifacts/deployment.yaml,../webhook-server/artifacts/service.yaml
envsubst < ../webhook-server/artifacts/webhook-registration-template-aks.yaml | kubectl apply -f -  

