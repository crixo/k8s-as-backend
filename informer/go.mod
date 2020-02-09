module github.com/crixo/k8s-as-backend/informer

go 1.13

require (
	github.com/crixo/k8s-as-backend/library v0.0.0
	go.uber.org/zap v1.13.0
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
)

replace github.com/crixo/k8s-as-backend/library v0.0.0 => ../library
