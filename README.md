# CloudForge

CloudForge is an enterprise-grade cloud provisioning and orchestration platform. Designed with an event-driven microservices architecture, it reliably orchestrates the lifecycle of complex cloud environments (Networking, Databases, Storage, Compute, and Secrets) ensuring high availability, eventual consistency, and complete fault tolerance.

Built on top of modern serverless and distributed systems patterns, CloudForge provides a scalable engine capable of handling high-throughput provisioning requests through robust SAGA choreographies.

## Core Capabilities

- **SAGA Orchestration**: Automated resource provisioning workflows with built-in state management, rollbacks, and compensating transactions.
- **Asynchronous Messaging**: Highly decoupled architecture utilizing message queues and event buses for reliable, non-blocking execution.
- **Resiliency & Fault Tolerance**: Native support for Dead Letter Queues (DLQ), automated retries with exponential backoff, and strict timeout controls.
- **Observability**: End-to-end distributed tracing and centralized logging for real-time monitoring of the provisioning lifecycle.
- **Infrastructure as Code**: Fully declarative deployment of the platform backbone via Terraform.

## Tech Stack

- **Backend**: Go (Golang) + Fiber
- **Cloud Infrastructure**: AWS (Lambda, API Gateway, DynamoDB, S3, EventBridge, SQS, SNS, SES)
- **IaC**: Terraform
- **Local Environment**: Docker & LocalStack (Full offline emulation)
- **CI/CD**: GitHub Actions
