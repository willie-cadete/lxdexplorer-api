package fetcher

import (
	"log"
	"lxdexplorer-api/config"
	"lxdexplorer-api/database"
	"os"
	"time"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
)

var conf, _ = config.LoadConfig()

type HostNode struct {
	CollectedAt time.Time   `bson:"collectedat"`
	Hostname    string      `bson:"hostname"`
	Containers  interface{} `bson:"containers"`
}

type Network struct {
	IPs     string `bson:"ips"`
	Netmask string `bson:"netmask"`
	CIDR    string `bson:"cidr"`
}

type OS struct {
	Distribution string `bson:"distribution"`
	Version      string `bson:"version"`
}

type Container struct {
	CollectedAt time.Time `bson:"collectedat"`
	Name        string    `bson:"name"`
	Hostnode    string    `bson:"hostnode"`
	Status      string    `bson:"status"`
	Network     Network   `bson:"network"`
	OS          OS        `bson:"os"`
	ImageID     string    `bson:"imageid"`
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

func ParseContainer(c api.ContainerFull, h string) Container {

	if c.State.Status == "Stopped" {
		return Container{
			Name:     c.Name,
			Hostnode: h,
			Status:   c.State.Status,
			Network: Network{
				IPs:     "N/A",
				Netmask: "N/A",
				CIDR:    "N/A",
			},
			OS: OS{
				Distribution: c.Config["image.os"],
				Version:      c.Config["image.release"],
			},
			ImageID: c.Config["volatile.base_image"][:6],
		}
	}

	return Container{
		Name:     c.Name,
		Hostnode: h,
		Status:   c.State.Status,
		Network: Network{
			IPs:     c.State.Network["eth0"].Addresses[0].Address,
			Netmask: c.State.Network["eth0"].Addresses[0].Netmask,
			CIDR:    c.State.Network["eth0"].Addresses[0].Address + "/" + c.State.Network["eth0"].Addresses[0].Netmask,
		},
		OS: OS{
			Distribution: c.Config["image.os"],
			Version:      c.Config["image.release"],
		},
		ImageID: c.Config["volatile.base_image"][:6],
	}
}

func AddLXDTTLs() {
	database.AddTTL("containers", "collectedat", int32(conf.Interval))
	log.Printf("Fetcher: Added TTL to containers collection: %d seconds", conf.Interval)

	log.Printf("Fetcher: Added TTL to history collection: %d days", conf.Retention)
	database.AddTTL("history", "collectedat", int32(conf.Retention*60*60*24))
}

func collect() {

	collectedAt := time.Now().UTC()

	for _, h := range getHostnodes() {
		c := Connect(h)
		if c == nil {
			continue
		}
		cs, _ := c.GetContainersFull()

		for _, c := range cs {
			container := ParseContainer(c, h)
			database.InsertMany("containers", []interface{}{Container{CollectedAt: collectedAt, Name: container.Name, Hostnode: container.Hostnode, Status: container.Status, Network: container.Network, OS: container.OS, ImageID: container.ImageID}})
		}
		log.Println("Fetcher: Inserted", len(cs), "containers from", h)

		database.InsertMany("history", []interface{}{HostNode{CollectedAt: collectedAt, Hostname: h, Containers: cs}})
		log.Println("Fetcher: Inserted", len(cs), "containers from hostnode:", h)

	}

	time.Sleep(time.Duration(conf.Interval) * time.Second)
}

func StartFetcher() {
	AddLXDTTLs()
	for {
		collect()
	}
}
