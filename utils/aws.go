package utils

import (
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func listInstances(profile string) []*ec2.Instance {
	cli := awsEc2Client(profile, "ap-northeast-1")
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}
	var instances []*ec2.Instance

	res, _ := cli.DescribeInstances(params)
	for _, v := range res.Reservations {
		if v != nil {
			for _, w := range v.Instances {
				instanceName := getInstanceName(w)
				if !regexp.MustCompile(`Vyos*`).MatchString(instanceName) {
					instances = append(instances, w)
				}
			}
		} else {
			break
		}
	}
	return instances
}

func awsEc2Client(profile string, region string) *ec2.EC2 {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}
	sess := session.New(&config)
	ec2Client := ec2.New(sess)
	return ec2Client
}

func getInstanceName(instance *ec2.Instance) (instanceName string) {
	for _, t := range instance.Tags {
		if *t.Key == "Name" {
			instanceName = *t.Value
		}
	}
	if instanceName == "" {
		instanceName = *instance.Tags[0].Value
	}
	return
}
