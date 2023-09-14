package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func (s *StartConfig) replaceDNS(hostedZoneID string, recordName *types.ResourceRecordSet) error {

	R53Client := route53.NewFromConfig(s.AWSConfig)

	change := types.Change{
		Action:            types.ChangeActionUpsert,
		ResourceRecordSet: recordName,
	}
	_, err := R53Client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(hostedZoneID),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{change},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
