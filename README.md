#Parte backend para proyecto final de Desarrollo de Aplicaciones Web

##requisitos: Golang, PostgreSQL

## Instalación
```
git clone https://github.com/alexanderbkl/meta-go-api.git
```

## Uso
```
cd meta-go-api
```

#Compilar la aplicación Go
```
go build .
```
Esto resultará en un archivo ejecutable llamado main (o main.exe en Windows)

#Ejecutar la aplicación de fondo (poner puerto 80 en main.go app.Listen(":80"))
```
nohup ./meta-go-api > meta-go-api.log 2>&1 &
```

#Verificar que la aplicación se esté ejecutando
```
ps aux | grep meta-go-api
```

#Terminar proceso o reiniciar en caso de cambios


```
kill <PID>

kill -HUP <PID>
```


#En caso de necesitar gestionar la aplicación de una manera más avanzada, se puede usar un process supervisor

```
systemctl
```
```
sudo apt-get install supervisord
```


#Configurar dominio en VPS Ubuntu:

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