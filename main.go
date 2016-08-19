package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	set "github.com/deckarep/golang-set"
	"github.com/pshima/consul-alb-sync/sync"
)

const (
	version = "0.0.1"
)

var path = flag.String("path", "", "Path to config")

func init() {
	flag.StringVar(path, "p", "", "Path to config")
}

func main() {
	log.Printf("Consul Attache %s", version)

	flag.Parse()

	log.Printf("Loading config from %s", *path)

	conf, _ := sync.GetConfig(*path)

	enabled, msg := conf.Validate()
	if enabled != true {
		log.Fatalf("Exiting: %v", msg)
	}

	consul, err := sync.ConsulClient()
	if err != nil {
		log.Fatalf("Unable to create consul client: %v", err)
	}

	service, _, err := consul.Catalog().Service(conf.ServiceName, "", nil)
	if err != nil {
		log.Fatalf("Unable to query consul services: %v", err)
	}

	if len(service) < 1 {
		log.Fatalf("No services found for %s, exiting", conf.ServiceName)
	}

	serviceaddrs := set.NewSet()
	for _, node := range service {
		log.Println(node.ServiceAddress)
		instanceid, err := sync.GetInstanceIDFromIP(node.ServiceAddress)
		if err != nil {
			log.Fatalf("Error getting target group %s attributes: %v", conf.ServiceName, err)
		}
		serviceaddrs.Add(fmt.Sprintf("%s:%d", instanceid, node.ServicePort))
	}

	tg, err := sync.GetTargetGroup(conf.ServiceName)
	if err != nil {
		log.Fatalf("Error getting target group %s: %v", conf.ServiceName, err)
	}

	tgARN := tg.TargetGroups[0].TargetGroupArn

	tgh, err := sync.GetTargetGroupHealth(*tgARN)
	if err != nil {
		log.Fatalf("Error getting target group %s attributes: %v", conf.ServiceName, err)
	}

	tgaddrs := set.NewSet()
	for _, tgnode := range tgh.TargetHealthDescriptions {
		tgaddrs.Add(fmt.Sprintf("%s:%s", *tgnode.Target.Id, fmt.Sprintf("%d", *tgnode.Target.Port)))
	}

	toremove := tgaddrs.Difference(serviceaddrs)
	toadd := serviceaddrs.Difference(tgaddrs)

	for _, i := range toremove.ToSlice() {
		removehost, removeport, err := net.SplitHostPort(i.(string))
		if err != nil {
			log.Fatalf("Unable to convert list to host and port!: %v", err)
		}
		removeport64, err := strconv.ParseInt(removeport, 10, 64)
		if err != nil {
			log.Fatalf("Unable to convert port to int64!: %v", err)
		}
		err = sync.RemoveFromTargetGroup(*tgARN, removehost, removeport64)
		log.Printf("Removed: %s\n", i)
	}
	for _, a := range toadd.ToSlice() {
		addhost, addport, err := net.SplitHostPort(a.(string))
		if err != nil {
			log.Fatalf("Unable to convert list to host and port!: %v", err)
		}
		addport64, err := strconv.ParseInt(addport, 10, 64)
		if err != nil {
			log.Fatalf("Unable to convert port to int64!: %v", err)
		}
		err = sync.AddToTargetGroup(*tgARN, addhost, addport64)
		log.Printf("Added: %s\n", a)
	}

}
