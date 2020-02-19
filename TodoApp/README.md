docker build -t crixo/todo-app:v0.0.0 .
docker run -it --rm -p 5000:80 --name todo-app crixo/todo-app:v0.0.0
# create cluster first using ../kind/create-cluster.sh
kind load docker-image todo-app:v0.0.0 --name k8s-as-backend

## Add nuget package
dotnet add TodoApi.csproj package Microsoft.Rest.ClientRuntime

## k8s api proxy
kubectl proxy --port=8080

