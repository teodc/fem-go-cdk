package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type FemGoCdkStackProps struct {
	awscdk.StackProps
}

func NewFemGoCdkStack(scope constructs.Construct, id string, props *FemGoCdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create the DynamoDB table
	dynamoDBTable := awsdynamodb.NewTable(stack, jsii.String("FemGoCdkTable"), &awsdynamodb.TableProps{
		TableName: jsii.String("FemGoCdkUsers"),
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),
			Type: awsdynamodb.AttributeType_STRING,
		},
	})

	// Create the Lambda function
	lambdaFunction := awslambda.NewFunction(stack, jsii.String("FemGoCdkFunction"), &awslambda.FunctionProps{
		FunctionName: jsii.String("FemGoCdkManageUsers"),
		Runtime:      awslambda.Runtime_PROVIDED_AL2023(),
		Code:         awslambda.AssetCode_FromAsset(jsii.String("lambda/func.zip"), nil),
		Handler:      jsii.String("main"),
	})

	// Allow the Lambda function to read/write to the DynamoDB table
	dynamoDBTable.GrantReadWriteData(lambdaFunction)

	// Create the API Gateway REST API
	api := awsapigateway.NewRestApi(stack, jsii.String("FemGoCdkRestApi"), &awsapigateway.RestApiProps{
		RestApiName: jsii.String("FemGoCdkUserApi"),
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
	})

	// Integrate Lambda with API Gateway
	lambdaIntegration := awsapigateway.NewLambdaIntegration(lambdaFunction, nil)

	// Create the API Gateway resources & methods
	usersResource := api.Root().AddResource(jsii.String("users"), nil)
	// - /users/register
	usersRegisterResource := usersResource.AddResource(jsii.String("register"), nil)
	usersRegisterResource.AddMethod(jsii.String("POST"), lambdaIntegration, nil)
	// - /users/login
	usersLoginResource := usersResource.AddResource(jsii.String("login"), nil)
	usersLoginResource.AddMethod(jsii.String("POST"), lambdaIntegration, nil)
	// - /users/protected
	usersProtectedResource := usersResource.AddResource(jsii.String("protected"), nil)
	usersProtectedResource.AddMethod(jsii.String("GET"), lambdaIntegration, nil)

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewFemGoCdkStack(app, "FemGoCdkStack", &FemGoCdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
