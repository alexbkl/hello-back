#Backend part of Hello.Storage app

##Requirements: Golang, PostgreSQL for a Database

## Installation
```
git clone https://github.com/Hello-Storage/hello-back.git
```

## Use
```
cd hello-back
```

# Environment variables

Environment variables are stored at .env file. You can change them if you need.

For database connecetion, DATABASE_URI has to be set. For example:

## Local:

var DATABASE_URI string = "host=localhoset user=postgres password=12345 dbname=metamask port=5432 sslmode=disable TimeZone=Europe/Madrid"

## Docker:

var DATABASE_URI string = "host=host.docker.internal user=postgres password=12345 dbname=metamask port=5432 sslmode=disable TimeZone=Europe/Madrid"

# Local

#Build and run the application on the local machine:
```
go build .
```
This will output an executable file named hello-back (or hello-back.exe on Windows)
#Build and run the application in local machine:

```
go run main.go
```

##Build and run in Docker container

##From the root directory of the project, run the following command:

```
docker build -t hello-back .
```

##Run the container:

```
docker run -p 8001:8001 hello-back
```

##Your app should now be running inside a Docker container, and it should be accessible at localhost:8080 on your host machine.


#Run in background in Linux

##Execute the background application (for example, put port 80 in main.go app.Listen(":80"))

```
nohup ./hello-back > hello-back.log 2>&1 &
```

#This will output the logs to hello-back.log

#Verify that the application is running

```
ps aux | grep meta-go-api
```

#Check the PID and the port number of the app:

```
lsof -i :80
```

#Stop the application or reset in case of changes


```
kill <PID>

kill -HUP <PID>
```


#In case of advanced management of the application, a process supervisor can be used


```
systemctl
```
```
sudo apt-get install supervisord
```


#Configure domain in VPS Ubuntu:


```
sudo apt update
sudo apt install nginx
sudo nano /etc/nginx/sites-available/ounn.space

```
server {
    listen 80;
    listen [::]:80;
    server_name example.com www.example.com;

    location / {
        proxy_pass http://raw_ip:8001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}
```
sudo ln -s /etc/nginx/sites-available/myapp /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your_domain

```

To create a cron job for autorenew with SSL certbot:

Open your crontab file by running the following command in the terminal:

$ sudo crontab -e

This will open your crontab file in the default text editor. You will need to add a line to this file that tells it to run certbot renew twice per day. A common practice is to run it at noon and midnight. To do this, add the following line:

$ 0 12,0 * * * /usr/bin/certbot renew --quiet


To make requests, use the https://domain.com/api/ endpoint. (WITHOUT the port) For example, baseUrl = https://ounn.space/api/