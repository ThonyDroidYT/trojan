package trojan

import (
	"crypto/sha256"
	"fmt"
	"trojan/core"
	"trojan/util"
)

// WebMenu menú de gestión web
func WebMenu() {
	fmt.Println()
	menu := []string{"Restablecer la contraseña del administrador web", "Modifique el nombre de dominio mostrado (no para solicitar un certificado)"}
	switch util.LoopInput("Por favor elige: 》", menu, true) {
	case 1:
		ResetAdminPass()
	case 2:
		SetDomain("")
	}
}

// ResetAdminPass Restablecer contraseña de administrador
func ResetAdminPass() {
	inputPass := util.Input("Ingrese la contraseña del usuario administrador: ", "")
	if inputPass == "" {
		fmt.Println("¡Deshaga los cambios!")
	} else {
		encryPass := sha256.Sum224([]byte(inputPass))
		err := core.SetValue("admin_pass", fmt.Sprintf("%x", encryPass))
		if err == nil {
			fmt.Println(util.Green("¡Contraseña de administrador restablecida correctamente!"))
		} else {
			fmt.Println(err)
		}
	}
}

// SetDomain Establecer el nombre de dominio mostrado
func SetDomain(domain string) {
	if domain == "" {
		domain = util.Input("Establecer el nombre de dominio mostrado: ", "")
	}
	if domain == "" {
		fmt.Println("¡Deshaga los cambios!")
	} else {
		core.WriteDomain(domain)
		Restart()
		fmt.Println("¡Dominio modificado con éxito!")
	}
}

// GetDomainAndPort Obtener el nombre de dominio y el puerto
func GetDomainAndPort() (string, int) {
	config := core.Load("")
	return config.SSl.Sni, config.LocalPort
}
