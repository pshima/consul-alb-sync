package sync

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
)

func ConsulClient() (*consulapi.Client, error) {
	consul, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		return nil, fmt.Errorf("Error initializing consul client: %v", err)
	}
	return consul, nil
}
