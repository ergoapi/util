// AGPL License
// Copyright (c) 2023 ysicing <i@ysicing.me>

package exssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/ergoapi/util/zos"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const sshAuthSock = "SSH_AUTH_SOCK"

func GetFileContent(path string) ([]byte, error) {
	buff, err := os.ReadFile(zos.ExpandPath(path))
	if err != nil {
		return []byte{}, err
	}
	return buff, nil
}

// SSHPrivateKeyPath returns ssh private key content from given path.
func SSHPrivateKeyPath(sshKey string) (string, error) {
	content, err := GetFileContent(sshKey)
	if err != nil {
		return "", fmt.Errorf("error while reading SSH key file: %v", err)
	}
	return string(content), nil
}

// SSHCertificatePath returns ssh certificate key content from given path
func SSHCertificatePath(sshCertPath string) (string, error) {
	content, err := GetFileContent(sshCertPath)
	if err != nil {
		return "", fmt.Errorf("error while reading SSH certificate file: %v", err)
	}
	return string(content), nil
}

// SSHConfigOptions contains options for SSH configuration.
type SSHConfigOptions struct {
	// HostKeyCallback is the callback for verifying host keys.
	// If nil, ssh.InsecureIgnoreHostKey() will be used (NOT RECOMMENDED for production).
	HostKeyCallback ssh.HostKeyCallback
}

// SSHConfigResult contains the SSH config and resources that need to be cleaned up.
type SSHConfigResult struct {
	Config *ssh.ClientConfig
	// AgentConn is the SSH agent connection that should be closed when done.
	// May be nil if agent auth is not used.
	AgentConn io.Closer
}

// Close closes any resources associated with the SSH config.
func (r *SSHConfigResult) Close() error {
	if r.AgentConn != nil {
		return r.AgentConn.Close()
	}
	return nil
}

// GetSSHConfig generate ssh config.
// Deprecated: Use GetSSHConfigWithOptions instead for proper resource management.
func GetSSHConfig(username, sshPrivateKeyString, passphrase, sshCert string, password string, timeout time.Duration, useAgentAuth bool) (*ssh.ClientConfig, error) {
	result, err := GetSSHConfigWithOptions(username, sshPrivateKeyString, passphrase, sshCert, password, timeout, useAgentAuth, nil)
	if err != nil {
		return nil, err
	}
	// Note: AgentConn will leak if useAgentAuth is true. Use GetSSHConfigWithOptions instead.
	return result.Config, nil
}

// GetSSHConfigWithOptions generates ssh config with options and returns resources that need cleanup.
// The caller is responsible for calling result.Close() when done with the SSH connection.
func GetSSHConfigWithOptions(username, sshPrivateKeyString, passphrase, sshCert string, password string, timeout time.Duration, useAgentAuth bool, opts *SSHConfigOptions) (*SSHConfigResult, error) {
	hostKeyCallback := ssh.InsecureIgnoreHostKey()
	if opts != nil && opts.HostKeyCallback != nil {
		hostKeyCallback = opts.HostKeyCallback
	}

	config := &ssh.ClientConfig{
		User:            username,
		Timeout:         timeout,
		HostKeyCallback: hostKeyCallback,
	}

	result := &SSHConfigResult{Config: config}

	if useAgentAuth {
		if sshAgentSock := os.Getenv(sshAuthSock); sshAgentSock != "" {
			sshAgent, err := net.Dial("unix", sshAgentSock)
			if err != nil {
				return nil, fmt.Errorf("cannot connect to SSH Auth socket %q: %s", sshAgentSock, err)
			}
			result.AgentConn = sshAgent
			config.Auth = append(config.Auth, ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers))
			return result, nil
		}
	} else if sshPrivateKeyString != "" {
		var (
			signer ssh.Signer
			err    error
		)
		if passphrase != "" {
			signer, err = parsePrivateKeyWithPassphrase(sshPrivateKeyString, passphrase)
		} else {
			signer, err = parsePrivateKey(sshPrivateKeyString)
		}
		if err != nil {
			return nil, err
		}

		if len(sshCert) > 0 {
			key, _, _, _, err := ssh.ParseAuthorizedKey([]byte(sshCert))
			if err != nil {
				return nil, fmt.Errorf("unable to parse SSH certificate: %v", err)
			}

			if _, ok := key.(*ssh.Certificate); !ok {
				return nil, fmt.Errorf("unable to cast public key to SSH certificate")
			}
			signer, err = ssh.NewCertSigner(key.(*ssh.Certificate), signer)
			if err != nil {
				return nil, err
			}
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(signer))
	} else if password != "" {
		config.Auth = append(config.Auth, ssh.Password(password))
	}

	return result, nil
}

func parsePrivateKey(keyBuff string) (ssh.Signer, error) {
	return ssh.ParsePrivateKey([]byte(keyBuff))
}

func parsePrivateKeyWithPassphrase(keyBuff, passphrase string) (ssh.Signer, error) {
	return ssh.ParsePrivateKeyWithPassphrase([]byte(keyBuff), []byte(passphrase))
}
