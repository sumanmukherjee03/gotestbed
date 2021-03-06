package sshutils

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

var (
	stdin  bytes.Buffer
	stdout bytes.Buffer
	stdio  bytes.Buffer
)

// NewSSHConnectCmd - is a function that generates a command to create a new connection and assign a pty
func NewSSHConnectCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "connect SERVER PORT",
		Short: "Connects to a remote server",
		Long: `Lets you connect to a remote server.
			And assigns a pseudo terminal`,
		Example: `
			### Available commands for ssh subcommands
			gotils ssh connect SERVER`,
		PreRun: func(cmd *cobra.Command, args []string) {
			if len(ssh_username) == 0 {
				ssh_username = viper.GetString("ssh.username")
			}
			if len(ssh_private_key_path) == 0 {
				ssh_private_key_path = viper.GetString("ssh.private_key_path")
			}
			if len(ssh_public_key_path) == 0 {
				ssh_public_key_path = viper.GetString("ssh.public_key_path")
			}
			if len(ssh_config_path) == 0 {
				ssh_config_path = viper.GetString("ssh.config_path")
			}
			if ssh_with_ssm && len(ssh_aws_token_serial_number) == 0 {
				ssh_aws_token_serial_number = viper.GetString("ssh.aws.token_serial_number")
			}
		},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return utils.RaiseCmdErr(cmd, "server needs to be provided")
			}
			ssh_host = args[0]
			if len(args) == 1 {
				ssh_port = 22
			}
			if len(args) == 2 {
				port, err := strconv.Atoi(args[1])
				if err != nil {
					return utils.RaiseCmdErr(cmd, "port needs to be an integer")
				}
				ssh_port = port
			}
			if len(args) > 2 {
				return utils.RaiseCmdErr(cmd, "Too many args")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if ssh_with_ssm {
				awsConn, err := NewAwsConn()
				if err != nil {
					utils.CheckErr(fmt.Sprintf("generating new aws session failed", err))
				}
				cmd, err := NewSsmShellConn(awsConn, true)
				if err != nil {
					utils.CheckErr(fmt.Sprintf("establishing ssm connection or generating ssh command failed", err))
				}
				if err = cmd.Run(); err != nil {
					utils.CheckErr(fmt.Sprintf("running shell for ssh command failed", err))
				}
			} else {
				r, err := NewSshShellConn(true)
				if err != nil {
					utils.CheckErr(fmt.Sprintf("assigning pseudo terminal failed: %s", err))
				}
				r.StartInteractiveShell()
			}
		},
	}
	return cmd
}
