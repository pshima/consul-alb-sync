package sync

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func GetTargetGroup(s string) (*elbv2.DescribeTargetGroupsOutput, error) {

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := elbv2.New(sess)

	params := &elbv2.DescribeTargetGroupsInput{
		Names: []*string{
			aws.String(s),
		},
	}

	resp, err := svc.DescribeTargetGroups(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetTargetGroupHealth(s string) (*elbv2.DescribeTargetHealthOutput, error) {

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	svc := elbv2.New(sess)

	params := &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(s),
	}

	resp, err := svc.DescribeTargetHealth(params)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetInstanceIDFromIP(s string) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}

	svc := ec2.New(sess)
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("private-ip-address"),
				Values: []*string{
					aws.String(s),
				},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		return "", err
	}
	if len(resp.Reservations) < 1 {
		return "", fmt.Errorf("Unable to find instance id for %s", s)
	}
	if len(resp.Reservations[0].Instances) < 1 {
		return "", fmt.Errorf("Unable to find instance id for %s", s)
	}
	return *resp.Reservations[0].Instances[0].InstanceId, nil
}

func RemoveFromTargetGroup(tgarn string, id string, port int64) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := elbv2.New(sess)

	params := &elbv2.DeregisterTargetsInput{
		TargetGroupArn: aws.String(tgarn),
		Targets: []*elbv2.TargetDescription{
			{
				Id:   aws.String(id),
				Port: aws.Int64(port),
			},
		},
	}
	_, err = svc.DeregisterTargets(params)
	if err != nil {
		return err
	}

	return nil
}

func AddToTargetGroup(tgarn string, id string, port int64) error {
	sess, err := session.NewSession()
	if err != nil {
		return err
	}

	svc := elbv2.New(sess)

	params := &elbv2.RegisterTargetsInput{
		TargetGroupArn: aws.String(tgarn),
		Targets: []*elbv2.TargetDescription{
			{
				Id:   aws.String(id),
				Port: aws.Int64(port),
			},
		},
	}
	_, err = svc.RegisterTargets(params)
	if err != nil {
		return err
	}

	return nil
}
