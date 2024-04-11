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

```
docker run -p 8080:8080 myapp-image
```

### Using Dockerfile for custom images Redis 

#### step 1: Create `redis.conf` file
#### step 2: Create `Dockerfile` as below
```
FROM redis
COPY redis.conf /usr/local/etc/redis/redis.conf
CMD ["redis-server", "/usr/local/etc/redis/redis.conf"]
```

#### step 3: Build new image from Dockerfile:
```
docker build -t my-custom-redis .
```

#### step 4: Run the container from this new image
```
docker run --name my-custom-redis-container -d my-custom-redis
```

#### run Powershell mount file redis.conf to file config redis on Docker image
```
docker run --name custom-redis -v /var/local/redis/redis.conf:/usr/loca
l/etc/redis/redis.conf -d redis redis-server /usr/local/etc/redis/redis.conf
```