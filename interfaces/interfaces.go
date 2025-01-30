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

func List(all bool) []string {

	var ifaces []string
	cmd := exec.Command("netsh", "interface", "show", "interface")
	var out strings.Builder
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil
	}
	longest := 0

	if all {
		for i, iface := range strings.Split(out.String(), "\n") {
			if iface == "" || len(iface) < 2 || !strings.Contains(iface, "Dedicado") {
				continue
			}
			coloredIface := rainbow.Color(strings.TrimSpace(iface))
			if strings.Contains(iface, "Desconectado") {
				coloredIface = strings.Replace(coloredIface, "Desconectado", "\033[31mDesconectado\033[0m", 1)
				ifaces = append(ifaces, fmt.Sprintf("│\033[31m[%d]\033[0m %s", i-2, coloredIface))
			} else {
				ifaces = append(ifaces, fmt.Sprintf("│\033[32m[%d]\033[0m %s", i-2, coloredIface))
			}
		}

	} else {
		for i, iface := range strings.Split(out.String(), "\n") {
			if iface == "" || len(iface) < 2 || !strings.Contains(iface, "Dedicado") || strings.Contains(iface, "Desconectado") {
				continue
			}
			coloredIface := rainbow.Color(strings.TrimSpace(iface))
			if strings.Contains(iface, "Desconectado") {
				coloredIface = strings.Replace(coloredIface, "Desconectado", "\033[31mDesconectado\033[0m", 1)
				ifaces = append(ifaces, fmt.Sprintf("│\033[31m[%d]\033[0m %s", i-2, coloredIface))
			} else {
				ifaces = append(ifaces, fmt.Sprintf("│\033[32m[%d]\033[0m %s", i-2, coloredIface))
			}
		}
	}

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
