# Karna API

## How it works

Karna API is an GraphQL API. It will request the AWS SDK and exposes its response under JSON format.

The karna API is available on `http://localhost:8000/graphql?query={...}`

## Endpoints

Each models are available in `ìnternal/api/models.go`.

- lambda: Return all Lambdas on format: https://github.com/karbonn/karna/blob/11331c5f9e32b1931f86781c20a3878e22eda5b8/internal/api/models.go#L76
- apigateway: Return all APIGateway RestAPIS on format: https://github.com/karbonn/karna/blob/11331c5f9e32b1931f86781c20a3878e22eda5b8/internal/api/models.go#L99
- ec2: Return all EC2 and VPCS on format: https://github.com/karbonn/karna/blob/11331c5f9e32b1931f86781c20a3878e22eda5b8/internal/api/models.go#L188

## Examples:

- Get all lambdas and their Layers:
  `http://localhost:8000/graphql?query={lambda{Layers{Arn} }}`

- Get all apigateway rest APIS, their names and their associated stages:
  `http://localhost:8000/graphql?query={apigateway{ API{Name}, Stages{Name,Stage} }}`
