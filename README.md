# Karna

Karna is a set of packages that allow you to manage and vizualise your lambda architecture.

## Installation

    go get -u github/karbonn/karna
    cd github/karna
    go install

### Deploy

This package allows you to pack, deploy, tag (alias), and clean the selected lambda function, based on a configuration file (karna.json). Karna Deploy exposes one command:

- deploy -a (alias) -t (functionName)

## Examples

All commands are avaiblable via: karna help or karna <subcommand> help
An example karna.json configuration file is presented in the examples / deploy folder
