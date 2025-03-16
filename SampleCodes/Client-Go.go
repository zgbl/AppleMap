package main

import (
	"context"
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// 连接 Kubernetes API
	config, err := rest.InClusterConfig() // 适用于 Pod 内部运行
	if err != nil {
		log.Fatalf("Failed to get in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}

	for {
		// 监控 Apple Maps 服务的 Deployment 副本数量
		deploymentsClient := clientset.AppsV1().Deployments("apple-maps-namespace")
		deployment, err := deploymentsClient.Get(context.TODO(), "maps-backend", metav1.GetOptions{})
		if err != nil {
			log.Printf("Failed to get deployment: %v", err)
			continue
		}

		fmt.Printf("Maps Backend Replicas: %d\n", *deployment.Spec.Replicas)

		// 如果 CPU 使用率高，可自动扩容
		if *deployment.Spec.Replicas < 5 {
			newReplicas := int32(*deployment.Spec.Replicas + 1)
			deployment.Spec.Replicas = &newReplicas
			_, err := deploymentsClient.Update(context.TODO(), deployment, metav1.UpdateOptions{})
			if err != nil {
				log.Printf("Failed to scale up: %v", err)
			} else {
				fmt.Println("Scaled up Maps Backend to", newReplicas)
			}
		}

		time.Sleep(10 * time.Second)
	}
}
