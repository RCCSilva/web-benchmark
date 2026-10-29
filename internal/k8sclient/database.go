package k8sclient

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type DatabaseResult struct {
	NodePort    int32
	ServiceName string
	Delete      func(ctx context.Context) error
}

// CreateDatabase creates a Postgres Pod plus a ClusterIP Service (used by the
// app pod as DB_HOST) and a NodePort Service (used by the orchestrator to
// reach it from outside the cluster), mirroring src/k8s/database.ts.
func (c *Client) CreateDatabase(ctx context.Context, slug, postgresUser, postgresPass, postgresDB string) (DatabaseResult, error) {
	podName := slug + "-db"
	serviceName := slug + "-db-srv"
	externalServiceName := slug + "-db-external-srv"

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
					Name:  "postgres",
					Image: "postgres:14-alpine",
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
						{Name: "POSTGRES_USER", Value: postgresUser},
						{Name: "POSTGRES_PASSWORD", Value: postgresPass},
						{Name: "POSTGRES_DB", Value: postgresDB},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return DatabaseResult{}, fmt.Errorf("creating database pod: %w", err)
	}

	_, err = services.Create(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: serviceName},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": podName},
			Ports: []corev1.ServicePort{
				{Name: "database-port", Protocol: corev1.ProtocolTCP, Port: 5432, TargetPort: intstr.FromInt(5432)},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return DatabaseResult{}, fmt.Errorf("creating database service: %w", err)
	}

	externalService, err := services.Create(ctx, &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: externalServiceName},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": podName},
			Type:     corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				{Name: "db-port", Protocol: corev1.ProtocolTCP, Port: 5432, TargetPort: intstr.FromInt(5432)},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return DatabaseResult{}, fmt.Errorf("creating external database service: %w", err)
	}

	if len(externalService.Spec.Ports) == 0 || externalService.Spec.Ports[0].NodePort == 0 {
		return DatabaseResult{}, fmt.Errorf("failed to set up %s", externalServiceName)
	}

	return DatabaseResult{
		NodePort:    externalService.Spec.Ports[0].NodePort,
		ServiceName: serviceName,
		Delete: func(ctx context.Context) error {
			podErr := pods.Delete(ctx, podName, metav1.DeleteOptions{})
			svcErr := services.Delete(ctx, serviceName, metav1.DeleteOptions{})
			extErr := services.Delete(ctx, externalServiceName, metav1.DeleteOptions{})
			for _, err := range []error{podErr, svcErr, extErr} {
				if err != nil {
					return err
				}
			}
			return nil
		},
	}, nil
}
