module github.com/karbonn/karna

require (
	github.com/aws/aws-sdk-go v1.25.11
	github.com/aws/aws-sdk-go-v2 v0.14.0
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/frankban/quicktest v1.5.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/graphql-go/graphql v0.7.8
	github.com/johnnadratowski/golang-neo4j-bolt-driver v0.0.0-20181101021923-6b24c0085aae

	github.com/logrusorgru/aurora v0.0.0-20190803045625-94edacc10f9b
	github.com/mholt/archiver v3.1.1+incompatible
	github.com/nwaples/rardecode v1.0.0 // indirect
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/spf13/cobra v0.0.5
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
)

replace github.com/karbonn/karna/cmd v0.0.1 => ./../github.com/karna/cmd

go 1.13
