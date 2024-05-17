package main

import (
	"context"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

// GetAllNamespace 获取所有namespace
//func GetAllNamespace(clientSet *kubernetes.Clientset) {
//	// 获取所有命名空间的列表
//	nsList, err := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
//	if err != nil {
//		panic(err.Error())
//	}
//
//	// 遍历每个命名空间并打印名称
//	for _, ns := range nsList.Items {
//		fmt.Printf("Namespace: %s\n", ns.Name)
//	}
//}

// GetAllIngress 获得所有namespace下的ingress
func GetAllIngress(clientSet *kubernetes.Clientset) *v1.IngressList {
	// 获取所有命名空间的 Ingress 列表
	ingressList, err := clientSet.NetworkingV1().Ingresses("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	return ingressList
}

// GetAllCWIngress 获得所有namespace下名为 cw 的ingress
func GetAllCWIngress(clientSet *kubernetes.Clientset) *v1.IngressList {
	// 使用 FieldSelector 来筛选名为 cw 的 Ingress
	ingresses, err := clientSet.NetworkingV1().Ingresses("").List(context.Background(), metav1.ListOptions{
		FieldSelector: "metadata.name=cw",
	})
	if err != nil {
		log.Fatalf("获取 Ingress 列表失败: %v", err)
	}

	return ingresses
}



// AddIPsToWhitelist 增加 IP 到 Ingress 注解的白名单
//func AddIPsToWhitelist(clientSet *kubernetes.Clientset, namespace, ingressName string, ips []string) error {
//	log.Println("正在增加ip到", namespace, "中的", ingressName, "IP为：", ips)
//	// 获取 Ingress 对象
//	ingress, err := clientSet.NetworkingV1().Ingresses(namespace).Get(context.Background(), ingressName, metav1.GetOptions{})
//	if err != nil {
//		return err
//	}
//
//	// 检查注解是否存在
//	whitelist, ok := ingress.Annotations["nginx.ingress.kubernetes.io/whitelist-source-range"]
//	if !ok {
//		log.Printf("注解 'nginx.ingress.kubernetes.io/whitelist-source-range' 不存在于 %s 的 Ingress %s，跳过添加操作n", namespace, ingressName)
//		return nil // 或者返回一个错误
//	}
//
//	// 获取当前白名单
//	ipsList := strings.Split(whitelist, ",")
//	ipsList = append(ipsList, ips...)
//	newWhitelist := strings.Join(ipsList, ",")
//
//	// 更新注解
//	ingress.Annotations["nginx.ingress.kubernetes.io/whitelist-source-range"] = newWhitelist
//
//	// 更新 Ingress 对象
//	_, err = clientSet.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
//	if err != nil {
//		return err
//	}
//
//	log.Println("加白完成:", namespace, "中的", ingressName, "IP为：", ips)
//
//	return nil
//}

func AddIPsToWhitelist(clientSet *kubernetes.Clientset, namespace, ingressName string, ips []string) error {
	log.Println("正在增加ip到", namespace, "中的", ingressName, "IP为：", ips)
	// 获取 Ingress 对象
	ingress, err := clientSet.NetworkingV1().Ingresses(namespace).Get(context.Background(), ingressName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// 删除现有的白名单
	ingress.Annotations["nginx.ingress.kubernetes.io/whitelist-source-range"] = ""

	// 添加新的 IP 到白名单
	newWhitelist := strings.Join(ips, ",")
	ingress.Annotations["nginx.ingress.kubernetes.io/whitelist-source-range"] = newWhitelist

	// 更新 Ingress 对象
	_, err = clientSet.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	log.Println("加白完成:", namespace, "中的", ingressName, "IP为：", ips)

	return nil
}

// RemoveIPsFromWhitelist 从 Ingress 注解的白名单中删除 IP
//func RemoveIPsFromWhitelist(clientSet *kubernetes.Clientset, namespace, ingressName string, ipsToRemove []string) error {
//	//log.Println("正在删除", namespace, "中的", ingressName, "IP为：", ipsToRemove)
//	// 获取 Ingress 对象
//	ingress, err := clientSet.NetworkingV1().Ingresses(namespace).Get(context.Background(), ingressName, metav1.GetOptions{})
//	if err != nil {
//		return err
//	}
//
//	// 获取当前白名单
//	//whitelist := ingress.Annotations["nginx.ingress.kubernetes.io/whitelist-source-range"]
//
//	// 获取当前白名单 并检查注解是否存在
//	whitelist, ok := ingress.Annotations["nginx.ingress.kubernetes.io/whitelist-source-range"]
//	if !ok {
//		log.Printf("注解 'nginx.ingress.kubernetes.io/whitelist-source-range' 不存在于 %s 的 Ingress %s，跳过删除操作n", namespace, ingressName)
//		return nil // 或者返回一个错误
//	}
//
//	ipsList := strings.Split(whitelist, ",")
//	var newIPsList []string
//
//	for _, existingIP := range ipsList {
//		shouldRemove := false
//		for _, ipToRemove := range ipsToRemove {
//			if strings.TrimSpace(existingIP) == strings.TrimSpace(ipToRemove) {
//				shouldRemove = true
//				break
//			}
//		}
//		if !shouldRemove {
//			newIPsList = append(newIPsList, strings.TrimSpace(existingIP))
//		}
//
//	}
//
//	// 更新白名单
//	newWhitelist := strings.Join(newIPsList, ",")
//	ingress.Annotations["nginx.ingress.kubernetes.io/whitelist-source-range"] = newWhitelist
//
//	// 更新 Ingress 对象
//	_, err = clientSet.NetworkingV1().Ingresses(namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
