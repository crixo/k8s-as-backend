docker build -t crixo/k8s-as-backend-todo-app:v0.0.0 .
docker run -it --rm -p 5000:80 --name todo-app crixo/todo-app:v0.0.0
# create cluster first using ../kind/create-cluster.sh
kind load docker-image crixo/todo-app:v0.0.0 --name k8s-as-backend --nodes="k8s-as-backend-worker,k8s-as-backend-worker2"

## Add nuget package
dotnet add TodoApi.csproj package Microsoft.Rest.ClientRuntime

## k8s api proxy
kubectl proxy --port=8080

## using kubectl as sidecar
running kubectl as sidecar, if you access the api using ```kubectl proxy --token ... --port 8080``` you need to supply the [bearer token](https://kubernetes.io/docs/tasks/access-application-cluster/access-cluster/#accessing-the-api-from-a-pod) mounted at /var/run/secrets/kubernetes.io/serviceaccount/token

curl -H "Authorization: Bearer from /var/...." http://localhost:8080/api

## direct access from pod to api via k8s service
TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
curl --cacert /var/run/secrets/kubernetes.io/serviceaccount/ca.crt -H "Authorization: Bearer $TOKEN" https://kubernetes.default.svc/api

# cluster browsing
http://localhost/todo-app/swagger-ui/index.html