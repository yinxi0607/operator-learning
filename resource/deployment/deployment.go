package deployment

import (
	appv1 "github.com/yinxi0607/operator-learning/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func New(app *appv1.App, env string) *appsv1.Deployment {
	labels := map[string]string{"app.example.com/v1": app.Name + "-" + env}
	selector := &metav1.LabelSelector{MatchLabels: labels}
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{APIVersion: "apps/v1", Kind: "Deployment"},
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.Name + "-" + env,
			Namespace: app.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(app, schema.GroupVersionKind{
					Group:   appv1.GroupVersion.Group,
					Version: appv1.GroupVersion.Version,
					Kind:    "App",
				}),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: app.Spec.Replicas,
			Selector: selector,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: newContainers(app, env),
				},
			},
		},
	}
}

func newContainers(app *appv1.App, env string) []corev1.Container {
	app.Spec.Envs = append(app.Spec.Envs, corev1.EnvVar{
		Name:  "envTest",
		Value: env,
	})
	var containerPorts []corev1.ContainerPort
	for _, servicePort := range app.Spec.Ports {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: servicePort.TargetPort.IntVal,
		})
	}
	return []corev1.Container{
		{
			Name:            app.Name,
			Image:           app.Spec.Image,
			Resources:       app.Spec.Resource,
			Env:             app.Spec.Envs,
			Ports:           containerPorts,
			ImagePullPolicy: corev1.PullIfNotPresent,
		},
	}
}
