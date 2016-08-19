package sync

import (
	//"fmt"

	"github.com/mitchellh/consulstructure"
)

type Config struct {
	Enabled     string
	ServiceName string
}

func GetConfig(prefix string) (*Config, error) {
	updateCh := make(chan interface{})
	errCh := make(chan error)
	decoder := &consulstructure.Decoder{
		Target:   &Config{},
		Prefix:   prefix,
		UpdateCh: updateCh,
		ErrCh:    errCh,
	}

	go decoder.Run()
	for {
		select {
		case v := <-updateCh:
			return v.(*Config), nil
		case err := <-errCh:
			return nil, err
		}
	}
}

func (c *Config) Validate() (bool, string) {
	if c.Enabled == "" {
		return false, "Enabled key empty, check your path or set enabled = true"
	} else if c.ServiceName == "" {
		return false, "ServiceName key empty, check your path or set servicename = nameofservice"
	} else {
		return true, ""
	}
}
