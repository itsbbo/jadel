package app

import (
	"io"
	"net"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type SSHStream struct {
	io.Reader
	io.Writer
	io.Closer
}

func ConnectSSH(user, host string, port int, privateKey string) (*ssh.Session, io.ReadWriteCloser, error) {
	authMethod, err := parsePrivateKey(privateKey)
	if err != nil {
		return nil, nil, err
	}

	config := ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := net.JoinHostPort(host, strconv.Itoa(port))
	client, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm", 80, 24, modes); err != nil {
		return nil, nil, err
	}

	sshStream, err := session.StdinPipe()
	if err != nil {
		return nil, nil, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := session.Shell(); err != nil {
		return nil, nil, err
	}

	return session, SSHStream{
		Reader: stdout,
		Writer: sshStream,
		Closer: session,
	}, nil
}

func parsePrivateKey(key string) (ssh.AuthMethod, error) {
	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}
