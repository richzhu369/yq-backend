package main

import (
	v1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
	"sync"
)
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func ChangeWhitelist(ips, act string, clientSet *kubernetes.Clientset) {
	ipsToAdd := strings.Split(ips, ",")

	ingressList := GetAllIngress(clientSet)

	// 使用并发控制
	var wg sync.WaitGroup
	updateChan := make(chan error)

	for _, ingress := range ingressList.Items {
		wg.Add(1)
		go func(ingress v1.Ingress) {
			defer wg.Done()
			if act == "add" {
				err := AddIPsToWhitelist(clientSet, ingress.Namespace, ingress.Name, ipsToAdd)
				if err != nil {
					updateChan <- err
					return
				}
			} else if act == "del" {
				err := RemoveIPsFromWhitelist(clientSet, ingress.Namespace, ingress.Name, ipsToAdd)
				if err != nil {
					updateChan <- err
					return
				}
			}

		}(ingress)
	}

	go func() {
		wg.Wait()
		close(updateChan)
	}()

	// 处理并发操作结果
	for err := range updateChan {
		if err != nil {
			log.Fatalf("Failed to update Ingress whitelist: %v", err)
		}
	}

	log.Println("Ingress whitelist update completed.")
}
