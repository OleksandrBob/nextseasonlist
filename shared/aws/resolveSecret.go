package aws

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func ResolveSecret(name string) string {
	envSecret := os.Getenv(name)
	if envSecret != "" {
		return envSecret
	}

	fmt.Println(name + " is unset. Trying to read from Secrets manager.")

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1")) // TODO: remove hardcode
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(name),
		VersionStage: aws.String("AWSCURRENT"),
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err.Error())
	}

	return *result.SecretString
}
