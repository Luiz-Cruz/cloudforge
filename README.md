# cloudforge ☁️


An enterprise-grade cloud provisioning and orchestration platform.

## Features 💻

- Automated Resource Provisioning
  - Networking
  - Databases
  - Storage
  - Compute & Secrets
- SAGA Orchestration
  - Rollbacks and compensating transactions
- Asynchronous Messaging
  - Dead Letter Queues (DLQ)
  - Exponential backoff retries
- Observability
  - Distributed tracing (AWS X-Ray)
  - Centralized logging

## Libraries ⚙️

- [Go Fiber](https://github.com/gofiber/fiber)
- [AWS SDK v2](https://github.com/aws/aws-sdk-go-v2)
- [aws-lambda-go-api-proxy](https://github.com/awslabs/aws-lambda-go-api-proxy)
- [Viper](https://github.com/spf13/viper)
- [LocalStack](https://github.com/localstack/localstack)

## Deploy ✈️

The application was built upon a Docker image but relies mostly on AWS resources to work. To deploy on AWS, simply configure **Terraform** with some vars described below, run the script, and it's all done!

### AWS ☁️

The following resources are used on AWS:

- **DynamoDB** to store provisioning states and metadata
- **S3** to store configurations and generated asset files
- **EventBridge** to CRON the orchestration jobs
- **SQS** to handle queues and DLQs for fault tolerance
- **SNS** to publish state changes and alarms
- **SES** to send transactional emails
- **Step Functions** to orchestrate the SAGA patterns
- **Lambda** to run the serverless application (both API and JOB)
- **API Gateway** to provide a RESTful interface

### Configuration 🛠

The following configuration are required through **Terraform vars**

|Terraform var|Environment variable|Description|
|---|---|---|
|Hard coded on Terraform|SERVER|Used to define the environment where the application will run. Defaults to **AWS**|
|Hard coded on Terraform|APPLICATION|Used to define the Lambda type: **API** or **JOB**|
|Hard coded on Terraform|CLOUD|Used to map clients. Values: **AWS** or **LOCAL**|
|Terraform takes it from S3 resource|STORAGE|The S3 bucket|
|Terraform takes it from SNS resource|REPORTS_TOPIC|The topic to notify alarms and events|

## Support ✉️

You can create a PR for it =)
