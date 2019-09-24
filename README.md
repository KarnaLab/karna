# Karna

Karna is a set of packages that allow, separately from the manager, to visualize and build visualization tools for your lambda architecture.

## Installation

    go get -u github/karbonn/karna
    cd github/karna
    go install

## Modules

Karna is composed of three packages. Each package has a README file that describes in more detail the possibilities offered by Karna.

### Viz

This package allows to visualize in Neo4J, the architecture of your lambdas. Karna Viz exposes two commands:

- show: Build the tree
- cleanup: Destroy the tree

### Deploy

This package allows you to pack, deploy, tag (alias), and clean the selected lambda function, based on a configuration file (karna.json). Karna Deploy exposes one command:

- deploy -a (alias) -t (functionName)

### API

This package allows you to mount a GraphQL API that exposes the following AWS resources in JSON format: Lambda, APIGateway, EC2.

- start api

## Examples

An example karna.json configuration file is presented in the examples / deploy folder
