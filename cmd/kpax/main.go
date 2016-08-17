package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"

	"github.com/jessevdk/go-flags"
	"h12.me/config"
	"h12.me/kpax/broker"
	"h12.me/kpax/cluster"
)

const (
	clientID = "h12.me/kpax/kafpro"
)

func parseJSON(cfg interface{}) error {
	var fileConfig struct {
		ConfigFile string `long:"config" default:"config.json"`
	}
	parser := flags.NewParser(&fileConfig, flags.IgnoreUnknown)
	if _, err := parser.Parse(); err != nil {
		return err
	}
	f, err := os.Open(fileConfig.ConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(cfg)
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

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

func toJSON(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "\t")
	return string(buf)
}
