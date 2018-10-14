package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"regexp"

	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	finder "github.com/b4b4r07/go-finder"
	shellwords "github.com/mattn/go-shellwords"
)

// BaseDir is the directory path of config, hosts file located
const BaseDir = "/.bastion_ec2ssh/"

// ConfigFileName is the file name of the config
const ConfigFileName = "default.yml"

// HostFileName is the file name of the hosts
const HostFileName = "hosts.yml"

// Config is
type Config struct {
	BastionUser           string `yaml:"bastion_user"`
	BastionPrivateKeyPath string `yaml:"bastion_private_key_path"`
	BastionHost           string `yaml:"bastion_host"`
	BastionPort           int    `yaml:"bastion_port"`
	HostUser              string `yaml:"host_user"`
	HostPort              int    `yaml:"host_port"`
	AwsCredentialProfile  string `yaml:"aws_credential_profile"`
}

func userDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func parseConfig() *Config {
	config := &Config{}

	buf, err := ioutil.ReadFile(userDir() + BaseDir + ConfigFileName)

	if err != nil {
		fmt.Println("not exists defaults.yml on " + BaseDir)
		os.Exit(2)
	}

	err = yaml.Unmarshal(buf, config)

	if err != nil {
		panic(err)
	}
	return config
}

func parseHosts() map[string]string {
	buf, err := ioutil.ReadFile(userDir() + BaseDir + HostFileName)

	if err != nil {
		fmt.Println("not exists hosts.yml on " + BaseDir)
		os.Exit(2)
	}

	var m map[string]string
	err = yaml.Unmarshal(buf, &m)

	if err != nil {
		panic(err)
	}

	return m
}

func runCmdStr(cmdstr string) error {
	c, err := shellwords.Parse(cmdstr)
	fmt.Println(c)

	if err != nil {
		return err
	}
	switch len(c) {
	case 0:
		return nil
	default:
		cmd := exec.Command(c[0], c[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// TODO: I wonder why it can not execute
		// err = cmd.Run()
	}
	if err != nil {
		return err
	}
	return nil
}

func ssh() {
	config := parseConfig()
	hosts := parseHosts()
	keys := make([]string, len(hosts))

	peco, _ := finder.New("peco")

	for k, v := range hosts {
		keys = append(keys, k)
		peco.Add(k, v)
	}
	selectedHosts, err := peco.Select()
	if err != nil {
		panic(err)
	}
	selectedHost := selectedHosts[0]
	fmt.Printf("Connecting to '%s' with '%s' via rcloud bastion with '%s'.\n", selectedHost, config.HostUser, config.BastionUser)
	cmdstr := fmt.Sprintf("ssh %s@%s -p %d -o ProxyCommand=\\'ssh %s@%s -p %d -i %s -W %%h:%%p\\'", config.HostUser, selectedHost, config.HostPort, config.BastionUser, config.BastionHost, config.BastionPort, config.BastionPrivateKeyPath)
	runCmdStr(cmdstr)
}

func awsEc2Client(profile string, region string) *ec2.EC2 {
	var config aws.Config
	if profile != "" {
		creds := credentials.NewSharedCredentials("", profile)
		config = aws.Config{Region: aws.String(region), Credentials: creds}
	} else {
		config = aws.Config{Region: aws.String(region)}
	}
	sess := session.New(&config)
	ec2Client := ec2.New(sess)
	return ec2Client
}

func getInstanceName(instance *ec2.Instance) (instanceName string) {
	for _, t := range instance.Tags {
		if *t.Key == "Name" {
			instanceName = *t.Value
		}
	}
	if instanceName == "" {
		instanceName = *instance.Tags[0].Value
	}
	return
}

func saveToFile(config *Config) {
	instances := listInstances(config.AwsCredentialProfile)
	filePath := userDir() + BaseDir + HostFileName
	hosts_file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer hosts_file.Close()

	if err != nil {
		log.Fatal(err)
	}

	writer := bufio.NewWriter(hosts_file)

	for _, v := range instances {
		lineStr := fmt.Sprintf("%s (%s): %s\n", getInstanceName(v), *v.InstanceId, *v.PrivateIpAddress)
		writer.WriteString(lineStr)
	}
	writer.Flush()
}

func listInstances(profile string) []*ec2.Instance {
	cli := awsEc2Client(profile, "ap-northeast-1")
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("running"),
					aws.String("pending"),
				},
			},
		},
	}
	var instances []*ec2.Instance

	res, _ := cli.DescribeInstances(params)
	for _, v := range res.Reservations {
		if v != nil {
			for _, w := range v.Instances {
				instanceName := getInstanceName(w)
				if !regexp.MustCompile(`Vyos*`).MatchString(instanceName) {
					instances = append(instances, w)
				}
			}
		} else {
			break
		}
	}
	return instances
}

func sync() {
	config := parseConfig()
	saveToFile(config)
}

func main() {
	if len(os.Args) < 2 {
		ssh()
	} else {
		subcommand := os.Args[1]
		if subcommand == "sync" {
			sync()
		} else {
			ssh()
		}
	}
}
