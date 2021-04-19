package trojan

import (
	"encoding/base64"
	"fmt"
	"trojan/core"
	"trojan/util"
)

var clientPath = "/root/config.json"

// GenClientJson Generar cliente json
func GenClientJson() {
	fmt.Println()
	var user core.User
	domain, port := GetDomainAndPort()
	mysql := core.GetMysql()
	userList, err := mysql.GetData()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if len(userList) == 1 {
		user = *userList[0]
	} else {
		UserList()
		choice := util.LoopInput("Seleccione el número de serie del usuario para generar el archivo de configuración: ", userList, true)
		if choice < 0 {
			return
		}
		user = *userList[choice-1]
	}
	pass, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		fmt.Println(util.Red("La decodificación Base64 falló: " + err.Error()))
		return
	}
	if !core.WriteClient(port, string(pass), domain, clientPath) {
		fmt.Println(util.Red("No se pudo generar el archivo de configuración!"))
	} else {
		fmt.Println("Archivo de configuración generado con éxito: " + util.Green(clientPath))
	}
}
