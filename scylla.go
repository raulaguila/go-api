package main

import (
	"fmt"
	"log"

	"github.com/gocql/gocql"
)

func main() {
	// Configura a conexão com o ScyllaDB
	cluster := gocql.NewCluster("127.0.0.1") // IP do nó do ScyllaDB
	cluster.Port = 9042
	cluster.Keyspace = "meubanco"      // Nome do keyspace
	cluster.Consistency = gocql.Quorum // Define o nível de consistência
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: "",
		Password: "",
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Erro ao conectar no ScyllaDB: %v", err)
	}
	defer session.Close()

	// Criar uma tabela de usuários
	err = session.Query(`
		CREATE TABLE IF NOT EXISTS usuarios (
			id UUID PRIMARY KEY,
			nome TEXT,
			idade INT
		)
	`).Exec()
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	fmt.Println("Tabela 'usuarios' criada com sucesso!")

	// Inserir um usuário no ScyllaDB
	id := gocql.TimeUUID() // Gera um UUID único
	err = session.Query(`
		INSERT INTO usuarios (id, nome, idade) VALUES (?, ?, ?)
	`, id, "João Silva", 30).Exec()
	if err != nil {
		log.Fatalf("Erro ao inserir usuário: %v", err)
	}

	fmt.Println("Usuário inserido com sucesso!")

	// Consultar todos os usuários
	var userID gocql.UUID
	var nome string
	var idade int

	iter := session.Query("SELECT id, nome, idade FROM usuarios").Iter()
	for iter.Scan(&userID, &nome, &idade) {
		fmt.Printf("Usuário encontrado: ID=%v, Nome=%s, Idade=%d\n", userID, nome, idade)
	}

	if err := iter.Close(); err != nil {
		log.Fatalf("Erro ao iterar resultados: %v", err)
	}
}
