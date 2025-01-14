package main

import (
	"context"
	"flag"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var logger = log.Default()

func runDistributedLock() {
	var name = flag.String("name", "", "give a name")
	flag.Parse()
	var err error
	var cli *clientv3.Client
	cli, err = clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// leaseSession obviously has a lease with ttl 3 seconds
	leaseSession, err := concurrency.NewSession(cli, concurrency.WithTTL(3))
	if err != nil {
		log.Fatal(err)
	}
	defer leaseSession.Close()

	// multiple services lock the same prefix
	mutex := concurrency.NewMutex(leaseSession, *name)
	ctx := context.TODO()

	if err = mutex.Lock(ctx); err != nil {
		log.Fatal("lock failed,", err)
	}
	logger.Printf("lock %s success\n", *name)
	time.Sleep(time.Second)
	if err = mutex.Unlock(ctx); err != nil {
		log.Fatal("unlock failed,", err)
	}
	logger.Println("released lock for ", *name)
}
