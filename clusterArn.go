package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
)

func clusterARN(clusterName string, cfg aws.Config) (string, error) {

	ClusterClient := eks.NewFromConfig(cfg)

	clusterOutput, err := ClusterClient.DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
	if err != nil {
		return "", err
	}
	return *clusterOutput.Cluster.Arn, nil
}
