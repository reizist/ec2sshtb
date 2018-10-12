package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"

	"gopkg.in/yaml.v2"
)

// const ConfigPath :=

type Config struct {
	BastionUser           string `yaml:"bastion_user"`
	BastionPrivateKeyPath string `yaml:"bastion_private_key_path"`
	HostUser              string `yaml:"host_user"`
	AwsCredentialProfile  string `yaml:"aws_credential_profile"`
}

func parseConfig() *Config {
	config := &Config{}

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	buf, err := ioutil.ReadFile(usr.HomeDir + "/.rcssh/default.yml")

	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf, config)

	if err != nil {
		panic(err)
	}
	return config
}

func main() {
	const BastionHost = "auth.rtc-rcloud.jp"
	const BastionPort = 443
	const HostPort = 14151
	// credentials
	const AwsConfig = ""
	const AwsDefaultRegion = "ap-northeast-1"
	// hosts書き込むパス
	const HostsPath = ""

	config := parseConfig()
	fmt.Printf("config: %+v", config)
}
