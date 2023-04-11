package main

import (
	"fmt"
	// "os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	// "gopkg.in/yaml.v2"
)

// type Config struct {
// 	TableName string `yaml:"tableName"`
	
// }

func CreateDb(config *Config, sess *session.Session) {
	// //  Open the configuration file
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
     tableName := config.TableName

	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))
	svc := dynamodb.New(sess)

	attributeDefinitions := []*dynamodb.AttributeDefinition{
		{
			AttributeName: aws.String("ID"),
			AttributeType: aws.String("S"),
		},
		{
			AttributeName: aws.String("Name"),
			AttributeType: aws.String("S"),
		},
	}

	keySchema := []*dynamodb.KeySchemaElement{
		{
			AttributeName: aws.String("ID"),
			KeyType:       aws.String("HASH"),
		},
		{
			AttributeName: aws.String("Name"),
			KeyType:       aws.String("RANGE"),
		},
	}

	provisionedThroughput := &dynamodb.ProvisionedThroughput{
		ReadCapacityUnits:  aws.Int64(10),
		WriteCapacityUnits: aws.Int64(10),
	}

	err = MakeTable(svc, attributeDefinitions, keySchema, provisionedThroughput, &tableName)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err)
		return
	}

	fmt.Println("Created the table", tableName)
}

func MakeTable(svc dynamodbiface.DynamoDBAPI, attributeDefinitions []*dynamodb.AttributeDefinition, keySchema []*dynamodb.KeySchemaElement, provisionedThroughput *dynamodb.ProvisionedThroughput, tableName *string) error {
	_, err := svc.CreateTable(&dynamodb.CreateTableInput{
		AttributeDefinitions:  attributeDefinitions,
		KeySchema:             keySchema,
		ProvisionedThroughput: provisionedThroughput,
		TableName:             tableName,
	})
	return err
}
