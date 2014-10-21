package util

import (
	"fmt"
	"log"
	"os/exec"
)

// 公司告警系统上报
func AgentWarn(msg string) {
	cmdstr := fmt.Sprintf("/usr/local/agenttools/agent/agentRepStr 14982 \"%s\"", msg)
	if err := exec.Command("bash", "-c", cmdstr).Run(); err != nil {
		log.Printf("WARN: AgentWarn fail, %s", err)
	}
}
