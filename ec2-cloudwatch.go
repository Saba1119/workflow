package main

import (
	"fmt"
	// "os"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	// "gopkg.in/yaml.v2"
)

// type Config struct {
// 	ImageID            string `yaml:"IMAGE_ID"`
// 	InstanceType       string `yaml:"INSTANCE_TYPE"`
// 	GroupName          string `yaml:"GROUP_NAME"`
// 	AlarmName          string `yaml:"ALARM_NAME"`
// 	SecurityGroupName  string `yaml:"SECURITYGROUP_NAME"`
// 	VPCID              string `yaml:"VPC_ID"`
// 	TargetGroupName1   string `yaml:"TARGETGROUP_NAME1"`
// 	TargetGroupName2   string `yaml:"TARGETGROUP_NAME2"`
// 	ALBName            string `yaml:"ALB_NAME"`
// 	AWS_REGION         string `yaml:"AWS_REGION"`
// 	SecurityGroup1FROMPort int64  `yaml:"Sgfrom_PORT1"`
// 	SecurityGroup1TOPort int64  `yaml:"Sgto_PORT1"`
// 	SecurityGroup2FROMPort int64  `yaml:"Sgfrom_PORT2"`
// 	SecurityGroup2TOPort int64  `yaml:"Sgto_PORT2"`
// 	TargetGroup1PORT     int64  `yaml:"TG1PORT"`
// 	TargetGroup2PORT     int64  `yaml:"TG2PORT"`
// 	Listener1PORT        int64  `yaml:"LISTENER1PORT"`
// 	Listener2PORT        int64  `yaml:"LISTENER2PORT"`


// 	             }


func Createec2(config *Config, sess *session.Session) {

	// // Open the configuration file
    // configFile, err := os.Open("config.yaml")
    // if err != nil {
    //     fmt.Println("Error opening configuration file:", err)
    //     os.Exit(1)
    // }
    // defer configFile.Close()

    // // Parse the configuration file
    // config := Config{}
    // err = yaml.NewDecoder(configFile).Decode(&config)
    // if err != nil {
    //     fmt.Println("Error parsing configuration file:", err)
    //     os.Exit(1)
    // }

    var err error // Declare the err variable

	// Create a new session using the default AWS configuration
	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))

	// Create an EC2 service client
	svc := ec2.New(sess)


	// Specify the parameters for the new EC2 instance
	params := &ec2.RunInstancesInput{
		ImageId:      aws.String(config.ImageID),
		InstanceType: aws.String(config.InstanceType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
        {
            DeviceName: aws.String("/dev/sdf"),
            Ebs: &ec2.EbsBlockDevice{
                VolumeSize: aws.Int64(10), // Size of the volume in GB
                VolumeType: aws.String("gp2"), // Type of the volume
                DeleteOnTermination: aws.Bool(true), // Automatically delete the volume when the instance is terminated
            },
        },
    },

}

    // Read environment variables
	groupName := config.GroupName 
	alarmName := config.AlarmName
	// instanceID := os.Getenv("INSTANCE_ID")
	// Create a new session to interact with AWS
	sess, err = session.NewSession(&aws.Config{
		Region: aws.String(config.AWS_REGION), // Replace with your desired region
	})
	if err != nil {
		fmt.Println("Error creating session: ", err)
		return }

	// Create the EC2 instance
	result, err := svc.RunInstances(params)
	if err != nil {
		fmt.Println("Error", err)
		return
	} 

	// Get the instance ID of the newly created EC2 instance
	instanceId := result.Instances[0].InstanceId
	// Specify the parameters for the new security group
		securityGroupDescription := "My security group description"

	// Create the security group
	createSecurityGroupResult, err := svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(config.SecurityGroupName ),
		Description: aws.String(securityGroupDescription),
		VpcId:       aws.String(config.VPCID),
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	// Get the ID of the newly created security group
	securityGroupId := createSecurityGroupResult.GroupId

	// Open port 8000 in the new security group
	_, err = svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:    securityGroupId,
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(config.SecurityGroup1FROMPort),
		ToPort:     aws.Int64(config.SecurityGroup1TOPort),
		CidrIp:     aws.String("0.0.0.0/0"),
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	// Open port 8080 in the new security group
	_, err = svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:    securityGroupId,
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(config.SecurityGroup2FROMPort),
		ToPort:     aws.Int64(config.SecurityGroup2TOPort),
		CidrIp:     aws.String("0.0.0.0/0"),
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	// Attach the security group to the instance
	_, err = svc.ModifyInstanceAttribute(&ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(*instanceId),
		Groups:     []*string{aws.String(*securityGroupId)},
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}
	// Add tags to the instance
	_, err = svc.CreateTags(&ec2.CreateTagsInput{
		Resources: []*string{aws.String(*instanceId)},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String("My EC2 Instance-sdk"),
			},
			
		},
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	} 


	// Create an EC2 service client
	ec2svc := ec2.New(sess)
	
	// Wait for the instance to reach the running state
	err = ec2svc.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(*instanceId)},
	})
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	// Create an ELBV2 service client
	elbv2Svc := elbv2.New(sess)

	// Create a target group1 for the ALB
	createTGOutput, err := elbv2Svc.CreateTargetGroup(&elbv2.CreateTargetGroupInput{
		Name:        aws.String(config.TargetGroupName1),
		Protocol:    aws.String("HTTP"),
		Port:        aws.Int64(config.TargetGroup1PORT),
		VpcId:       aws.String(config.VPCID), // replace with your own VPC ID
		TargetType:  aws.String("instance"),
		HealthCheckProtocol: aws.String("HTTP"),
		HealthCheckPath:     aws.String("/healthcheck"),
		HealthCheckIntervalSeconds: aws.Int64(30),
		HealthyThresholdCount:      aws.Int64(2),
		UnhealthyThresholdCount:    aws.Int64(2),
	})
	if err != nil {
		fmt.Println("Error creating target group:", err)
		return
	}
	tgArn1 := createTGOutput.TargetGroups[0].TargetGroupArn
	fmt.Println("Target group1 created successfully")
	fmt.Println("ARN:", *tgArn1)


	// Create a target group2 for the ALB
	createTGOutput, err = elbv2Svc.CreateTargetGroup(&elbv2.CreateTargetGroupInput{
		Name:        aws.String(config.TargetGroupName2),
		Protocol:    aws.String("HTTP"),
		Port:        aws.Int64(config.TargetGroup2PORT),
		VpcId:       aws.String(config.VPCID), // replace with your own VPC ID
		TargetType:  aws.String("instance"),
		HealthCheckProtocol: aws.String("HTTP"),
		HealthCheckPath:     aws.String("/healthcheck"),
		HealthCheckIntervalSeconds: aws.Int64(30),
		HealthyThresholdCount:      aws.Int64(2),
		UnhealthyThresholdCount:    aws.Int64(2),
	})
	if err != nil {
		fmt.Println("Error creating target group:", err)
		return
	}
	tgArn2 := createTGOutput.TargetGroups[0].TargetGroupArn
	fmt.Println("Target group2 created successfully")
	fmt.Println("ARN:", *tgArn2)


    // create the first subnet
    subnet1, err := svc.CreateSubnet(&ec2.CreateSubnetInput{
        CidrBlock: aws.String("172.31.128.0/20"),
        VpcId:     aws.String(config.VPCID),
        AvailabilityZone: aws.String("us-east-1b"),
    })
    if err != nil {
        fmt.Println("Error creating subnet 1:", err)
        return
    }
    fmt.Println("Subnet 1 created with ID:", *subnet1.Subnet.SubnetId)

    // create the second subnet
    subnet2, err := svc.CreateSubnet(&ec2.CreateSubnetInput{
        CidrBlock: aws.String("172.31.96.0/20"),
        VpcId:     aws.String(config.VPCID),
        AvailabilityZone: aws.String("us-east-1c"),
    })
    if err != nil {
        fmt.Println("Error creating subnet 2:", err)
        return
    }
    fmt.Println("Subnet 2 created with ID:", *subnet2.Subnet.SubnetId)
    Subnet1 := subnet1.Subnet.SubnetId
    Subnet2 := subnet2.Subnet.SubnetId

	// Create a new ALB
	createLBOutput, err := elbv2Svc.CreateLoadBalancer(&elbv2.CreateLoadBalancerInput{
		Name:           aws.String(config.ALBName),
		Subnets:        []*string{aws.String(*Subnet1), aws.String(*Subnet2)}, // replace with your own subnet IDs
		SecurityGroups: []*string{aws.String(*securityGroupId)}, // replace with your own security group IDs
		IpAddressType:  aws.String("ipv4"),
	})
	if err != nil {
		fmt.Println("Error creating ALB:", err)
		return
	}
	lbArn := createLBOutput.LoadBalancers[0].LoadBalancerArn
	lbDns := createLBOutput.LoadBalancers[0].DNSName
	fmt.Println("ALB created successfully")
	fmt.Println("ARN:", *lbArn)
	fmt.Println("DNS:", *lbDns)

	// Register the EC2 instance with the target group1
	_, err = elbv2Svc.RegisterTargets(&elbv2.RegisterTargetsInput{
		TargetGroupArn: tgArn1,
		Targets: []*elbv2.TargetDescription{
			{
				Id: aws.String(*instanceId),
			},
		},
	})
	if err != nil {
		fmt.Println("Error registering target:", err)
		return
	}
	fmt.Println("Target registered successfully")


	// Register the EC2 instance with the target group2
	_, err = elbv2Svc.RegisterTargets(&elbv2.RegisterTargetsInput{
		TargetGroupArn: tgArn2,
		Targets: []*elbv2.TargetDescription{
			{
				Id: aws.String(*instanceId),
			},
		},
	})
	if err != nil {
		fmt.Println("Error registering target:", err)
		return
	}
	fmt.Println("Target registered successfully")


	// Create a listener1 for the ALB
	_, err = elbv2Svc.CreateListener(&elbv2.CreateListenerInput{
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("forward"),
				TargetGroupArn: tgArn1,
			},
		},
		LoadBalancerArn: lbArn,
		Protocol:        aws.String("HTTP"),
		Port:            aws.Int64(config.Listener1PORT),
	})
	if err != nil {
		fmt.Println("Error creating listener1:", err)
		return
	}
	fmt.Println("Listener1 created successfully")



	// Create a listener2 for the ALB
	_, err = elbv2Svc.CreateListener(&elbv2.CreateListenerInput{
		DefaultActions: []*elbv2.Action{
			{
				Type: aws.String("forward"),
				TargetGroupArn: tgArn2,
			},
		},
		LoadBalancerArn: lbArn,
		Protocol:        aws.String("HTTP"),
		Port:            aws.Int64(config.Listener2PORT),
	})
	if err != nil {
		fmt.Println("Error creating listener2:", err)
		return
	}
	fmt.Println("Listener2 created successfully")


	// Create a CloudWatch Logs client
	logsSvc := cloudwatchlogs.New(sess)

	// Create a new log group
	_, err = logsSvc.CreateLogGroup(&cloudwatchlogs.CreateLogGroupInput{
		LogGroupName: aws.String(groupName),
	})
	if err != nil {
		fmt.Println("Error creating log group: ", err)
		return
	}

	// Define a filter to monitor logs for the word "error" or "exception"
	filterName := "my-log-filter"
	filterPattern := "error"
	_, err = logsSvc.PutMetricFilter(&cloudwatchlogs.PutMetricFilterInput{
		LogGroupName:  aws.String(groupName),
		FilterName:    aws.String(filterName),
		FilterPattern: aws.String(filterPattern),
		MetricTransformations: []*cloudwatchlogs.MetricTransformation{
			{
				MetricName:      aws.String("ErrorCount"),
				MetricNamespace: aws.String("MyApplication"),
				MetricValue:     aws.String("1"),
			},
		},
	})
	if err != nil {
		fmt.Println("Error creating log filter: ", err)
		return
	}

	// Create a CloudWatch client
	cwSvc := cloudwatch.New(sess)

	// Define an alarm to trigger when CPU usage reaches 5% on an EC2 instance
	metricName := "CPUUtilization"
	namespace := "AWS/EC2"
	_, err = cwSvc.PutMetricAlarm(&cloudwatch.PutMetricAlarmInput{
		AlarmName:          aws.String(alarmName),
		ComparisonOperator: aws.String("GreaterThanThreshold"),
		EvaluationPeriods:  aws.Int64(1),
		MetricName:         aws.String(metricName),
		Namespace:          aws.String(namespace),
		Period:             aws.Int64(60),
		Statistic:          aws.String("Average"),
		Threshold:          aws.Float64(5.0),
		ActionsEnabled:     aws.Bool(true),
		AlarmActions: []*string{
			aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:pipeline", aws.StringValue(sess.Config.Region), "554248189203")),
		},
		OKActions: []*string{
			aws.String(fmt.Sprintf("arn:aws:sns:%s:%s:pipeline", aws.StringValue(sess.Config.Region), "554248189203")),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(*instanceId),
			},
		},
	})
	if err != nil {
		fmt.Println("Error creating alarm:", err)
		return
	}

	fmt.Println("Successfully created CloudWatch log group, filter, and alarm!")
} 