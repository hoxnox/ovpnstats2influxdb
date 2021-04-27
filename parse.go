package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const splitCharacter = ","

// ClientInfo represents a CLIENT_LIST entry
// HEADER,CLIENT_LIST,Common Name,Real Address,Virtual Address,Virtual IPv6 Address,Bytes Received,Bytes Sent,Connected Since,Connected Since (time_t),Username,Client ID,Peer ID
type ClientInfo struct {
	CommonName     string
	RealAddress    string
	BytesReceived  int
	BytesSent      int
	ConnectedSince time.Time
}

// RoutingInfo represents a ROUTING_TABLE entry
// HEADER,ROUTING_TABLE,Virtual Address,Common Name,Real Address,Last Ref,Last Ref (time_t)
type RoutingInfo struct {
	VirtualAddress string
	CommonName     string
	RealAddress    string
	LastRef        time.Time
}

type CurrentParsingType int

const (
	NoneType CurrentParsingType = iota
	ClientInfoType
	RoutingInfoType
)

func parseClientListEntry(line string) (ClientInfo, error) {
	parts := strings.Split(line, splitCharacter)
	bytesReceived, err := strconv.Atoi(parts[2])
	if err != nil {
		return ClientInfo{}, err
	}
	bytesSent, err := strconv.Atoi(parts[3])
	if err != nil {
		return ClientInfo{}, err
	}
	connectedSince, err := time.Parse(time.ANSIC, parts[4])
	if err != nil {
		return ClientInfo{}, err
	}
	info := ClientInfo{
		CommonName:     parts[0],
		RealAddress:    parts[1],
		BytesReceived:  bytesReceived,
		BytesSent:      bytesSent,
		ConnectedSince: connectedSince,
	}
	return info, nil
}

func parseRoutingTableEntry(line string) (RoutingInfo, error) {
	parts := strings.Split(line, splitCharacter)
	lastRef, err := time.Parse(time.ANSIC, parts[3])
	if err != nil {
		return RoutingInfo{}, err
	}
	info := RoutingInfo{
		VirtualAddress: parts[0],
		CommonName:     parts[1],
		RealAddress:    parts[2],
		LastRef:        lastRef,
	}
	return info, nil
}

// ParseStatusFile parses the openvpn-status.log file at `filename` and returns a corresponding slice of ClientInfo and RoutingInfo objects
func ParseStatusFile(filename string) ([]ClientInfo, []RoutingInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("Cannot close status file ", err)
		}
	}(file)

	var clients []ClientInfo
	var routes []RoutingInfo

	var currentParsingType = NoneType
	var parsedHeader = false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "GLOBAL STATS") && !strings.Contains(line, ",") {
			break
		}
		if strings.Contains(line, "CLIENT LIST") && !strings.Contains(line, ",") {
			currentParsingType = ClientInfoType
			parsedHeader = false
			continue
		}
		if line == "Common Name,Real Address,Bytes Received,Bytes Sent,Connected Since" {
			parsedHeader = true
			continue
		}
		if strings.Contains(line, "ROUTING TABLE") && !strings.Contains(line, ",") {
			currentParsingType = RoutingInfoType
			parsedHeader = false
			continue
		}
		if line == "Virtual Address,Common Name,Real Address,Last Ref" {
			parsedHeader = true
			continue
		}
		if currentParsingType == ClientInfoType && parsedHeader {
			info, err := parseClientListEntry(line)
			if err != nil {
				return nil, nil, err
			}
			clients = append(clients, info)
			continue
		}
		if currentParsingType == RoutingInfoType && parsedHeader {
			info, err := parseRoutingTableEntry(line)
			if err != nil {
				return nil, nil, err
			}
			routes = append(routes, info)
			continue
		}
	}
	return clients, routes, nil
}
