package main

import (
	"encoding/json"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Retornado chave=Host valor=Address
func (s *StartConfig) fetchIngress() (map[string]string, error) {

	kclient, err := kubernetes.NewForConfig(s.EKSClient)
	if err != nil {
		return nil, err
	}

	allNamespaces, err := kclient.CoreV1().Namespaces().List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var mapIngress = make(map[string]string)

	for _, ns := range allNamespaces.Items {
		ingressList, err := kclient.NetworkingV1().Ingresses(ns.Name).List(ctx, v1.ListOptions{})
		if err != nil {
			return nil, err
		}

		jsonR, err := json.Marshal(ingressList)
		if err != nil {
			return nil, err
		}

		var data Ingress
		err = json.Unmarshal(jsonR, &data)
		if err != nil {
			return nil, err
		}

		for _, items := range data.Items {
			for _, rule := range items.Spec.Rules {
				if rule.Host != "" {
					//fmt.Printf("Host: %s\nAddress: %s\n", rule.Host, items.Status.LoadBalancer.Ingress[0].Hostname)
					mapIngress[rule.Host+"."] = items.Status.LoadBalancer.Ingress[0].Hostname
				}
			}
		}
	}

	return mapIngress, nil
}
