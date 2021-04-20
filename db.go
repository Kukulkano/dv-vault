package main

// Q: Why are the passwords not encrypted/hashed?
// A: If the attacker gets the database with the passwords, he also
//    got the whole content. Thus, he also has all data.
//    By this, he no longer needs the passwords as he already has
//    all data. Therefore, we do not need to secure that here.
//    Also, consider it being much faster to not always calculating
//    bcrypt or sha2 hashed passwords for every request.

import (
	"fmt"
	"runtime"
	"time"

	"github.com/jackc/pgx"
)

var DB *pgx.ConnPool

func InitDatabase() bool {

	fmt.Print("Connect CockroachDB… ")

	// Set client connection
	var poolConfig pgx.ConnPoolConfig
	config, err := pgx.ParseConnectionString(cfg.ConnectionString)
	if err != nil {
		panic("Can not parse your connection string")
	}

	poolConfig.ConnConfig = config
	poolConfig.AcquireTimeout = time.Minute

	fmt.Print(config.Host + " ")

	maxConn := cfg.MaxConnections
	if maxConn < 1 {
		// <1 is auto, which uses CPU cores (incl. hyperthreading) * 3
		// https://www.cockroachlabs.com/docs/v21.1/connection-pooling.html#sizing-connection-pools
		// "Many workloads perform best when the number of connections was
		// between 2 and 4 times the number of CPU cores in the cluster."
		maxConn = runtime.NumCPU() * 3
	}

	poolConfig.MaxConnections = maxConn

	// Connect to CockroachDB
	DB, err = pgx.NewConnPool(poolConfig)
	if err != nil {
		fmt.Printf("\n%v\n", err)
		panic("Can not connect new pool to CockroachDB")
	}

	// Check the connection
	var w int
	row := DB.QueryRow("SELECT COUNT(*) FROM dv.provider").Scan(&w)
	if row != nil {
		panic("Test query to 'dv.providers' table failed. Maybe no entries?")
	}

	fmt.Printf("(maxConnections: %v) Done\n", maxConn)
	return true
}

// cleanupHeartBeat is called async to find and delete expired
// published entries in the database every two hours.
func cleanupHeartBeat() {
	for range time.Tick(time.Hour * 2) {
		if cfg.DebugMode > 0 {
			fmt.Printf("Cleanup expired and published payloads.\n")
		}
		/*
			// slower version, but more easy to read
			sql := `DELETE FROM dv.data
			          WHERE DURATION > 0 AND
					    NOW() > CREATIONDATE + CONCAT(DURATION::text, ' days')::INTERVAL`
		*/
		sql := `DELETE FROM dv.data 
				  WHERE DURATION > 0 AND 
				    CAST(NOW() - CREATIONDATE AS INT) > DURATION * 86400`
		_, err := DB.Exec(sql)
		if err != nil {
			LogInternalf("Failed to delete published and expired data (cleanupHeartBeat). %v",
				err)
		}
	}
}
