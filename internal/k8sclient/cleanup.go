package k8sclient

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Cleanup deletes every pod and service in the namespace, mirroring
// src/k8s/cleanup.ts. It logs and continues past individual delete failures
// rather than aborting, matching the previous best-effort behaviour.
func (c *Client) Cleanup(ctx context.Context) error {
	pods := c.Clientset.CoreV1().Pods(c.Namespace)
	services := c.Clientset.CoreV1().Services(c.Namespace)

	podList, err := pods.List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing pods: %w", err)
	}
	serviceList, err := services.List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing services: %w", err)
	}

	fmt.Printf("Cleaning up before exiting (%d pods, %d services)...\n", len(podList.Items), len(serviceList.Items))

	for _, pod := range podList.Items {
		fmt.Printf(" Deleting pod %s...\n", pod.Name)
		if err := pods.Delete(ctx, pod.Name, metav1.DeleteOptions{}); err != nil {
			fmt.Printf(" Failed to delete pod %s\n", pod.Name)
		}
	}

	for _, service := range serviceList.Items {
		fmt.Printf(" Deleting service %s...\n", service.Name)
		if err := services.Delete(ctx, service.Name, metav1.DeleteOptions{}); err != nil {
			fmt.Printf("  Failed to delete service %s\n", service.Name)
		}
	}

	return nil
}
