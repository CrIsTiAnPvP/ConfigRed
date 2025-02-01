package main

import (
	"bufio"
	"configred/interfacesv4"
	"configred/rainbow"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

func getHostname() string {
	// getHostname retrieves the hostname of the current machine.
	// It returns "Unknown" if there is an error while executing the command.
	cmd := exec.Command("hostname")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "Unknown"
	}

	return strings.TrimSpace(out.String())
}

func clear() {
	// clear clears the terminal screen by executing the "cls" command
	// in a Windows command prompt. It sets the command's standard output
	// to the current process's standard output.
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func menu() {
	clear()
	println("💻|", rainbow.Color("Configuración de Red"), "🌐")
	println()
	println("1. Listar interfacesv4")
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
	var interfaz string
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
			bar := progressbar.NewOptions(100, progressbar.OptionClearOnFinish(), progressbar.OptionSetPredictTime(false), progressbar.OptionSetDescription("Buscando interfacesv4"))
			for i := 0; i < 100; i++ {
				bar.Add(1)
				time.Sleep(8 * time.Millisecond)
			}
			interfacesv4.PrintList(true)
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
						interfacesv4.PrintList(true)
						println()
						print(fmt.Sprintf("\033[38;5;38mSeleccione una interfaz (1-%d): ", len(interfacesv4.List(true))))
						ifaceB, _ := reader.ReadString('\n')
						ifaceB = strings.TrimSpace(ifaceB)
						iface, _ := strconv.Atoi(ifaceB)

						for iface < 1 || iface > len(interfacesv4.List(true)) {
							println("\033[31mInterfaz no válida, intente de nuevo.\033[0m")
							print(fmt.Sprintf("\033[38;5;38mSeleccione una interfaz (1-%d): ", len(interfacesv4.List(true))))
							ifaceB, _ = reader.ReadString('\n')
							ifaceB = strings.TrimSpace(ifaceB)
							iface, _ = strconv.Atoi(ifaceB)
						}

						if iface >= 1 && iface <= len(interfacesv4.List(true)) {
							interfaz = strings.Split(interfacesv4.List(true)[iface-1], "         ")[1]
							interfaz = strings.Split(interfacesv4.RemoveANSI(interfaz), "│")[0]
							interfaz = strings.TrimSpace(interfaz)
						}
					} else {
						clear()
						println("💻|", rainbow.Color("Configuración de Red ➡️ Configurar una red"), "🌐")
						println()
						interfacesv4.PrintList(false)
						println()
						print(fmt.Sprintf("\033[38;5;38mSeleccione una interfaz (1-%d): ", len(interfacesv4.List(false))))
						ifaceB, _ := reader.ReadString('\n')
						ifaceB = strings.TrimSpace(ifaceB)
						iface, _ := strconv.Atoi(ifaceB)

						for iface < 1 || iface > len(interfacesv4.List(false)) {
							println("\033[31mInterfaz no válida, intente de nuevo.\033[0m")
							print(fmt.Sprintf("\033[38;5;38mSeleccione una interfaz (1-%d): ", len(interfacesv4.List(false))))
							ifaceB, _ = reader.ReadString('\n')
							ifaceB = strings.TrimSpace(ifaceB)
							iface, _ = strconv.Atoi(ifaceB)
						}

						if iface >= 1 && iface <= len(interfacesv4.List(false)) {
							interfaz = strings.Split(interfacesv4.List(false)[iface-1], "         ")[1]
							interfaz = strings.Split(interfacesv4.RemoveANSI(interfaz), "│")[0]
							interfaz = strings.TrimSpace(interfaz)
						}
					}
					config := interfacesv4.ParseConfig(interfacesv4.GetConfig(interfaz))
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración de Interfaz | "+interfaz), "🌐")
					println()
					interfacesv4.PrintConfig(config)
					os.Exit(1)
				}
			}
		case 4:
			return
		default:
			println("Opción no válida")
		}
	}
}
