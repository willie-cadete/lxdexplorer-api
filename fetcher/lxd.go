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
	for _, h := range getHostnodes() {
		c := Connect(h)
		if c == nil {
			continue
		}
		cs, _ := c.GetContainersFull()

		for _, c := range cs {
			database.InsertOne("containers", c)
		}

		log.Println("Inserted", len(cs), "containers from", h)
	}

	time.Sleep(time.Duration(conf.Interval) * time.Second)
}
