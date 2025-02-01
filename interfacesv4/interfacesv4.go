package interfacesv4

import (
	"configred/rainbow"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func RemoveANSI(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}

func List(everything bool) []string {
	var ifaces []string
	cmd := exec.Command("netsh", "interface", "show", "interface")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil
	}
	for _, iface := range strings.Split(out.String(), "\n") {
		if everything {
			if iface == "" || len(iface) < 2 || !strings.Contains(iface, "Dedicado") {
				continue
			}
			ifaces = append(ifaces, iface)
		} else {
			if iface == "" || len(iface) < 2 || !strings.Contains(iface, "Dedicado") || strings.Contains(iface, "Desconectado") {
				continue
			}
			ifaces = append(ifaces, iface)
		}
	}
	return ifaces
}

func PrintList(all bool) {
	ifaces := ParseList(List(all), all)
	fmt.Println("➡ ", rainbow.Color("Lista De Interfaces De Red"), "🌐")

	longest := 0
	for _, iface := range ifaces {
		length := len(RemoveANSI(iface))
		if length > longest {
			longest = length
		}
	}

	longest -= 6
	tline, bline, bbline := "┌"+strings.Repeat("─", longest)+"┐", "└"+strings.Repeat("─", longest)+"┘", "└"+strings.Repeat("─", longest)+"┘"

	println(tline)
	cabecera := " Nº Estado Adm.    Estado\t    Tipo\t    Nombre" + strings.Repeat(" ", longest-58)
	println("│", rainbow.Color(cabecera), "│")
	println(bline)
	println(tline)
	println(strings.Join(ifaces, "\n"))
	println(bbline)
}

func ParseList(ifaces []string, everything bool) []string {
	var parsed []string
	if everything {
		for i, iface := range ifaces {
			i += 1
			if iface == "" || len(iface) < 2 || !strings.Contains(iface, "Dedicado") {
				continue
			}
			coloredIface := rainbow.Color(strings.TrimSpace(iface))
			if strings.Contains(iface, "Desconectado") {
				coloredIface = strings.Replace(coloredIface, "Desconectado", "\033[31mDesconectado\033[0m", 1)
				parsed = append(parsed, fmt.Sprintf("│\033[31m[%d]\033[0m %s", i, coloredIface))
			} else {
				parsed = append(parsed, fmt.Sprintf("│\033[32m[%d]\033[0m %s", i, coloredIface))
			}
		}

	} else {
		saltos := 0
		for i, iface := range ifaces {
			i += 1
			if iface == "" || len(iface) < 2 || !strings.Contains(iface, "Dedicado") || strings.Contains(iface, "Desconectado") {
				saltos += 1
				continue
			}
			coloredIface := rainbow.Color(strings.TrimSpace(iface))
			parsed = append(parsed, fmt.Sprintf("│\033[32m[%d]\033[0m %s", i-saltos, coloredIface))
		}
	}

	return makeBorders(parsed)
}

func makeBorders(items []string) []string {
	var longest int
	for _, line := range items {
		length := len(RemoveANSI(line))
		if strings.Contains(line, "\t") {
			length -= 4
		}
		if length > longest {
			longest = length
		}
	}

	for i, iface := range items {
		rem := longest - len(RemoveANSI(iface)) - func() int {
			if strings.Contains(iface, "\t") {
				return 4
			} else {
				return 0
			}
		}()
		if rem > 0 {
			items[i] = iface + strings.Repeat(" ", rem) + "│"
		} else {
			items[i] = iface + "│"
		}
	}
	return items
}

func GetConfig(iface string) []string {
	var config []string

	cmd := exec.Command("netsh", "interface", "ip", "show", "config", iface)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return []string{"Error al obtener la configuración de la interfaz."}
	}
	for _, line := range strings.Split(out.String(), "\n") {
		if line == "" || len(line) <= 1 || strings.Contains(line, iface) {
			continue
		} else {
			line = strings.TrimSpace(line)
			config = append(config, line)
		}
	}
	return config
}

func ParseConfig(config []string) []string {
	var parsed []string
	var value string
	for i, line := range config {
		if strings.Contains(line, "DHCP habilitado") {
			value = strings.TrimSpace(strings.Split(line, ":")[1])
			if strings.Contains(value, "S") {
				parsed = append(parsed, "│"+rainbow.Color(" DHCP habilitado ")+"Si")
			} else {
				parsed = append(parsed, "│"+rainbow.Color(" DHCP deshabilitado ")+"No")
			}
		}
		if strings.Contains(line, "IP") {
			value = strings.TrimSpace(strings.Split(line, ":")[1])
			parsed = append(parsed, "│"+rainbow.Color(" Direccion IP --> ")+value)
		}
		if strings.Contains(line, "Prefijo") {
			value = strings.Split(strings.Split(strings.TrimSpace(strings.Split(line, ":")[1]), "/")[1], " ")[0]
			parsed = append(parsed, "│"+rainbow.Color(" Prefijo de subred (CIDR) --> ")+"/"+value)
			value = strings.Split(strings.Split(strings.TrimSpace(strings.Split(line, "(")[1]), ")")[0], " ")[1]
			parsed = append(parsed, "│\t└"+rainbow.Color(" Mascara de subred --> ")+value)
		}
		if strings.Contains(line, "Puerta") {
			value = strings.TrimSpace(strings.Split(line, ":")[1])
			parsed = append(parsed, "│"+rainbow.Color(" Puerta de enlace predeterminada --> ")+value)
		}
		if strings.Contains(line, "DNS") {
			value = strings.TrimSpace(strings.Split(line, ":")[1])
			parsed = append(parsed, "│"+rainbow.Color(" Servidores DNS "))
			parsed = append(parsed, "│\t├"+rainbow.Color(" DNS Primario --> ")+value)
			if regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`).MatchString(strings.TrimSpace(config[i+1])) {
				parsed = append(parsed, "│\t└"+rainbow.Color(" DNS Secundario --> ")+strings.TrimSpace(config[i+1]))
			} else {
				parsed = append(parsed, "│\t└"+rainbow.Color(" DNS Secundario --> ")+"Sin configurar")
			}
		}
		if strings.Contains(line, "WINS") {
			value = strings.TrimSpace(strings.Split(line, ":")[1])
			parsed = append(parsed, "│"+rainbow.Color(" Servidores WINS configurados --> ")+value)
		}
	}
	return makeBorders(parsed)
}

func PrintConfig(config []string) {
	var longest int
	for _, line := range config {
		length := len(RemoveANSI(line))
		if strings.Contains(line, "\t") {
			length -= 4
		}
		if length > longest {
			longest = length
		}
	}

	longest -= 6

	tline := "┌" + strings.Repeat("─", longest) + "┐"
	bline := "└" + strings.Repeat("─", longest) + "┘"

	println(tline)
	println(strings.Join(config, "\n"))
	println(bline)
}
