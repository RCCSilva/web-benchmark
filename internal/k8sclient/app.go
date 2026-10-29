package k8sclient

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type AppSpec struct {
	Slug         string
	Image        string
	DBHost       string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
}

type AppResult struct {
	NodePort int32
	Delete   func(ctx context.Context) error
}

// CreateApp creates the app Pod and its NodePort Service, mirroring
// src/k8s/app.ts.
func (c *Client) CreateApp(ctx context.Context, spec AppSpec) (AppResult, error) {
	podName := spec.Slug + "-app"
	serviceName := spec.Slug + "-srv"

	pods := c.Clientset.CoreV1().Pods(c.Namespace)
	services := c.Clientset.CoreV1().Services(c.Namespace)

	_, err := pods.Create(ctx, &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   podName,
			Labels: map[string]string{"app": podName},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  podName,
					Image: spec.Image,
					Resources: corev1.ResourceRequirements{
						Requests: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("1Gi"),
							corev1.ResourceCPU:    resource.MustParse("1000m"),
						},
						Limits: corev1.ResourceList{
							corev1.ResourceMemory: resource.MustParse("1Gi"),
							corev1.ResourceCPU:    resource.MustParse("1000m"),
						},
					},
					Env: []corev1.EnvVar{
						{Name: "DB_USER", Value: spec.PostgresUser},
						{Name: "DB_PASSWORD", Value: spec.PostgresPass},
						{Name: "DB_NAME", Value: spec.PostgresDB},
						{Name: "DB_HOST", Value: spec.DBHost},
						{Name: "DB_PORT", Value: "5432"},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return AppResult{}, fmt.Errorf("creating app pod: %w", err)
	}

	service, err := services.Create(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: serviceName},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": podName},
			Type:     corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{Name: "app-port", Protocol: corev1.ProtocolTCP, Port: 3000, TargetPort: intstr.FromInt(3000)},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return AppResult{}, fmt.Errorf("creating app service: %w", err)
	}

	if len(service.Spec.Ports) == 0 || service.Spec.Ports[0].NodePort == 0 {
		return AppResult{}, fmt.Errorf("failed to set up %s", serviceName)
	}

	return AppResult{
		NodePort: service.Spec.Ports[0].NodePort,
		Delete: func(ctx context.Context) error {
			podErr := pods.Delete(ctx, podName, metav1.DeleteOptions{})
			svcErr := services.Delete(ctx, serviceName, metav1.DeleteOptions{})
			if podErr != nil {
				return podErr
			}
			return svcErr
		},
	}, nil
}
