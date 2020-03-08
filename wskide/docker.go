package wskide

import (
	"fmt"
	"strings"
)

func dockerVersion() (string, error) {
	return SysErr("@docker version --format {{.Server.Version}}")
}

func dockerStatus(container string) {
	res, err := SysErr("@docker inspect --format {{.State.Status}} " + container)
	if err != nil {
		res = "not running\n"
	}
	fmt.Print(container, ": ", res)
}

func dockerCreateNetwork(network string) error {
	_, err := SysErr("@docker network inspect --format='{{.Driver}}' " + network)
	if err != nil {
		_, err = SysErr("@docker network create " + network)
		fmt.Printf("Network %s created\n", network)
	} else {
		fmt.Printf("Network %s exists\n", network)
	}
	return err
}

func dockerDeleteNetwork(network string) error {
	res, err := SysErr("@docker network inspect --format={{.Driver}} " + network)
	if strings.TrimSpace(res) == "bridge" {
		fmt.Printf("Network %s deleted\n", network)
		_, err = SysErr("@docker network rm " + network)
	}
	return err
}
