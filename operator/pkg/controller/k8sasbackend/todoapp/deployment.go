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

	deploymentName := common.CreateUniqueSecondaryResourceName(i, BaseName)
	nsn := types.NamespacedName{Name: deploymentName, Namespace: i.Namespace}
	found := &appsv1.Deployment{}
	err := common.EnsureResource(found, nsn, i, resUtils)
	log.Info("Found Deployment", "Namespace", found.Namespace, "Name", found.Name)
	if err == nil && found.Name != "" {
		size := i.Spec.Size
		// TODO: handle i.Spec.ProductVersion changes
		if size != *found.Spec.Replicas {
			log.Info("Reconciling Size", "DesiredSize", size, "CurrentSize", &found.Spec.Replicas)
			found.Spec.Replicas = &size
			//TODO:
			return common.UpdateResource(resUtils.Client, found)
		}
	}
	return nil, err
}

func createDeployment(resNamespacedName types.NamespacedName, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	image := common.CreateImageName(i, todoAppImage)
	informerImage := common.CreateImageName(i, common.InformerImage)
	var replicas int32 = i.Spec.Size
	matchingLabels := common.CreateMatchingLabels(i, BaseName)
	serviceAccountName := common.CreateUniqueSecondaryResourceName(i, authz.BaseName)
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
					ServiceAccountName: serviceAccountName,
					Containers: []corev1.Container{{
						Image: image,
						//ImagePullPolicy: corev1.PullAlways,
						Name: "todo-app",
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(todoPort),
							//Name:          "visitors",
						}},
						Env: []corev1.EnvVar{
							//TODO: create env variable mapping the meta.namespace to have the api calling the ralted ns while is creating the todo-CR
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
							{
								Name: "NAMESPACE",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
							{
								Name: "POD_NAME",
								ValueFrom: &corev1.EnvVarSource{
									FieldRef: &corev1.ObjectFieldSelector{
										FieldPath: "metadata.name",
									},
								},
							},
							{
								//TODO: modify C# TodoApp using this var to create the todo resource for the api-server
								Name:  "OPERATOR_NAME",
								Value: i.Name,
							},
						},
					},
						{
							Image: informerImage,
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
									Name:  "TODO_APP_SVC",
									Value: fmt.Sprintf("http://localhost:%d", todoPort),
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
							Image: kubectlImage,
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
