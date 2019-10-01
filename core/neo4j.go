package core

import (
	bolt "github.com/johnnadratowski/golang-neo4j-bolt-driver"
)

const (
	protocol = "bolt"
	username = "neo4j"
	password = ""
	host     = "localhost"
	port     = "7687"
)

func parseConfiguration(configuration *KarnaNeo4JConfiguration) {
	if configuration.Host == "" {
		configuration.Host = host
	}
	if configuration.Password == "" {
		configuration.Password = password
	}
	if configuration.Port == "" {
		configuration.Port = port
	}
	if configuration.Username == "" {
		configuration.Username = username
	}
}

func (neo4j *KarnaNeo4J) createConnection() bolt.Conn {
	driver := bolt.NewDriver()

	parseConfiguration(&neo4j.Configuration)

	con, err := driver.OpenNeo(protocol + "://" + neo4j.Configuration.Username + ":" + neo4j.Configuration.Password + "@" + neo4j.Configuration.Host + ":" + neo4j.Configuration.Port)
	handleError(err)
	return con
}

//Bulk => Bulk import in Neo4j database.
func (neo4j *KarnaNeo4J) Bulk(queries []string, args []map[string]interface{}) {
	conn := neo4j.createConnection()
	defer conn.Close()

	pipeline, err := conn.PreparePipeline(queries...)

	if err != nil {
		LogErrorMessage(err.Error())

	}

	_, err = pipeline.ExecPipeline(args...)

	if err != nil {
		LogErrorMessage(err.Error())

	}
}

//CleanUp => Remove all database entities.
func (neo4j *KarnaNeo4J) CleanUp() {
	conn := neo4j.createConnection()
	defer conn.Close()

	stmt, err := conn.PrepareNeo("MATCH (n) DETACH DELETE n")

	if err != nil {
		LogErrorMessage(err.Error())

	}

	_, err = stmt.QueryNeo(nil)

	if err != nil {
		LogErrorMessage(err.Error())

	}
}

func handleError(err error) {
	if err != nil {
		LogErrorMessage(err.Error())

	}
}
