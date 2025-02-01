package main

import (
	"bufio"
	"configred/interfaces"
	"configred/rainbow"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func getHostname() string {
	cmd := exec.Command("hostname")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "error"
	}

	return strings.TrimSpace(out.String())
}

func clear() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func menu() {
	clear()
	println("💻|", rainbow.Color("Configuración de Red"), "🌐")
	println()
	println("1. Listar Interfaces")
	println("2. Configurar una interfaz")
	println("3. Configurar una ruta")
	println("4. Salir")
	println()
}

func menuRed() {
	clear()
	println("💻|", rainbow.Color("Configuración de Red ➡️ Configurar una red"), "🌐")
	println()
	println("1. Configuración estática")
	println("2. Configuración dinámica")
	println("3. Configuración actual")
	println("4. Volver")
	println()
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		menu()
		print("Seleccione una opción: ")
		opcionStr, _ := reader.ReadString('\n')
		opcionStr = strings.TrimSpace(opcionStr)
		opcion, _ := strconv.Atoi(opcionStr)
		switch opcion {
		case 1:
			clear()
			bar := progressbar.NewOptions(100, progressbar.OptionClearOnFinish(), progressbar.OptionSetPredictTime(false), progressbar.OptionSetDescription("Buscando interfaces"))
			for i := 0; i < 100; i++ {
				bar.Add(1)
				time.Sleep(8 * time.Millisecond)
			}
			interfaces.PrintList(true)
			println("Presione Enter para continuar...")
			reader.ReadString('\n')
		case 2:
			for {
				menuRed()
				print(rainbow.Color(fmt.Sprintf("root🛜%s: ", getHostname())))
				optStr, _ := reader.ReadString('\n')
				optStr = strings.TrimSpace(optStr)
				opt, _ := strconv.Atoi(optStr)
				switch opt {
				case 1:
					return
				case 2:
					return
				case 3:
					println()
					print("¿Mostrar todas las interfaces S/n? ")
					showOpt, _ := reader.ReadString('\n')
					show := strings.TrimSpace(showOpt)
					if show == "S" || show == "s" || show == "" {
						clear()
						println("💻|", rainbow.Color("Configuración de Red ➡️ Configurar una red"), "🌐")
						println()
						interfaces.PrintList(true)
						println()
						print(fmt.Sprintf("\033[38;5;38mSeleccione una interfaz (1-%d): ", len(interfaces.List(true))))
						ifaceB, _ := reader.ReadString('\n')
						ifaceB = strings.TrimSpace(ifaceB)
						iface, _ := strconv.Atoi(ifaceB)

						if iface >= 1 && iface <= len(interfaces.List(true)) {
							interfaz := strings.Split(interfaces.List(true)[iface-1], "         ")[1]
							interfaz = strings.Split(interfaces.RemoveANSI(interfaz), "│")[0]
							interfaz = strings.TrimSpace(interfaz)
							interfaces.GetConfig(interfaz)
							os.Exit(1)
						}
					} else {
						interfaces.PrintList(false)
						reader.ReadString('\n')
					}
				}
			}
		case 4:
			return
		default:
			println("Opción no válida")
		}
	}
}
