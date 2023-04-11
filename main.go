package main

import (
	"fmt"
	"os"

	 // "github.com/aws/aws-sdk-go/aws"
	 "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/dynamodb"
	// "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"gopkg.in/yaml.v2"
     )

type Config struct {
	TableName string `yaml:"tableName"`
	AWS_REGION string `yaml:"AWS_REGION"`
    BUCKET_NAME string `yaml:"BUCKET_NAME"`
    KMS_KEY string `yaml:"KMS_KEY"`
    ImageID            string `yaml:"IMAGE_ID"`
	InstanceType       string `yaml:"INSTANCE_TYPE"`
	GroupName          string `yaml:"GROUP_NAME"`
	AlarmName          string `yaml:"ALARM_NAME"`
	SecurityGroupName  string `yaml:"SECURITYGROUP_NAME"`
	VPCID              string `yaml:"VPC_ID"`
	TargetGroupName1   string `yaml:"TARGETGROUP_NAME1"`
	TargetGroupName2   string `yaml:"TARGETGROUP_NAME2"`
	ALBName            string `yaml:"ALB_NAME"`
	SecurityGroup1FROMPort int64  `yaml:"Sgfrom_PORT1"`
	SecurityGroup1TOPort int64  `yaml:"Sgto_PORT1"`
	SecurityGroup2FROMPort int64  `yaml:"Sgfrom_PORT2"`
	SecurityGroup2TOPort int64  `yaml:"Sgto_PORT2"`
	TargetGroup1PORT     int64  `yaml:"TG1PORT"`
	TargetGroup2PORT     int64  `yaml:"TG2PORT"`
	Listener1PORT        int64  `yaml:"LISTENER1PORT"`
	Listener2PORT        int64  `yaml:"LISTENER2PORT"`


	InstanceID          string `yaml:"INSTANCE_ID"`
    TargetGroupARN1   string  `yaml:"TARGETGROUP_ARN1"`
    TargetGroupARN2   string  `yaml:"TARGETGROUP_ARN2"`
    // KMSID            string `yaml:"KMS_ID"`

}

func main()  {
	//  Open the configuration file
    configFile, err := os.Open("config.yaml")
    if err != nil {
        fmt.Println("Error opening configuration file:", err)
        os.Exit(1)
    }
    defer configFile.Close()

    // Parse the configuration file
    config := Config{}
    err = yaml.NewDecoder(configFile).Decode(&config)
    if err != nil {
        fmt.Println("Error parsing configuration file:", err)
        os.Exit(1)
    }   
    sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
    if err != nil {
        fmt.Println("Error creating session:", err)
        os.Exit(1)
     }
	// CreateS3(&config, sess)
	// Createec2(&config, sess)
	// CreateDb(&config, sess)	
	teardown(&config, sess)
}
