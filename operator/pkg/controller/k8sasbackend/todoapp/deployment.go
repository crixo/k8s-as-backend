package todoapp

import (
	"fmt"

	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	authz "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/authz"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func (t TodoApp) ensureDeployment(i *k8sasbackendv1alpha1.K8sAsBackend) (*reconcile.Result, error) {
	resUtils := &common.ResourceUtils{
		Scheme:          t.Manager.GetScheme(),
		Client:          t.Manager.GetClient(),
		ResourceFactory: createDeployment,
	}

	nsn := types.NamespacedName{Name: deploymentName, Namespace: i.Namespace}
	found := &appsv1.Deployment{}
	return nil, common.EnsureResource(found, nsn, i, resUtils)
}

func createDeployment(resNamespacedName types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	image := "crixo/k8s-as-backend-todo-app:v0.0.0"
	var replicas int32 = 1

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resNamespacedName.Name,
			Namespace: resNamespacedName.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: matchingLabels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: matchingLabels,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName: authz.ServiceAccountName,
					Containers: []corev1.Container{{
						Image: image,
						//ImagePullPolicy: corev1.PullAlways,
						Name: deploymentName,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(todoPort),
							//Name:          "visitors",
						}},
						Env: []corev1.EnvVar{
							{
								Name:  "FullBasePath",
								Value: fmt.Sprintf("http://localhost%s", getAppBaseUrl(i)),
							},
							{
								Name:  "RoutePrefix",
								Value: "swagger-ui",
							},
							{
								Name:  "RelativeBasePath",
								Value: common.TrimFirstRune(getAppBaseUrl(i)), //todoAppUrlSegmentIdentifier,
							},
							{
								Name:  "UseSwagger",
								Value: "1",
							},
							//TODO: use the env var within the app
							{
								Name:  "KUBECTL_PROXY_PORT",
								Value: fmt.Sprint(kubectlApiPort),
							},
						},
					},
						{
							Image: "crixo/k8s-as-backend-informer:v.0.0.0",
							Name:  "informer",
							Ports: []corev1.ContainerPort{{
								ContainerPort: kubectlApiPort,
							}},
							Env: []corev1.EnvVar{
								{
									Name:  "USER",
									Value: "/root",
								},
								{
									Name: "NAMESPACE",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "metadata.namespace",
										},
									},
								},
							},
						},
						{
							Image: "bitnami/kubectl:1.16",
							Name:  "kubectl",
							Ports: []corev1.ContainerPort{{
								ContainerPort: kubectlApiPort,
							}},
							Command: []string{
								"/bin/sh",
								"-c",
								"kubectl proxy --port=8080",
							},
						}},
				},
			},
		},
	}
}
