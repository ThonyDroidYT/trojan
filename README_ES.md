# trojan
![](https://img.shields.io/github/v/release/Jrohy/trojan.svg) 
![](https://img.shields.io/docker/pulls/jrohy/trojan.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Jrohy/trojan)](https://goreportcard.com/report/github.com/Jrohy/trojan)
[![Downloads](https://img.shields.io/github/downloads/Jrohy/trojan/total.svg)](https://img.shields.io/github/downloads/Jrohy/trojan/total.svg)
[![License](https://img.shields.io/badge/license-GPL%20V3-blue.svg?longCache=true)](https://www.gnu.org/licenses/gpl-3.0.en.html)


programa de implementación de gestión multiusuario troyan

## Características
- Administre los multiusuarios troyanos de dos maneras: página web en línea y línea de comandos
- Iniciar / detener / reiniciar el servidor troyan
- Revisar las estadísticas de tráfico y el límite de tráfico
- Gestión del modo de línea de comandos, finalización de comandos de soporte
- Aplicación de certificado acme.sh integrada
- Generar archivo de configuración del cliente
- Ver registros de troyan en línea en tiempo real
- Troyan en línea y conmutador troyan-go en cualquier momento
- Soporte de URL trojan://enlace para compartir y compartir código QR (el código QR solo está disponible en páginas web)
- Limitar el período de uso del usuario

## Metodo de instalacion
*Trojan Prepare el nombre de dominio disponible para el servidor con anticipación*  

###  a. Instalación de script con un clic
```
#Instalar actualización
source <(curl -sL https://www.tdproyectos.tk/trojan/install.sh)

#Desinstalar
source <(curl -sL https://www.tdproyectos.tk/trojan/install.sh) --remove

```
Después de la instalación, ingrese el comando 'trojan' para ingresar al programa de administración
Acceso al navegador https://NombreDelDominio, para administrar usuarios trojan desde la página web en línea
Dirección del código fuente de la página de inicio: [trojan-web](https://github.com/Jrohy/trojan-web)

### b.correr con docker
1. Instalar mysql

Debido a que el uso de la memoria mariadb es al menos la mitad que el de mysql, se recomienda utilizar la base de datos mariadb
```
docker run --name trojan-mariadb --restart=always -p 3306:3306 -v /home/mariadb:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=trojan -e MYSQL_ROOT_HOST=% -e MYSQL_DATABASE=trojan -d mariadb:10.2
```
El puerto, la contraseña de root y el directorio persistente se pueden cambiar a otros

2. Instalar trojan
```
docker run -it -d --name trojan --net=host --restart=always --privileged jrohy/trojan init
```
Entrar en el contenedor después de correr `docker exec -it trojan bash`, Luego escribe 'trojan' Puede iniciar la instalación inicial   

Iniciar servicio web: `systemctl start trojan-web`   

Configurar el inicio automático: `systemctl enable trojan-web`

Programa de gestión de actualizaciones: `source <(curl -sL https://git.io/trojan-install)`

## Captura de pantalla de Ejecución
![avatar](asset/1.png)
![avatar](asset/2.png)

## Línea de comando
```
Uso:
  trojan [flags]
  trojan [comando]

Comandos disponibles:
  add           Agregar usuario
  clean         Borrar el tráfico de usuarios especificado
  completion    Finalización automática de comandos (soporte bash y zsh)
  del           eliminar usuarios
  help          Ayuda sobre cualquier comando
  info          Lista de información de usuario
  log           Ver registros de trojan
  restart       Reiniciar trojan
  start         Iniciar trojan
  status        Verificar el estado de trojan
  stop          Detener trojan
  tls           Instalación de certificado
  update        Actualizar trojan
  updateWeb     Actualizar el programa de gestión de trojan
  version       Mostrar número de versión
  import [path] Importar archivo sql
  export [path] Exportar archivo sql
  web           Empezar por web

Flags:
  -h, --help   ayuda para trojan
```

## Nota
Después de instalar trojan, se recomienda encarecidamente activar BBR y otras aceleraciones: [Linux-NetSpeed](https://github.com/chiakge/Linux-NetSpeed)  

Cliente trojan recomendados: 
   - pc: [Trojan-Qt5](https://github.com/TheWanderingCoel/Trojan-Qt5)
   - ios: [shadowrocket](https://apps.apple.com/us/app/shadowrocket/id932747118)
   - android: [igniter](https://github.com/trojan-gfw/igniter)
