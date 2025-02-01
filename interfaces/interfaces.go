package interfaces

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

func makeBorders(ifaces []string) []string {
	var longest int
	for _, iface := range ifaces {
		length := len(RemoveANSI(iface))
		if length > longest {
			longest = length
		}
	}

	for i, iface := range ifaces {
		rem := longest - len(RemoveANSI(iface))
		if rem > 0 {
			ifaces[i] = iface + strings.Repeat(" ", rem) + "│"
		} else {
			ifaces[i] = iface + "│"
		}
	}
	return ifaces
}

func GetConfig(iface string) {
	cmd := exec.Command("netsh", "interface", "ip", "show", "config", iface)
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return
	}

	fmt.Println(out.String())
}
