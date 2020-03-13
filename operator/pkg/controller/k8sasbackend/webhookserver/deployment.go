package webhookserver

import (
	k8sasbackendv1alpha1 "github.com/crixo/k8s-as-backend/operator/pkg/apis/k8sasbackend/v1alpha1"
	common "github.com/crixo/k8s-as-backend/operator/pkg/controller/k8sasbackend/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func (ws WebhookServer) ensureDeployment(i *k8sasbackendv1alpha1.K8sAsBackend) error {
	resUtils := &common.ResourceUtils{
		Scheme:          ws.Scheme,
		Client:          ws.Client,
		ResourceFactory: createDeployment,
	}

	found := &appsv1.Deployment{}
	return common.EnsureResource(found, "k8s-as-backend-webhook-server", i, resUtils)
}

func createDeployment(resourceName string, i *k8sasbackendv1alpha1.K8sAsBackend) runtime.Object {
	labels := map[string]string{
		"app": "k8s-as-backend-webhook-server",
	}
	image := "crixo/k8s-as-backend-webhook-server:v.0.0.0"
	var replicas int32 = 1

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resourceName,
			Namespace: i.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						//ImagePullPolicy: corev1.PullAlways,
						Name: "k8s-as-backend-webhook-server",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 443,
							//Name:          "visitors",
						}},
						Args: []string{
							"-tls-cert-file=/etc/webhook/certs/cert.pem",
							"-tls-private-key-file=/etc/webhook/certs/key.pem",
							"-v=2",
						},
						// Env: []corev1.EnvVar{
						// 	{
						// 		Name:  "MYSQL_DATABASE",
						// 		Value: "visitors",
						// 	},
						// },
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "webhook-certs",
							MountPath: "/etc/webhook/certs",
							ReadOnly:  true,
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: "webhook-certs",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: secretName,
							},
						},
					}},
				},
			},
		},
	}
}
