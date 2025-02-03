package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func RemoveANSI(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return re.ReplaceAllString(input, "")
}

func MakeBorders(items []string) []string {
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

func convertToBinary(ip string) string {
	var binary string
	for _, octet := range strings.Split(ip, ".") {
		octetInt, _ := strconv.Atoi(octet)
		binary += fmt.Sprintf("%08b", octetInt)
	}
	return binary
}

func convertToIp(binary string) string {
	var ip string
	for i := 0; i < 32; i += 8 {
		octet, _ := strconv.ParseInt(binary[i:i+8], 2, 64)
		ip += fmt.Sprintf("%d.", octet)
	}
	return strings.TrimSuffix(ip, ".")
}

func TranslateCIDR(mask string) string {
	var decimalMask []string
	cidr, _ := strconv.Atoi(strings.TrimPrefix(mask, "/"))
	binaryMask := strings.Repeat("1", cidr) + strings.Repeat("0", 32-cidr)

	for i := 0; i < 32; i += 8 {
		octet, _ := strconv.ParseInt(binaryMask[i:i+8], 2, 64)
		decimalMask = append(decimalMask, fmt.Sprintf("%d", octet))
	}

	return strings.Join(decimalMask, ".")

}

func CalculateHosts(ip string, mask string) map[string]string {
	binaryIP := convertToBinary(ip)
	binaryMask := convertToBinary(mask)
	netBits := strings.Count(binaryMask, "1")
	network := binaryIP[:netBits]
	hostBits := 32 - netBits
	maxHostNet := 2 ^ hostBits - 2
	broadcast := network + strings.Repeat("1", hostBits)
	firstHost := network + strings.Repeat("0", hostBits-1) + "1"
	lastHost := network + strings.Repeat("1", hostBits-1) + "0"

	return map[string]string{
		"network":    convertToIp(network + strings.Repeat("0", 32-netBits)),
		"firstHost":  convertToIp(firstHost),
		"lastHost":   convertToIp(lastHost),
		"broadcast":  convertToIp(broadcast),
		"totalHosts": strconv.Itoa(maxHostNet),
	}
}

func SuggestMask(ip string) string {
	if strings.HasPrefix(ip, "10.") {
		return "255.0.0.0 | /8"
	} else if strings.HasPrefix(ip, "172.") {
		return "255.255.0.0 | /16"
	}
	return "255.255.255.0 | /24"
}

func SuggestGW(ip string, mask string) string {
	data := CalculateHosts(ip, mask)
	if ip == data["firstHost"] {
		return data["lastHost"]
	} else if ip == data["lastHost"] {
		return data["firstHost"]
	}
	return data["firstHost"] + " | " + data["lastHost"]
}

func IsInNet(ip string, networkIp string, mask string) bool {
	binaryIP := convertToBinary(ip)
	binaryNetwork := convertToBinary(networkIp)
	binaryNetwork = binaryNetwork[:strings.Count(convertToBinary(mask), "1")]
	binaryMask := convertToBinary(mask)
	netBits := strings.Count(binaryMask, "1")
	networkIP := binaryIP[:netBits]
	if networkIP != binaryNetwork {
		return false
	} else {
		firstHost := binaryNetwork + strings.Repeat("0", 32-netBits-1) + "1"
		lastHost := binaryNetwork + strings.Repeat("1", 32-netBits-1) + "0"
		if binaryIP >= firstHost && binaryIP <= lastHost {
			return true
		} else {
			return false
		}
	}
}
