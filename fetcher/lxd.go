package fetcher

import (
	"log"
	"lxdexplorer-api/config"
	"lxdexplorer-api/database"
	"os"
	"time"

	lxd "github.com/canonical/lxd/client"
)

var conf, _ = config.LoadConfig()

type HostNode struct {
	CollectedAt time.Time   `bson:"collectedat"`
	Hostname    string      `bson:"hostname"`
	Containers  interface{} `bson:"containers"`
}

type Container struct {
	CollectedAt time.Time   `bson:"collectedat"`
	Name        string      `bson:"name"`
	Container   interface{} `bson:"container"`
	Backups     interface{} `bson:"backups"`
	State       interface{} `bson:"state"`
	Snapshots   interface{} `bson:"snapshots"`
}

func connectionOptions() *lxd.ConnectionArgs {
	c := conf

	TLSCertificate, _ := os.ReadFile(c.LXD.TLSCertificate)
	TLSKey, _ := os.ReadFile(c.LXD.TLSKey)

	args := lxd.ConnectionArgs{
		TLSClientCert:      string(TLSCertificate),
		TLSClientKey:       string(TLSKey),
		InsecureSkipVerify: !c.LXD.CertificateVerify,
		SkipGetServer:      false,
	}

	return &args

}

func Connect(h string) lxd.InstanceServer {
	args := connectionOptions()

	cnn, err := lxd.ConnectLXD("https://"+h+":8443", args)
	if err != nil {
		log.Println(err)
	}
	return cnn
}

func getHostnodes() []string {
	c := conf
	return c.HostNodes
}

func Run() {

	collectedAt := time.Now().UTC()

	for _, h := range getHostnodes() {
		c := Connect(h)
		if c == nil {
			continue
		}
		cs, _ := c.GetContainersFull()

		// for _, c := range cs {
		// 	database.InsertOne("containers", c)
		// }
		// log.Println("Inserted", len(cs), "containers from", h)
		for _, c := range cs {
			database.InsertMany("containers", []interface{}{Container{CollectedAt: collectedAt, Name: c.Name, Container: c.Container, Backups: c.Backups, State: c.State, Snapshots: c.Snapshots}})
		}
		log.Println("Inserted", len(cs), "containers from", h)

		database.InsertMany("hostnodes", []interface{}{HostNode{CollectedAt: collectedAt, Hostname: h, Containers: cs}})
		log.Println("Inserted", len(cs), "containers from hostnode:", h)

	}

	database.AddTTL("containers", "collectedat", int32(conf.Interval*2))
	database.AddTTL("hostnodes", "collectedat", conf.Retention*24*60)

	time.Sleep(time.Duration(conf.Interval) * time.Second)
}
