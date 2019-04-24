package route53

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	sdkaws "github.com/blend/go-sdk/aws"
	"github.com/blend/go-sdk/ex"
)

// AddDNSEntry adds the dns entry to route53
func AddDNSEntry(fqdn, target string, ttl int64, cfg sdkaws.Config) error {
	svc := route53.New(sdkaws.MustNewSession(cfg))
	zone, err := getHostedZoneFromAWS(fqdn, svc)
	if err != nil {
		return ex.New(err)
	}
	rs := getRecordSetForCNAME(fqdn, target, ttl)
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  getChangeBatchForRecordSet(route53.ChangeActionUpsert, rs),
		HostedZoneId: zone.Id,
	}
	_, err = svc.ChangeResourceRecordSets(input)
	return ex.New(err)
}

func getHostedZoneFromAWS(fqdn string, svc *route53.Route53) (*route53.HostedZone, error) {
	zones, err := listHostedZones(svc)
	if err != nil {
		return nil, ex.New(err)
	}
	return hostedZoneForFQDN(fqdn, zones)
}

func hostedZoneForFQDN(fqdn string, zones []*route53.HostedZone) (*route53.HostedZone, error) {
	var bestMatch *route53.HostedZone
	for _, zone := range zones {
		if zone != nil && zone.Name != nil {
			name := strings.TrimSuffix(*zone.Name, ".") //remove trailing periods
			// check to see if this zone is a parent domain (suffix of) the fqdn, and then if it is a longer match
			// than the current best matching domain
			if strings.HasSuffix(fqdn, name) && (bestMatch == nil || len(name) > len(*bestMatch.Name)) {
				bestMatch = zone
			}
		}
	}
	if bestMatch != nil && aws.StringValue(bestMatch.Name) == "" {
		bestMatch = nil
	}
	return bestMatch, nil
}

func listHostedZones(svc *route53.Route53) ([]*route53.HostedZone, error) {
	input := route53.ListHostedZonesInput{}
	output, err := svc.ListHostedZones(&input)
	if err != nil {
		return nil, ex.New(err)
	}
	zones := output.HostedZones
	for output.IsTruncated != nil && *output.IsTruncated {
		input.Marker = output.NextMarker
		output, err = svc.ListHostedZones(&input)
		if err != nil {
			return nil, ex.New(err)
		}
		zones = append(zones, output.HostedZones...)
	}
	return zones, nil
}

func getRecordSetForCNAME(fqdn, elb string, ttl int64) *route53.ResourceRecordSet {
	records := []*route53.ResourceRecord{
		&route53.ResourceRecord{
			Value: aws.String(elb),
		},
	}
	return &route53.ResourceRecordSet{
		Type:            aws.String(route53.RRTypeCname),
		TTL:             aws.Int64(ttl),
		Name:            aws.String(fqdn),
		ResourceRecords: records,
	}
}

func getChangeBatchForRecordSet(action string, rs *route53.ResourceRecordSet) *route53.ChangeBatch {
	return &route53.ChangeBatch{
		Changes: []*route53.Change{
			&route53.Change{
				Action:            &action,
				ResourceRecordSet: rs,
			},
		},
	}
}
