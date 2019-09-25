# Karna Viz

## How it works

Karna Viz will request the AWS SDK to build a graph in Neo4J composed of lambdas functions, their triggers, their endpoints
in APIGateway and any VPCS to which they belong.

A docker-compose file is available in the examples / viz folder.

## Commands

`karna viz show`

This command will build the tree in Neo4J.

### Options

You can set the following args:

`--port [-p]` => Neo4J port.

`--credentials [-c]` => Neo4J credentials (username/password).

`--host` => Neo4J host.

`karna viz cleanup`

This command destroy in Neo4J.
