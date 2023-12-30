package controller

import (
	api "github.com/shahriarKhanpiyal/custom-resource/api/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newDeployment(crd *api.CustomResource, deploymentName string) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: crd.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(crd, api.GroupVersion.WithKind("CustomResource")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: crd.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":  crd.Name,
					"Kind": "CustomResource",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":  crd.Name,
						"Kind": "CustomResource",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  crd.Name,
							Image: crd.Spec.Container.Image,
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: crd.Spec.Container.Port,
								},
							},
						},
					},
				},
			},
		},
	}
}

func newService(crd *api.CustomResource, serviceName string) *corev1.Service {
	labels := map[string]string{
		"app":  crd.Name,
		"Kind": "CustomResource",
	}
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: crd.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(crd, api.GroupVersion.WithKind("CustomResource")),
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Port:       crd.Spec.Container.Port,
					TargetPort: intstr.FromInt(int(crd.Spec.Container.Port)),
					Protocol:   "TCP",
					NodePort:   crd.Spec.Service.ServiceNodePort,
				},
			},
			Type: func() corev1.ServiceType {
				if crd.Spec.Service.ServiceType == "NodePort" {
					return corev1.ServiceTypeNodePort
				} else {
					return corev1.ServiceTypeClusterIP
				}
			}(),
		},
	}
}
