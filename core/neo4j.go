package core

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

const (
	protocol = "bolt"
	username = "neo4j"
	password = "password"
	host     = "localhost"
	port     = "7687"
)

func createConnection() bolt.Conn {
	driver := bolt.NewDriver()
	con, err := driver.OpenNeo(protocol + "://" + username + ":" + password + "@" + host + ":" + port)
	handleError(err)
	return con
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

//Bulk => Bulk import in Neo4j database.
func Bulk(queries []string, args []map[string]interface{}) {
	conn := createConnection()
	defer conn.Close()

	pipeline, err := conn.PreparePipeline(queries...)

	if err != nil {
		panic(err.Error())
	}

	_, err = pipeline.ExecPipeline(args...)

	if err != nil {
		panic(err.Error())
	}
}

//CleanUp => Remove all database entities.
func CleanUp() {
	conn := createConnection()
	defer conn.Close()

	stmt, err := conn.PrepareNeo("MATCH (n) DETACH DELETE n")

	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.QueryNeo(nil)

	if err != nil {
		panic(err.Error())
	}
}
