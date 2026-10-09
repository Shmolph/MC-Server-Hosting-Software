package main

import (
	"fmt"
	"strings"
	"strconv"
	"os/exec"
)

// Create server =DDD
func createcontainer(s serverdiagram) error {
	name := ("mc-"+(strconv.Itoa(s.Id)))
	args := []string{
		"create",
		"--name", name,
		"-e", "EULA=TRUE",
		"-e", "TYPE=" + strings.ToUpper(s.Type),
		"-e", "VERSION=" + s.Version,
		"-e", "MEMORY=" + strconv.Itoa(s.RamMB) + "M",
		"-p", strconv.Itoa(s.Port) + ":25565",
		"-v", "/home/mc-servers/" + strconv.Itoa(s.Id) + ":/data",
		"--memory", strconv.Itoa(s.RamMB+512) + "m",
		"--cpus", "2",
		"itzg/minecraft-server",
	}

	cmd := exec.Command("docker", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("DOCKERMAN: Failed to create container:", string(out), err)
		return err
	}
	return nil
}

// Check Server Status
func servercontainerstatus(serverid int) string {
	name := ("mc-"+(strconv.Itoa(serverid)))
	args := []string{
	"inspect",
	"-f", "{{.State.Status}} {{if .State.Health}}{{.State.Health.Status}}{{end}}",
	name,
	}
	cmd := exec.Command("docker", args...)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("DOCKERMAN: Server no exist:", string(out), err)
		return "not setup"
	}
	health := ""
	fields := strings.Fields(string(out))
	if len(fields) == 0 {
		return "offline"
	}
	switch fields[0] {
	case "running":
		if len(fields) > 1 {
			health = fields[1]
		}
		if health == "" || health == "healthy" {
			return "online"
		}
		return "starting"
	default:
		return "offline"
	}
}

// Change status
func changestatus(serverid int, dockerverb string) error {
	name := "mc-" + strconv.Itoa(serverid)
	args := []string{
		dockerverb,
		name,
	}
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("DOCKERMAN: Failed to change status:", string(out), err)
		return err
	}
	return nil
}

// Terminal
func sendcommand(serverid int, command string) (string, error) {
	name := "mc-" + strconv.Itoa(serverid)
	args := []string{
		"exec",
		name,
		"rcon-cli",
		command,
	}
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("DOCKERMAN: Failed to send command", string(out), err)
		return string(out), err
	}
	return string(out), nil
}

// Logs
func getlogs(serverid int) (string, error) {
	name := "mc-" + strconv.Itoa(serverid)
	args := []string{
		"logs",
		"--tail",
		"200",
		name,
	}
	cmd := exec.Command("docker", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("DOCKERMAN: Failed to get logs", string(out), err)
		return string(out), err
	}
	return string(out), nil
}
