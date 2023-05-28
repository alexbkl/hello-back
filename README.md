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
go build main.go
```
Esto resultará en un archivo ejecutable llamado main (o main.exe en Windows)

#Ejecutar la aplicación de fondo
```
nohup ./main > meta-go-api.log 2>&1 &
```

#Verificar que la aplicación se esté ejecutando
```
ps aux | grep main
```

#En caso de necesitar gestionar la aplicación de una manera más avanzada, se puede usar un process supervisor

```
systemctl
```
```
sudo apt-get install supervisord
```