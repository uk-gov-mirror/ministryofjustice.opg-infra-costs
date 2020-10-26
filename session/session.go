package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Assumed creates a session using the ARN and generates an authenticated session based on assuming the role
func Assumed(roleArn string, region string) (*session.Session, error) {
	baseSess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	stsSvc := sts.New(baseSess)
	sessionName := "cost-cli-session"
	assumedRole, err := stsSvc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
	})
	if err != nil {
		return nil, err
	}

	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(
			*assumedRole.Credentials.AccessKeyId,
			*assumedRole.Credentials.SecretAccessKey,
			*assumedRole.Credentials.SessionToken),
		Region: aws.String(region),
	})
}

// Identity returns a caller identity object based on the account and role string passed in
func Identity(roleArn string, region string) (*sts.GetCallerIdentityOutput, error) {
	sess, err := Assumed(roleArn, region)
	if err != nil {
		return nil, err
	}
	stsClient := sts.New(sess)
	request := sts.GetCallerIdentityInput{}
	return stsClient.GetCallerIdentity(&request)
}
