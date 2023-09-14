package main

import (
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/route53"
)

func (s *StartConfig) fetchR53HostedZones() ([]string, error) {
	R53Client := route53.NewFromConfig(s.AWSConfig)

	HZList, err := R53Client.ListHostedZones(ctx, &route53.ListHostedZonesInput{})
	if err != nil {
		return nil, err
	}

	var hostedZones []string

	for _, hostedZone := range HZList.HostedZones {
		hostedZoneNumber := strings.Replace(*hostedZone.Id, "/hostedzone/", "", 1)
		hostedZones = append(hostedZones, hostedZoneNumber)
	}

	return hostedZones, nil
}
