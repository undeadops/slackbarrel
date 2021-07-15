# Barrel - Configuration Management

The Barrel minimal Config management System.

## Basic Architecture

Requirements were for a Chef or Puppet like system.  Which requires a central management server and agents that poll that management server.  This is the basic architecture of the Barrel config management system.  It has a Server and Agent(s).  Configuration is for both server and agents is a simple yaml config file.  The management server exposes a set of HTTP endpoints which should be accessible to the agents.  

## Installation

Installing these process is fairly straight forward.  Being written in Go, their are two binaries.  The Server is compiled for both MacOS and Linux, and should run fine on either.  The Agent is a bit more purpose driven and will require a Linux Host(and requirements were for specifically for managaging Ubuntu Linux Hosts). 

Compiling from source, has a provided `Makefile`.  So all that is required to build, is Go, and make.  This will produce a few binaries in the `build` directory. (I have however included them with this) 

The agent should be deployed to the managed hosts, at `/usr/local/bin/barrel-agent`.
Since these are Ubuntu Hosts with `systemd`, I've also provided a systemd config file to launch this agent.  You can copy the local `barrel-agent.service` to `/etc/systemd/system/barrel-agent.service` on the managed hosts.  Reload systemd with `systemctl daemon-reload`.  Start barrel-agent with the command `systemctl start barrel-agent` and enable it at start with `systemctl enable barrel-agent`. The, config for the agent should live in `/etc/barrel/agent.conf` An example one is also provided named `agent.yaml`. 

The server is a bit more free-form.  There is a systemd config file for it as well, and its config should also live in `/etc/barrel/server.conf` as well as a Data Directory for the YAML used to configure the agents.  However, they're also runable from the local system by just issuing the command: `./barrel-server` which will automaticaly look for its config file at `./config.yaml` and the data directory at `./data`.  Default Port binding is `5000`.   Which is all adjustable with the `-help` flag issued to either binaries. 

To be able to run the server from behind a firewall, I've been testing with the use of ngrok[https://ngrok.com/].  And just updating the agent server url accordingly.

With the free plan: `ngrok http 5000` - Follow ngrok's documentation for how to setup/configure it for your use.

And running the agent with the following command `agent -config agent.yaml -server https://8b3430e1b482.ngrok.io`
 
## Configuration of the Server

Configuration for the Management Server is done via yaml file, and placing accompanying data files in the data directory. 

Config file example looks like:

```yaml
---
apps:
  foobar:
    package:
    - name: vim
      state: absent
    - name: apache2
      state: present
    - name: libapache2-mod-php
      state: present
    - name: php
      state: present
    - name: nginx
      state: absent
    files:
    - path: /var/www/html/index.html
      state: absent
    - src: foobar/index.php
      path: /var/www/html/index.php
      owner: www-data
      group: www-data
      mode: 644
      state: present
    service:
    - name: apache2
      action: reload
```

Top level, is apps, it allows for a map set, of multiple apps.
Within each app, has keys, `package`, `files`, and `service`.  Packages are basic package names you would apt-get install on an ubuntu host. with a state of present or absent, to describe if it should install the package or remove it from the system. 

Files, are basic files you want to copy to the managed hosts.  These files will be pulled via an HTTP endpoint from the management host. The Source Path (`src`), is relative to the DataDir configuration of the server. The `path`, is the destination path on the agent system.  This path will be overwritten when the `src` sha is different than the `path` destination.  Permissions are fixed up on these files as well.  Each file has a state key as well, which will describe if the file should be present on the system or if it should be absent from the system.  Only `path` and `state` are required to be defined. 

Service, should be a service managed via systemd. Managed services will be restarted on changes to files. Action would describe a systemctl reload or systemctl restart on the service. 

### Known Defeciencies

I ran out of time to work on it, and there are a few rough edges around what its doing.  Reloading/ReReading the config from the server,
requires restarting the server, which is less than ideal.


