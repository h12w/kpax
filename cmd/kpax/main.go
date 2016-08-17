package main

import (
	"log"

	"h12.me/config"
	"h12.me/kpax/broker"
	"h12.me/kpax/cluster"
)

const (
	clientID = "h12.me/kpax/kafpro"
)

func main() {
	var cfg Config
	cmd, err := config.ParseCommand(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	c := cluster.New(broker.New, cfg.Brokers)
	switch cmd.Name {
	case "consume":
		err = cfg.Consume.Exec(c)
	case "offset":
		err = cfg.Offset.Exec(c)
	case "rollback":
		err = cfg.Rollback.Exec(c)
	case "tail":
		err = cfg.Tail.Exec(c)
	case "produce":
		err = cfg.Produce.Exec(c)
	case "meta":
		err = cfg.Meta.Exec(c)
	default:
		log.Fatal("unkown command " + cmd.Name)
	}
	if err != nil {
		log.Fatal(err)
	}
}
