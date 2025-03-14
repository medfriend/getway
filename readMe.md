## getway medfriend
___
### variables de ambiente 
existen dos archivos de configuracion los cuales son .env para el ambiente de desarrollo desplegando los microservicios desde desarrollo o sea utilizando go build dentro del ide o el terminal, para el ambiente donde los microservicios se encuentran empaquetados es el archivo .env-devprod

---
### definicion de nombre de servicios

para los nombre de los componentes debe de medfri-[nombre de la llave en minisculas]

---

### definicion de address 
para los servicios que son de tipo rest estos servicios son localhost debido a que estos servicios crean sus propios server, para los servicios que son alcanzados desde zeromq el address es el server de zeromq o sea la ip de la traza-go o del serivicio que este administrando esas configuraciones