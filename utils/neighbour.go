package utils

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func IsFoundHost(host string, port uint16) bool {
    target := fmt.Sprintf("%s:%d", host, port)

    _, err := net.DialTimeout("tcp", target, 1*time.Second)
    if err != nil {
        fmt.Printf("%s %v\n", target, err)
        return false
    }
    return true
}

var PATTERN = regexp.MustCompile(`^((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}`)

func FindNeighbours(myHost string, myPort uint16, startIp uint8, endIp uint8, startPort uint16, endPort uint16) []string {
    address := fmt.Sprintf("%s:%d", myHost, myPort)

    // Split the IP address into parts
    ipParts := strings.Split(myHost, ".")
    if len(ipParts) != 4 {
        fmt.Println("Invalid IP address format")
        return nil
    }

    // Extract the last octet of the IP address
    lastIp, err := strconv.Atoi(ipParts[3])
    if err != nil {
        fmt.Printf("Error converting IP part to integer: %v\n", err)
        return nil
    }

    prefixHost := strings.Join(ipParts[:3], ".") + "."

    neighbours := make([]string, 0)

    for port := startPort; port <= endPort; port++ {
        for ip := int(startIp); ip <= int(endIp); ip++ {
            newLastOctet := lastIp + ip
            if newLastOctet > 255 {
                continue // Skip invalid IP addresses
            }
            guessHost := fmt.Sprintf("%s%d", prefixHost, newLastOctet)
            guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
            if guessTarget != address && IsFoundHost(guessHost, port) {
                neighbours = append(neighbours, guessTarget)
            }
        }
    }
    return neighbours
}

func GetHost() string {
    hostname, err := os.Hostname()
    if err != nil {
        return "127.0.0.1"
    }
    address, err := net.LookupHost(hostname)
    if err != nil {
        return "127.0.0.1"
    }
    return address[0]
}
