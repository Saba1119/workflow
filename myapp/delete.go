package main

import (
    "fmt"
    "os"


    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/ec2"
    "github.com/aws/aws-sdk-go/service/s3"
    // "github.com/aws/aws-sdk-go/service/kms"
    "github.com/aws/aws-sdk-go/service/elbv2"
    "github.com/aws/aws-sdk-go/service/cloudwatch"
    "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
    // "gopkg.in/yaml.v2"
    // "io/ioutil"
)

func teardown(config *Config, sess *session.Session) {
    // Initialize AWS service clients
    ec2Svc := ec2.New(sess)
    elbv2Svc := elbv2.New(sess)
    cwSvc := cloudwatch.New(sess)
    cwlSvc := cloudwatchlogs.New(sess)

    // Delete EC2 instances
    _, err := ec2Svc.TerminateInstances(&ec2.TerminateInstancesInput{
        InstanceIds: []*string{aws.String(config.InstanceID)},
    })
    if err != nil {
        fmt.Println("Error deleting EC2 instances:", err)
        os.Exit(1)
    }
     // Wait until the EC2 instance is deleted.
    fmt.Println("Waiting for instance to be deleted...")
    err = ec2Svc.WaitUntilInstanceTerminated(&ec2.DescribeInstancesInput{
        InstanceIds: []*string{aws.String(config.InstanceID)},
    })
    if err != nil {
        fmt.Println("Error waiting for instance to be deleted:", err)
        return
    }

    fmt.Println("EC2 instances deleted successfully.")

    // Delete ALB resources
    albName := config.ALBName

    // Describe the ALB with the specified name to get its ARN
    describeInput := &elbv2.DescribeLoadBalancersInput{
        Names: []*string{aws.String(albName)},
    }
    describeOutput, err := elbv2Svc.DescribeLoadBalancers(describeInput)
    if err != nil {
        fmt.Println("Failed to describe ALB:", err)
        os.Exit(1)
    }

    // Check if the ALB exists
    if len(describeOutput.LoadBalancers) == 0 {
        fmt.Println("ALB not found with name:", albName)
        os.Exit(1)
    }

    // Delete the ALB using its ARN
    albArn := describeOutput.LoadBalancers[0].LoadBalancerArn
    deleteInput := &elbv2.DeleteLoadBalancerInput{
        LoadBalancerArn: albArn,
    }
    _, err = elbv2Svc.DeleteLoadBalancer(deleteInput)
    if err != nil {
        fmt.Println("Failed to delete ALB:", err)
        os.Exit(1)
    }

    // Wait for the ALB to be deleted
     err = elbv2Svc.WaitUntilLoadBalancersDeleted(&elbv2.DescribeLoadBalancersInput{
        LoadBalancerArns: []*string{albArn},
    })
    if err != nil {
        fmt.Println("Failed to wait for ALB deletion:", err)
        os.Exit(1)
    }

    fmt.Println("ALB deleted successfully!")

    
    // Create CloudWatch and CloudWatch Logs service clients
    // cwSvc := cloudwatch.New(sess)
    // cwlSvc := cloudwatchlogs.New(sess)

    // Delete the CloudWatch Alarm
    _, err = cwSvc.DeleteAlarms(&cloudwatch.DeleteAlarmsInput{
        AlarmNames: []*string{aws.String(config.AlarmName)},
    })
    if err != nil {
        fmt.Println("Error deleting CloudWatch Alarm:", err)
        os.Exit(1)
    }

    fmt.Println("CloudWatch Alarm deleted successfully.")

    // Delete the CloudWatch Logs Log Group
    _, err = cwlSvc.DeleteLogGroup(&cloudwatchlogs.DeleteLogGroupInput{
        LogGroupName: aws.String(config.GroupName),
    })
    if err != nil {
        fmt.Println("Error deleting CloudWatch Logs Log Group:", err)
        os.Exit(1)
    }

    fmt.Println("CloudWatch Logs Log Group deleted successfully.")
    // Delete S3 bucket and objects
    s3Svc := s3.New(sess)
    _, err = s3Svc.DeleteBucket(&s3.DeleteBucketInput{
        Bucket: aws.String(config.BUCKET_NAME ),
    })
    if err != nil {
        fmt.Println("Error deleting S3 bucket:", err)
        os.Exit(1)
    }
    fmt.Println("S3 bucket deleted successfully.")

    // Create a KMS service client
    // svc := kms.New(sess)

    // Specify the KMS key ID to delete
   

    // // Delete the KMS key
    // _, err = svc.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
    //     KeyId:               aws.String(keyId),
    // })
    // if err != nil {
    //     fmt.Println("Error deleting KMS key:", err)
    //     os.Exit(1)
    // }
    // fmt.Println("KMS key deleted successfully.")

    // Delete DynamoDB table
    dynamoSvc := dynamodb.New(sess)
    _, err = dynamoSvc.DeleteTable(&dynamodb.DeleteTableInput{
        TableName: aws.String(config.TableName),
    })
    if err != nil {
        fmt.Println("Error deleting DynamoDB table:", err)
        os.Exit(1)
    }
    fmt.Println("DynamoDB table deleted successfully.")

    // Delete the target groups
    _, err = elbv2Svc.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
        TargetGroupArn: aws.String(config.TargetGroupARN1), // Set the target group name
    })
    if err != nil {
        fmt.Println("Error deleting ALB target group:", err)
        os.Exit(1)
    }
    fmt.Println("ALB target group deleted successfully.")

    _, err = elbv2Svc.DeleteTargetGroup(&elbv2.DeleteTargetGroupInput{
        TargetGroupArn: aws.String(config.TargetGroupARN2), // Set the target group name
    })
    if err != nil {
        fmt.Println("Error deleting ALB target group2:", err)
        os.Exit(1)
    }
    fmt.Println("ALB target group2 deleted successfully.")

     // Delete the subnets
    _, err = ec2Svc.DeleteSubnet(&ec2.DeleteSubnetInput{
        SubnetId: aws.String("subnet-04bb6e9cb783d95e7"),
    })
    if err != nil {
        fmt.Println("Error deleting subnet:", err)
        os.Exit(1)
    }
    fmt.Println("Subnet1 deleted successfully.")

    _, err = ec2Svc.DeleteSubnet(&ec2.DeleteSubnetInput{
        SubnetId: aws.String("subnet-08d2dba7c4b2e35c2"),
    })
    if err != nil {
        fmt.Println("Error deleting subnets:", err)
        os.Exit(1)
    }
    fmt.Println("Subnet2 deleted successfully.")
}