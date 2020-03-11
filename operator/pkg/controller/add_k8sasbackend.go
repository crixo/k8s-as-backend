package controller

import (
	"github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, k8sasbackend.Add)
}
