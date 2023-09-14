package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var ctx = context.TODO()

type StartConfig struct {
	AWSConfig aws.Config
	EKSClient *rest.Config
}

func startEKSConfig(currentContext string) *StartConfig {
	var kubeconfig string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfig = os.Getenv("KUBECFG_PATH")
	}

	apiConf, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	apiConf.CurrentContext = currentContext

	conf, err := clientcmd.NewDefaultClientConfig(*apiConf, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		log.Fatal(err)
	}

	return &StartConfig{
		EKSClient: conf,
	}
}

func startAWSConfig(profile, region string) *StartConfig {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(profile), config.WithDefaultRegion(region))
	if err != nil {
		log.Fatal(err)
	}

	return &StartConfig{
		AWSConfig: cfg,
	}
}

const (
	profile = "default"
	region  = "us-east-1"
)

func main() {

	clusterName := os.Args[1]

	RunAWS := startAWSConfig(profile, region)

	clusterARN, err := clusterARN(clusterName, RunAWS.AWSConfig)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(clusterARN)

	hostedZones, err := RunAWS.fetchR53HostedZones()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(hostedZones)

	RunEKSClient := startEKSConfig(clusterARN)

	keyValueIngress, err := RunEKSClient.fetchIngress()
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range keyValueIngress {
		fmt.Printf("%s - %s\n", key, value)
	}

	for _, hostedZone := range hostedZones {
		fmt.Printf("Iniciando processo na hostedZone: %s\n", hostedZone)
		recordSetItems, err := RunAWS.fetchAllDNS(hostedZone)
		if err != nil {
			log.Fatal(err)
		}
		for _, recordSetItem := range recordSetItems {
			value, ok := keyValueIngress[*recordSetItem.Name]
			if ok {
				fmt.Printf("Substituindo: %v - %v para -> %v - %v\n", *recordSetItem.Name, *recordSetItem.ResourceRecords[0].Value, *recordSetItem.Name, value)
				recordSetItem.ResourceRecords[0].Value = aws.String(value)
				err = RunAWS.replaceDNS(hostedZone, &recordSetItem)
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("Endereço: %s não foi encontrado na lista de endereços do Cluster: %s\n", *recordSetItem.Name, hostedZone)
			}
		}
	}
}
