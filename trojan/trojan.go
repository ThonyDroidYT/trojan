package trojan

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"trojan/core"
	"trojan/util"
)

// ControllMenu Menú de control de trojan
func ControllMenu() {
	fmt.Println()
	tType := Type()
	if tType == "trojan" {
		tType = "trojan-go"
	} else {
		tType = "trojan"
	}
	menu := []string{"Iniciar trojan", "Detener trojan", "Reiniciar trojan", "Verificar el estado de trojan", "Ver registro de trojan"}
	menu = append(menu, "Cambiar a"+tType)
	switch util.LoopInput("Por favor elige: ", menu, true) {
	case 1:
		Start()
	case 2:
		Stop()
	case 3:
		Restart()
	case 4:
		Status(true)
	case 5:
		go Log(300)
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		//cuadra
		<-c
	case 6:
		_ = core.SetValue("trojanType", tType)
		InstallTrojan()
	}
}

// Restart 重启trojan
func Restart() {
	if err := util.ExecCommand("systemctl restart trojan"); err != nil {
		fmt.Println(util.Red("No se pudo reiniciar trojan!"))
	} else {
		fmt.Println(util.Green("Trojan reiniciado con éxito!"))
	}
}

// Start Iniciar trojan
func Start() {
	if err := util.ExecCommand("systemctl start trojan"); err != nil {
		fmt.Println(util.Red("¡No se pudo iniciar trojan!"))
	} else {
		fmt.Println(util.Green("Trojan iniciado con éxito!"))
	}
}

// Stop Detener trojan
func Stop() {
	if err := util.ExecCommand("systemctl stop trojan"); err != nil {
		fmt.Println(util.Red("¡Detección de trojan fallida!"))
	} else {
		fmt.Println(util.Green("Trojan detenido con éxito!"))
	}
}

// Status Obtener el estado de troyano
func Status(isPrint bool) string {
	result := util.ExecCommandWithResult("systemctl status trojan")
	if isPrint {
		fmt.Println(result)
	}
	return result
}

// RunTime Tiempo de ejecución del trojan
func RunTime() string {
	result := strings.TrimSpace(util.ExecCommandWithResult("ps -Ao etime,args|grep -v grep|grep /usr/local/etc/trojan/config.json"))
	resultSlice := strings.Split(result, " ")
	if len(resultSlice) > 0 {
		return resultSlice[0]
	}
	return ""
}

// TrojanVersion Versión de trojan
func Version() string {
	flag := "-v"
	if Type() == "trojan-go" {
		flag = "-version"
	}
	result := strings.TrimSpace(util.ExecCommandWithResult("/usr/bin/trojan/trojan " + flag))
	if len(result) == 0 {
		return ""
	}
	firstLine := strings.Split(result, "\n")[0]
	tempSlice := strings.Split(firstLine, " ")
	return tempSlice[len(tempSlice)-1]
}

// TrojanType Tipo de trojan
func Type() string {
	tType, _ := core.GetValue("trojanType")
	if tType == "" {
		if strings.Contains(Status(false), "trojan-go") {
			tType = "trojan-go"
		} else {
			tType = "trojan"
		}
		_ = core.SetValue("trojanType", tType)
	}
	return tType
}

// Log Imprima registros de troyanos en tiempo real
func Log(line int) {
	result, _ := LogChan("-n "+strconv.Itoa(line), make(chan byte))
	for line := range result {
		fmt.Println(line)
	}
}

// LogChan registro de troyano en tiempo real, volver a chan
func LogChan(param string, closeChan chan byte) (chan string, error) {
	cmd := exec.Command("bash", "-c", "journalctl -f -u trojan -o cat "+param)

	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println("Error: El comando tiene un error: ", err.Error())
		return nil, err
	}
	ch := make(chan string, 100)
	stdoutScan := bufio.NewScanner(stdout)
	go func() {
		for stdoutScan.Scan() {
			select {
			case <-closeChan:
				stdout.Close()
				return
			default:
				ch <- stdoutScan.Text()
			}
		}
	}()
	return ch, nil
}
