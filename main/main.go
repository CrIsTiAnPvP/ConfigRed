package main

import (
	"bufio"
	"configred/interfacesv4"
	"configred/rainbow"
	"configred/utils"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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
	println("1. Listar interfaces")
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
			bar := progressbar.NewOptions(100, progressbar.OptionClearOnFinish(), progressbar.OptionSetPredictTime(false), progressbar.OptionSetDescription("Buscando interfaces"))
			for i := 0; i < 100; i++ {
				bar.Add(1)
				time.Sleep(8 * time.Millisecond)
			}
			interfacesv4.PrintList(true)
			println("Presione Enter para continuar...")
			reader.ReadString('\n')
		case 2:
		ifacev4:
			for {
				menuRed()
				print(rainbow.Color(fmt.Sprintf("root🛜%s: ", getHostname())))
				optStr, _ := reader.ReadString('\n')
				optStr = strings.TrimSpace(optStr)
				opt, _ := strconv.Atoi(optStr)
				switch opt {
				case 1:
					var ip, mask, gw, dns1, dns2 string
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración estática"), "🌐")
					println()
					print("¿Mostrar todas las interfaces S/n? ")
					show, _ := reader.ReadString('\n')
					show = strings.TrimSpace(show)
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración estática"), "🌐")
					println()
					if show == "S" || show == "s" || show == "" {
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
						interfaz = strings.TrimSpace(strings.Split(utils.RemoveANSI(strings.Split(interfacesv4.List(true)[iface-1], "         ")[1]), "│")[0])
					} else {
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
						interfaz = strings.TrimSpace(strings.Split(utils.RemoveANSI(strings.Split(interfacesv4.List(false)[iface-1], "         ")[1]), "│")[0])
					}
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración estática"), "🌐")
					println()
					print(rainbow.Color("Dirección IP: "))
					ip, _ = reader.ReadString('\n')
					if regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`).MatchString(strings.TrimSpace(ip)) {
						sgmask := utils.SuggestMask(strings.TrimSpace(ip))
						println()
						print(rainbow.Color(fmt.Sprintf("Máscara de subred (%s) | CIDR: ", sgmask)))
						mask, _ = reader.ReadString('\n')
						if regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`).MatchString(strings.TrimSpace(mask)) {
							continue
						} else if regexp.MustCompile(`^\/([1-2][0-9]|3[0-2]|[1-9])$`).MatchString(strings.TrimSpace(mask)) {
							mask = utils.TranslateCIDR(strings.TrimSpace(mask))
						} else if strings.TrimSpace(mask) == "" {
							mask = sgmask
						} else {
							println("\033[31mMáscara de subred no válida, usando la sugerida.\033[0m")
							mask = sgmask
						}
						println()
						sggw := utils.SuggestGW(strings.TrimSpace(ip), strings.TrimSpace(mask))
						print(rainbow.Color(fmt.Sprintf("Puerta de enlace predeterminada (%s): ", sggw)))
						gw, _ = reader.ReadString('\n')
						if regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`).MatchString(strings.TrimSpace(gw)) {
							check := utils.IsInNet(strings.TrimSpace(gw), utils.CalculateHosts(ip, mask)["network"], strings.TrimSpace(mask))
							if !check {
								print("\033[31mLa puerta de enlace no pertenece a la subred, ¿estás seguro que quieres usarla?, puedes no tener conexión a internet (S/n)\033[0m")
								confirm, _ := reader.ReadString('\n')
								confirm = strings.TrimSpace(confirm)
								if confirm == "S" || confirm == "s" || confirm == "" {
									gw = strings.TrimSpace(gw)
								} else {
									println("\033[31mUsando la puerta de enlace sugerida.\033[0m")
									gw = sggw
								}
							} else if gw == strings.TrimSpace(ip) {
								println("\033[31mLa puerta de enlace no puede ser la misma que la dirección IP, usando la ultima IP\033[0m")
								gw = utils.CalculateHosts(ip, mask)["lastHost"]
							}
						}
						println()
						print(rainbow.Color("¿Configurar DNS? S/n: "))
						dns1, _ = reader.ReadString('\n')
						dns1 = strings.TrimSpace(dns1)
						if dns1 == "S" || dns1 == "s" || dns1 == "" {
							println()
							print(rainbow.Color("DNS Primario: "))
							dns1, _ = reader.ReadString('\n')
							if !regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`).MatchString(strings.TrimSpace(dns1)) {
								println("\033[31mDirección IP no válida, usando 8.8.8.8\033[0m")
								dns1 = "8.8.8.8"
							}
							println()
							print(rainbow.Color("¿Configurar DNS Secundario? S/n: "))
							dns2, _ = reader.ReadString('\n')
							dns2 = strings.TrimSpace(dns2)
							if dns2 == "S" || dns2 == "s" || dns2 == "" {
								println()
								print(rainbow.Color("DNS Secundario: "))
								dns2, _ = reader.ReadString('\n')
								if !regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`).MatchString(strings.TrimSpace(dns2)) {
									println("\033[31mDirección IP no válida, usando 1.1.1.1\033[0m")
									dns2 = "1.1.1.1"
								}
							}
						}
						interfacesv4.SetStatic(interfaz, ip, mask, gw, dns1, dns2)
						clear()
						println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración estática"), "🌐")
						println()
						interfacesv4.PrintConfig(interfacesv4.ParseConfig(interfacesv4.GetConfig(interfaz)))
						println()
						println("Presione Enter para volver al menu")
						reader.ReadString('\n')
					} else {
						println("\033[31mDirección IP no válida.\033[0m")
						print("Presione Enter para volver al menu")
						reader.ReadString('\n')
					}
				case 2:
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración dinámica"), "🌐")
					println()
					print("¿Mostrar todas las interfaces S/n? ")
					show, _ := reader.ReadString('\n')
					show = strings.TrimSpace(show)
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración dinámica"), "🌐")
					println()
					if show == "S" || show == "s" || show == "" {
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
						interfaz = strings.TrimSpace(strings.Split(utils.RemoveANSI(strings.Split(interfacesv4.List(true)[iface-1], "         ")[1]), "│")[0])
					} else {
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
						interfaz = strings.TrimSpace(strings.Split(utils.RemoveANSI(strings.Split(interfacesv4.List(false)[iface-1], "         ")[1]), "│")[0])
					}
					interfacesv4.SetDinamic(interfaz)
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración dinámica"), "🌐")
					println()
					println(rainbow.Color("Configuración dinámica establecida correctamente"))
					println()
					interfacesv4.PrintConfig(interfacesv4.ParseConfig(interfacesv4.GetConfig(interfaz)))
					println()
					println("Presione Enter para volver al menu")
					reader.ReadString('\n')

				case 3:
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración actual"), "🌐")
					println()
					print("¿Mostrar todas las interfaces S/n? ")
					show, _ := reader.ReadString('\n')
					show = strings.TrimSpace(show)
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración actual"), "🌐")
					println()
					if show == "S" || show == "s" || show == "" {
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

						interfaz = strings.TrimSpace(strings.Split(utils.RemoveANSI(strings.Split(interfacesv4.List(true)[iface-1], "         ")[1]), "│")[0])

					} else {
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
						interfaz = strings.TrimSpace(strings.Split(utils.RemoveANSI(strings.Split(interfacesv4.List(false)[iface-1], "         ")[1]), "│")[0])
					}
					config := interfacesv4.ParseConfig(interfacesv4.GetConfig(interfaz))
					clear()
					println("💻|", rainbow.Color("Configuración de Red ➡️ Configuración de Interfaz | "+interfaz), "🌐")
					println()
					interfacesv4.PrintConfig(config)
					println()
					println("Presione Enter para volver al menu")
					reader.ReadString('\n')
				case 4:
					break ifacev4
				}
			}
		case 4:
			return
		default:
			println("Opción no válida")
		}
	}
}
