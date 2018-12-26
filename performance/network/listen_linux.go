package network

import (
	"bytes"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func getListenPorts(ports *ListenCollection) {
	var stdout bytes.Buffer
	cmd := exec.Command("ss", "-ltn")
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return
	}

	for {
		line, err := stdout.ReadString('\n')
		if err == io.EOF {
			break
		}
		if len(line) < 6 {
			continue
		}
		if line[0] != 'L' || line[5] != 'N' {
			continue
		}

		fields := make([]string, 0)
		for _, field := range strings.Split(line, " ") {
			val := strings.TrimSpace(field)
			if len(val) < 1 {
				continue
			}
			fields = append(fields, val)
		}
		if len(fields) < 4 {
			continue
		}
		ipPort := fields[3]
		pos := strings.LastIndex(ipPort, ":")
		if pos < 1 {
			continue
		}
		ip := ipPort[0:pos]
		port := ipPort[pos+1:]
		portVal, err := strconv.Atoi(port)
		if err != nil {
			continue
		}

		listen := &Listen{
			Address:  ip,
			Port:     portVal,
			Protocol: "tcp",
		}
		*ports = append(*ports, listen)
	}

}
