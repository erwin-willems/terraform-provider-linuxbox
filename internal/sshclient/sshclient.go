package sshclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	serrors "errors"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type Config struct {
	Host             string
	Port             int
	User             string
	Password         string
	PrivateKey       string
	PrivateKeyPhrase string
	PrivEscPassword  string
	UseSSHAgent      bool
	PrivEscMethod    string
}

var SessionLimit = 5

type sshClient struct {
	*ssh.Client
	sessionsInUse int
	mu            *sync.Mutex
	cond          *sync.Cond
}

func newSSHClient(sc *ssh.Client) *sshClient {
	mu := new(sync.Mutex)
	return &sshClient{
		Client:        sc,
		sessionsInUse: 0,
		mu:            mu,
		cond:          sync.NewCond(mu),
	}
}

type sshSession struct {
	*ssh.Session
	cl *sshClient
}

func (s *sshClient) NewSession() (*sshSession, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	for s.sessionsInUse >= SessionLimit {
		s.cond.Wait()
	}

	cs, err := s.Client.NewSession()
	if err != nil {
		return nil, err
	}

	s.sessionsInUse++
	return &sshSession{
		Session: cs,
		cl:      s,
	}, nil
}

func (s *sshSession) Close() error {
	defer func() {
		s.cl.mu.Lock()
		s.cl.sessionsInUse--
		s.cl.cond.Broadcast()
		s.cl.mu.Unlock()
	}()
	return s.Session.Close()
}

// ErrTimeout is returned when there was a timeout when connecting to SSH daemon.
var ErrTimeout = serrors.New("timed out connecting to ssh daemon")

func newClientFuture() *clientFuture {
	mu := new(sync.Mutex)
	return &clientFuture{
		mu:  mu,
		cnd: sync.NewCond(mu),
	}
}

type clientFuture struct {
	client *sshClient
	err    error
	mu     *sync.Mutex
	cnd    *sync.Cond
}

func (cf *clientFuture) getClient() (*sshClient, error) {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	for cf.client == nil && cf.err == nil {
		cf.cnd.Wait()
	}
	return cf.client, cf.err
}

func (cf *clientFuture) createClientInternal(c *Config) (*sshClient, error) {
	var auths []ssh.AuthMethod
	var sshAgent net.Conn
	keys := []ssh.Signer{}

	key, err := ioutil.ReadFile(c.PrivateKey)
	if err != nil {
		return nil, errors.Wrap(err, "while reading private ssh_key")
	}

	if c.Password != "" {
		auths = append(auths, ssh.Password(c.Password))
	}

	ssh_auth_sock := os.Getenv("SSH_AUTH_SOCK")
	if c.UseSSHAgent && ssh_auth_sock != "" {
		if sshAgent, err = net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
			signers, err := agent.NewClient(sshAgent).Signers()
			if err == nil {
				keys = append(keys, signers...)
			}
		}
	}

	if !c.UseSSHAgent {
		if c.PrivateKeyPhrase != "" {
			signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(c.PrivateKeyPhrase))
			if err == nil {
				keys = append(keys, signer)
			}
		} else {
			signer, err := ssh.ParsePrivateKey(key)
			if err == nil {
				keys = append(keys, signer)
			}
		}
	}
	auths = append(auths, ssh.PublicKeys(keys...))

	server := c.Host
	port := c.Port
	addr := fmt.Sprintf("%s:%d", server, port)

	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         15 * time.Second,
	}

	var client *ssh.Client
	deadline := time.Now().Add(time.Minute)
	for {
		client, err = ssh.Dial("tcp", addr, config)
		if err == nil {
			break
		}

		if IsConnectTimeout(err) {
			return nil, ErrTimeout
		}

		if time.Now().Before(deadline) {
			time.Sleep(1 * time.Second)
			continue
		}

		return nil, err
	}

	return newSSHClient(client), nil

}

func IsConnectTimeout(err error) bool {
	if err == nil {
		return false
	}

	if err == ErrTimeout {
		return true
	}

	msg := err.Error()
	return strings.Contains(msg, "timed out while connecting to ssh")

}

func (cf *clientFuture) createClient(c *Config) {
	cl, err := cf.createClientInternal(c)
	cf.mu.Lock()
	cf.client = cl
	cf.err = err
	cf.cnd.Broadcast()
	cf.mu.Unlock()
}

var clientPool = map[*Config]*clientFuture{}
var clientPoolMu = new(sync.Mutex)

func getClient(c *Config) (*sshClient, error) {

	didCreateFuture := false

	clientPoolMu.Lock()
	cf := clientPool[c]
	if cf == nil {
		didCreateFuture = true
		cf = newClientFuture()
		clientPool[c] = cf
	}
	clientPoolMu.Unlock()

	if didCreateFuture {
		cf.createClient(c)
	}

	return cf.getClient()
}

func Run(c *Config, sudo bool, cmd string, stdin string) (string, string, error) {
	if sudo && c.PrivEscMethod == "sudo" {
		if c.PrivEscPassword == "" {
			cmd = fmt.Sprintf("sudo %s", cmd)
		} else {
			cmd = fmt.Sprintf("echo '%s' | sudo --prompt='' --stdin sh -c \"%s\"", c.PrivEscPassword, cmd)
		}
	}

	cl, err := getClient(c)
	if err != nil {
		return "", "", err
	}
	session, err := cl.NewSession()
	if err != nil {
		return "", "", errors.Wrap(err, "while open ssh session")
	}
	defer session.Close()

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	stdinPipe, err := session.StdinPipe()
	if err != nil {
		return "", "", errors.Wrap(err, "Unable to setup stdin for session")
	}
	defer stdinPipe.Close()
	if stdin != "" {
		stdinPipe.Write([]byte(stdin))
		stdinPipe.Close()
	}

	session.Stdout = stdout
	session.Stderr = stderr
	fmt.Printf("Run command: %s \n", cmd)

	err = session.Run(cmd)
	return stdout.String(), stderr.String(), err

}

func Check(c *Config) error {
	_, err := getClient(c)
	return err
}

func IsExecError(err error) bool {
	if err == nil {
		return false
	}

	msg := err.Error()
	return strings.Contains(msg, "Process exited with status")
}
