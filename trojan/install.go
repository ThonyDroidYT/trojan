package trojan

import (
	"fmt"
	"github.com/gobuffalo/packr/v2"
	"net"
	"strconv"
	"strings"
	"time"
	"trojan/core"
	"trojan/util"
)

var (
	dockerInstallUrl1 = "https://get.docker.com"
	dockerInstallUrl2 = "https://git.io/docker-install"
	dbDockerRun       = "docker run --name trojan-mariadb --restart=always -p %d:3306 -v /home/mariadb:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=%s -e MYSQL_ROOT_HOST=%% -e MYSQL_DATABASE=trojan -d mariadb:10.2"
)

// InstallMenu 安装目录
func InstallMenu() {
	fmt.Println()
	menu := []string{"Actualizar trojan", "Solicitud de certificado", "Instalar mysql"}
	switch util.LoopInput("por favor elige: ", menu, true) {
	case 1:
		InstallTrojan()
	case 2:
		InstallTls()
	case 3:
		InstallMysql()
	default:
		return
	}
}

// InstallDocker 安装docker
func InstallDocker() {
	if !util.CheckCommandExists("docker") {
		util.RunWebShell(dockerInstallUrl1)
		if !util.CheckCommandExists("docker") {
			util.RunWebShell(dockerInstallUrl2)
		} else {
			util.ExecCommand("systemctl enable docker")
			util.ExecCommand("systemctl start docker")
		}
		fmt.Println()
	}
}

// InstallTrojan Instalar trojan
func InstallTrojan() {
	fmt.Println()
	box := packr.New("trojan-install", "../asset")
	data, err := box.FindString("trojan-install.sh")
	if err != nil {
		fmt.Println(err)
	}
	if util.ExecCommandWithResult("systemctl list-unit-files|grep trojan.service") != "" && Type() == "trojan-go" {
		data = strings.ReplaceAll(data, "TYPE=0", "TYPE=1")
	}
	util.ExecCommand(data)
	util.OpenPort(443)
	util.ExecCommand("systemctl restart trojan")
	util.ExecCommand("systemctl enable trojan")
}

// InstallTls Instale el certificado
func InstallTls() {
	domain := ""
	fmt.Println()
	choice := util.LoopInput("Elija el método de certificado: ", []string{"Certificado Let's Encrypt Automático", "Ruta de certificado personalizado"}, true)
	if choice < 0 {
		return
	} else if choice == 1 {
		localIP := util.GetLocalIP()
		fmt.Printf("IP nativa: %s\n", localIP)
		for {
			domain = util.Input("Ingrese el nombre de dominio para el certificado: ", "")
			ipList, err := net.LookupIP(domain)
			fmt.Printf("%s IP analizada: %v\n", domain, ipList)
			if err != nil {
				fmt.Println(err)
				fmt.Println("El nombre de dominio es incorrecto, vuelva a ingresar")
				continue
			}
			checkIp := false
			for _, ip := range ipList {
				if localIP == ip.String() {
					checkIp = true
				}
			}
			if checkIp {
				break
			} else {
				fmt.Println("El nombre de dominio ingresado no coincide con la IP local, ¡vuelva a ingresar!")
			}
		}
		util.InstallPack("socat")
		if !util.IsExists("/root/.acme.sh/acme.sh") {
			util.RunWebShell("https://get.acme.sh")
		}
		util.ExecCommand("systemctl stop trojan-web")
		util.OpenPort(80)
		util.ExecCommand(fmt.Sprintf("bash /root/.acme.sh/acme.sh --issue -d %s --debug --standalone --keylength ec-256", domain))
		crtFile := "/root/.acme.sh/" + domain + "_ecc" + "/fullchain.cer"
		keyFile := "/root/.acme.sh/" + domain + "_ecc" + "/" + domain + ".key"
		core.WriteTls(crtFile, keyFile, domain)
	} else if choice == 2 {
		crtFile := util.Input("Ingrese la ruta del archivo cert del certificado: ", "")
		keyFile := util.Input("Introduzca la ruta del archivo key del certificado: ", "")
		if !util.IsExists(crtFile) || !util.IsExists(keyFile) {
			fmt.Println("El archivo key o certificado ingresado no existe!")
		} else {
			domain = util.Input("Ingrese el nombre de dominio correspondiente a este certificado: ", "")
			if domain == "" {
				fmt.Println("El nombre de dominio ingresado está vacío!")
				return
			}
			core.WriteTls(crtFile, keyFile, domain)
		}
	}
	Restart()
	util.ExecCommand("systemctl restart trojan-web")
	fmt.Println()
}

// InstallMysql Instalar mysql
func InstallMysql() {
	var (
		mysql  core.Mysql
		choice int
	)
	fmt.Println()
	if util.IsExists("/.dockerenv") {
		choice = 2
	} else {
		choice = util.LoopInput("por favor elige: ", []string{"Instalar la versión docker de mysql (mariadb)", "Ingrese una conexión mysql personalizada"}, true)
	}
	if choice < 0 {
		return
	} else if choice == 1 {
		mysql = core.Mysql{ServerAddr: "127.0.0.1", ServerPort: util.RandomPort(), Password: util.RandString(5), Username: "root", Database: "trojan"}
		InstallDocker()
		fmt.Println(fmt.Sprintf(dbDockerRun, mysql.ServerPort, mysql.Password))
		if util.CheckCommandExists("setenforce") {
			util.ExecCommand("setenforce 0")
		}
		util.OpenPort(mysql.ServerPort)
		util.ExecCommand(fmt.Sprintf(dbDockerRun, mysql.ServerPort, mysql.Password))
		db := mysql.GetDB()
		for {
			fmt.Printf("%s mariadb está comenzando, por favor espere...\n", time.Now().Format("2006-01-02 15:04:05"))
			err := db.Ping()
			if err == nil {
				db.Close()
				break
			} else {
				time.Sleep(2 * time.Second)
			}
		}
		fmt.Println("mariadb comenzó con éxito!")
	} else if choice == 2 {
		mysql = core.Mysql{}
		for {
			for {
				mysqlUrl := util.Input("Ingrese la dirección de conexión de mysql (formato: host:port), La dirección de conexión predeterminada es 127.0.0.1:3306, Usar retorno directo, De lo contrario, ingrese una dirección de conexión personalizada: ",
					"127.0.0.1:3306")
				urlInfo := strings.Split(mysqlUrl, ":")
				if len(urlInfo) != 2 {
					fmt.Printf("El %s ingresado no coincide con el formato coincidente(host:port)\n", mysqlUrl)
					continue
				}
				port, err := strconv.Atoi(urlInfo[1])
				if err != nil {
					fmt.Printf("%s No un número\n", urlInfo[1])
					continue
				}
				mysql.ServerAddr, mysql.ServerPort = urlInfo[0], port
				break
			}
			mysql.Username = util.Input("Ingrese el nombre de usuario de mysql (ingrese para usar root): ", "root")
			mysql.Password = util.Input(fmt.Sprintf("Ingrese la contraseña del usuario mysql %s: ", mysql.Username), "")
			db := mysql.GetDB()
			if db != nil && db.Ping() == nil {
				mysql.Database = util.Input("Introduzca el nombre de la base de datos utilizada (se puede crear automáticamente si no existe, Ingrese para usar trojan): ", "trojan")
				db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;", mysql.Database))
				break
			} else {
				fmt.Println("No se pudo conectar a mysql, vuelva a ingresar")
			}
		}
	}
	mysql.CreateTable()
	core.WriteMysql(&mysql)
	if userList, _ := mysql.GetData(); len(userList) == 0 {
		AddUser()
	}
	Restart()
	fmt.Println()
}
