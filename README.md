# ec2sshtb

## DESCRIPTION
the helper ssh to ec2 instance through bastion host.

ec2sshtb has only 2commands: sync, ssh.

```
$ ec2sshtb sync
$ ec2sshtb ssh
```

## HOW TO USE
Just you need to prepare default.yml on home dir.

### Install

```
go get github.com/reizist/ec2sshtb
```

### Prepare configuration

```
....
~/
|-- .ec2sshtb/
|   `-- default.yml
....
```

like this:

```yml
---
bastion_user: reizist
bastion_private_key_path: ~/.ssh/id_rsa
bastion_port: 2222
bastion_host: auth.bastion.example.com
host_user: ec2-user
aws_credential_profile: ec2_profile
host_port: 22
```


Or clone the https://github.com/reizist/ec2sshtb and run:

```
go build
```

## AUTHOR

reizist <reizist@gmail.com>