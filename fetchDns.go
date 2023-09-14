package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (s *StartConfig) fetchAllDNS(hostedZoneID string) ([]types.ResourceRecordSet, error) {

	R53Client := route53.NewFromConfig(s.AWSConfig)

	var cnameRecords []types.ResourceRecordSet

	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
	}

	for {
		recSetOutput, err := R53Client.ListResourceRecordSets(ctx, input)
		if err != nil {
			return nil, err
		}

		for _, record := range recSetOutput.ResourceRecordSets {
			cnameRecords = append(cnameRecords, record)
		}

		if recSetOutput.IsTruncated {
			input.StartRecordName = recSetOutput.NextRecordName
		} else {
			break
		}
	}

	return cnameRecords, nil
}
