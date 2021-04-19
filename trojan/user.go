package trojan

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"trojan/core"
	"trojan/util"
)

// UserMenu Menú de gestión de usuarios
func UserMenu() {
	fmt.Println()
	menu := []string{"Crear usuario", "Eliminar usuarios", "Limitar el tráfico", "Tráfico vacío", "Establecer fecha límite", "Cancelar fecha límite"}
	switch util.LoopInput("Por favor elige: 》", menu, false) {
	case 1:
		AddUser()
	case 2:
		DelUser()
	case 3:
		SetUserQuota()
	case 4:
		CleanData()
	case 5:
		SetupExpire()
	case 6:
		CancelExpire()
	}
}

// AddUser Agregar usuario
func AddUser() {
	randomUser := util.RandString(4)
	randomPass := util.RandString(8)
	inputUser := util.Input(fmt.Sprintf("Generar nombre de usuario aleatorio: %s, Usar retorno directo, De lo contrario, ingrese un nombre de usuario personalizado: ", randomUser), randomUser)
	if inputUser == "admin" {
		fmt.Println(util.Yellow("No se puede crear un nuevo nombre de usuario.'admin'Usuario!"))
		return
	}
	mysql := core.GetMysql()
	if user := mysql.GetUserByName(inputUser); user != nil {
		fmt.Println(util.Yellow("Nombre de usuario existente: " + inputUser + " ¡Usuario!"))
		return
	}
	inputPass := util.Input(fmt.Sprintf("Generar contraseña aleatoria: %s, Usar retorno directo, De lo contrario, ingrese una contraseña personalizada: ", randomPass), randomPass)
	base64Pass := base64.StdEncoding.EncodeToString([]byte(inputPass))
	if user := mysql.GetUserByPass(base64Pass); user != nil {
		fmt.Println(util.Yellow("La contraseña existente es: " + inputPass + " ¡Usuario!"))
		return
	}
	if mysql.CreateUser(inputUser, base64Pass, inputPass) == nil {
		fmt.Println("Usuario agregado exitosamente!")
	}
}

// DelUser eliminar usuarios
func DelUser() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Seleccione el número de serie del usuario que desea eliminar: 》", userList, true)
	if choice == -1 {
		return
	}
	if mysql.DeleteUser(userList[choice-1].ID) == nil {
		fmt.Println("Usuario eliminado correctamente!")
	}
}

// SetUserQuota Limita el tráfico de usuarios
func SetUserQuota() {
	var (
		limit int
		err   error
	)
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Seleccione el número de serie del usuario cuyo tráfico se restringirá: ", userList, true)
	if choice == -1 {
		return
	}
	for {
		quota := util.Input("Por favor ingrese usuario"+userList[choice-1].Username+"Tamaño de tráfico restringido (byte unitario)", "")
		limit, err = strconv.Atoi(quota)
		if err != nil {
			fmt.Printf("%s No un número, por favor ingrese de nuevo!\n", quota)
		} else {
			break
		}
	}
	if mysql.SetQuota(userList[choice-1].ID, limit) == nil {
		fmt.Println("Usuario configurado correctamente" + userList[choice-1].Username + "Limitar el tráfico" + util.Bytefmt(uint64(limit)))
	}
}

// CleanData Tráfico de usuarios claro
func CleanData() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Seleccione el número de serie del usuario para borrar el tráfico: ", userList, true)
	if choice == -1 {
		return
	}
	if mysql.CleanData(userList[choice-1].ID) == nil {
		fmt.Println("Tráfico vaciado correctamente!")
	}
}

// CancelExpire Cancelar de fecha límite
func CancelExpire() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Seleccione el número de serie del usuario para cancelar la fecha límite: ", userList, true)
	if choice == -1 {
		return
	}
	if userList[choice-1].UseDays == 0 {
		fmt.Println(util.Yellow("¡El usuario seleccionado no ha establecido una fecha límite!"))
		return
	}
	if mysql.CancelExpire(userList[choice-1].ID) == nil {
		fmt.Println("¡Fecha límite cancelado con éxito!")
	}
}

// SetupExpire Establecer fecha límite
func SetupExpire() {
	userList := UserList()
	mysql := core.GetMysql()
	choice := util.LoopInput("Seleccione el número de serie del usuario para establecer una fecha límite: ", userList, true)
	if choice == -1 {
		return
	}
	useDayStr := util.Input("Ingrese la cantidad de días que se restringirán: ", "")
	if useDayStr == "" {
		return
	} else if strings.Contains(useDayStr, "-") {
		fmt.Println(util.Yellow("El número de días no puede ser negativo."))
		return
	} else if !util.IsInteger(useDayStr) {
		fmt.Println(util.Yellow("¡La entrada no es un número entero!"))
		return
	}
	useDays, _ := strconv.Atoi(useDayStr)
	if mysql.SetExpire(userList[choice-1].ID, uint(useDays)) == nil {
		fmt.Println("Fecha límite establecida con éxito!")
	}
}

// CleanDataByName Borrar el tráfico de usuario especifico
func CleanDataByName(usernames []string) {
	mysql := core.GetMysql()
	if err := mysql.CleanDataByName(usernames); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("¡Vaciar el tráfico con éxito!")
	}
}

// UserList Pegue a lista de usuários e imprima-a
func UserList(ids ...string) []*core.User {
	mysql := core.GetMysql()
	userList, err := mysql.GetData(ids...)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	domain, port := GetDomainAndPort()
	for i, k := range userList {
		pass, err := base64.StdEncoding.DecodeString(k.Password)
		if err != nil {
			pass = []byte("")
		}
		fmt.Printf("%d.\n", i+1)
		fmt.Println("Nombre de usuario: 》" + k.Username)
		fmt.Println("Contraseña: 》" + string(pass))
		fmt.Println("Tráfico de subida: 》" + util.Cyan(util.Bytefmt(k.Upload)))
		fmt.Println("Tráfico de descarga: 》" + util.Cyan(util.Bytefmt(k.Download)))
		if k.Quota < 0 {
			fmt.Println("Límite de flujo: 》" + util.Cyan("Ilimitado"))
		} else {
			fmt.Println("Límite de flujo: 》" + util.Cyan(util.Bytefmt(uint64(k.Quota))))
		}
		if k.UseDays == 0 {
			fmt.Println("Fecha de expiración: 》" + util.Cyan("Ilimitado"))
		} else {
			fmt.Println("Fecha de expiración: 》" + util.Cyan(k.ExpiryDate))
		}
		fmt.Println("Compartir enlace: 》" + util.Green(fmt.Sprintf("trojan://%s@%s:%d", string(pass), domain, port)))
		fmt.Println()
	}
	return userList
}
