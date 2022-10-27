package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() {
	var err error
	// TEMPLATE "postgres://username:password@localhost:5432/database"
	databaseUrl := "postgres://postgres:admin@localhost:5432/personal_web"

	Conn, err = pgx.Connect(context.Background(), databaseUrl)

	// ERROR HANDLING
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database %v\n", err)
		os.Exit(1)
	} else {
		fmt.Println("Success connect to database")
	}
}
