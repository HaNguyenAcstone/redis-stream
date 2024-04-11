# Redis Stream

## Setup Eviroment
### Step 1: Setup server Linux (Unbutu 22.04 Server LTS)

### Step 2: Create SSH key from computer to connect the Server
#### Powershell:
```
ssh-keygen -t rsa -b 4096
```

#### Access path `./C/users/my-computer/.ssh` added <b>ssh</b>

### | Setup Docker On Unbutu Server
#### Access the link to setup docker on unbutu: <a href="https://www.digitalocean.com/community/tutorials/how-to-install-and-use-docker-on-ubuntu-22-04">Setup Docker On Linux (Unbutu)</a>


### | Setup Redis on Docker

Powershell: 
```
    ~root@name: docker pull redis
    ~root@name: docker run --name redis-name -p 6379:6379 -d redis
```

### <b>Note*: redis-name for the name images create redis on docker</b>

