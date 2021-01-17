package proxmox

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

type configPool struct {
	Name 		string `json:"poolid"`
	Description string `json:"comment"`
}

func NewConfigPool() configPool {
	return configPool{}
}

func (config configPool) CreatePool(poolId string, client *Client) (err error) {
	params := map[string]interface{}{
		"poolid":  config.Name,
		"comment": config.Description,
	}

	exitStatus, err := client.CreatePool(poolId, params)
	if err != nil {
		return fmt.Errorf("Error creating Pool: %v, error status: %s (params: %v)", err, exitStatus, params)
	}

	return
}

func (config configPool) UpdateConfig(poolId string, client *Client) (err error) {
	configParams := map[string]interface{}{
		"poolid":  config.Name,
		"comment": config.Description,
	}

	_, err = client.SetPoolConfig(poolId, configParams)
	if err != nil {
		log.Print(err)
		return err
	}

	return err
}

func NewConfigPoolFromJson(io io.Reader) (config *configPool, err error) {
	config = &configPool{}
	err = json.NewDecoder(io).Decode(config)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	log.Println(config)
	return
}

func NewConfigPoolFromApi(poolId string, client *Client) (config *configPool, err error) {
	var poolConfig map[string]interface{}
	for ii := 0; ii < 3; ii++ {
		poolConfig, err = client.GetPoolConfig(poolId)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}

		time.Sleep(8 * time.Second)
	}

	name := ""
	if _, isSet := poolConfig["poolid"]; isSet {
		name = poolConfig["poolid"].(string)
	}
	description := ""
	if _, isSet := poolConfig["comment"]; isSet {
		description = poolConfig["comment"].(string)
	}

	config = &configPool{
		Name:        name,
		Description: description,
	}

	return
}

func (c configPool) String() string {
	jsConf, _ := json.Marshal(c)
	return string(jsConf)
}
