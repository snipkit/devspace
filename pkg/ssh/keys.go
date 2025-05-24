package ssh

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"path/filepath"
	"sync"

	"dev.khulnasoft.com/pkg/provider"
	"dev.khulnasoft.com/pkg/util"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

var (
	DevSpaceSSHHostKeyFile    = "id_devspace_rsa_host"
	DevSpaceSSHPrivateKeyFile = "id_devspace_rsa"
	DevSpaceSSHPublicKeyFile  = "id_devspace_rsa.pub"
)

var keyLock sync.Mutex

func rsaKeyGen() (privateKey string, publicKey string, err error) {
	privateKeyRaw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", errors.Errorf("generate private key: %v", err)
	}

	return generateKeys(pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKeyRaw),
	}, privateKeyRaw)
}

func generateKeys(block pem.Block, cp crypto.Signer) (privateKey string, publicKey string, err error) {
	pkBytes := pem.EncodeToMemory(&block)
	privateKey = string(pkBytes)

	publicKeyRaw := cp.Public()
	p, err := ssh.NewPublicKey(publicKeyRaw)
	if err != nil {
		return "", "", err
	}
	publicKey = string(ssh.MarshalAuthorizedKey(p))

	return privateKey, publicKey, nil
}

func makeHostKey() (string, error) {
	privKey, _, err := rsaKeyGen()
	if err != nil {
		return "", err
	}

	return privKey, err
}

func makeSSHKeyPair() (string, string, error) {
	privKey, pubKey, err := rsaKeyGen()
	if err != nil {
		return "", "", err
	}

	return pubKey, privKey, err
}

func GetPrivateKeyRaw(context, workspaceID string) ([]byte, error) {
	workspaceDir, err := provider.GetWorkspaceDir(context, workspaceID)
	if err != nil {
		return nil, err
	}

	return GetPrivateKeyRawBase(workspaceDir)
}

func GetDevSpaceKeysDir() string {
	dir, err := util.UserHomeDir()
	if err == nil {
		tempDir := filepath.Join(dir, ".devspace", "keys")
		err = os.MkdirAll(tempDir, 0755)
		if err == nil {
			return tempDir
		}
	}

	tempDir := os.TempDir()
	return filepath.Join(tempDir, "devspace-ssh")
}

func GetDevSpaceHostKey() (string, error) {
	tempDir := GetDevSpaceKeysDir()
	return GetHostKeyBase(tempDir)
}

func GetDevSpacePublicKey() (string, error) {
	tempDir := GetDevSpaceKeysDir()
	return GetPublicKeyBase(tempDir)
}

func GetDevSpacePrivateKeyRaw() ([]byte, error) {
	tempDir := GetDevSpaceKeysDir()
	return GetPrivateKeyRawBase(tempDir)
}

func GetHostKey(context, workspaceID string) (string, error) {
	workspaceDir, err := provider.GetWorkspaceDir(context, workspaceID)
	if err != nil {
		return "", err
	}

	return GetHostKeyBase(workspaceDir)
}

func GetPrivateKeyRawBase(dir string) ([]byte, error) {
	keyLock.Lock()
	defer keyLock.Unlock()

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	// check if key pair exists
	privateKeyFile := filepath.Join(dir, DevSpaceSSHPrivateKeyFile)
	publicKeyFile := filepath.Join(dir, DevSpaceSSHPublicKeyFile)
	_, err = os.Stat(privateKeyFile)
	if err != nil {
		pubKey, privateKey, err := makeSSHKeyPair()
		if err != nil {
			return nil, errors.Wrap(err, "generate key pair")
		}

		err = os.WriteFile(publicKeyFile, []byte(pubKey), 0644)
		if err != nil {
			return nil, errors.Wrap(err, "write public ssh key")
		}

		err = os.WriteFile(privateKeyFile, []byte(privateKey), 0600)
		if err != nil {
			return nil, errors.Wrap(err, "write private ssh key")
		}
	}

	// read private key
	out, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "read private ssh key")
	}

	return out, nil
}

func GetHostKeyBase(dir string) (string, error) {
	keyLock.Lock()
	defer keyLock.Unlock()

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	// check if key pair exists
	hostKeyFile := filepath.Join(dir, DevSpaceSSHHostKeyFile)
	_, err = os.Stat(hostKeyFile)
	if err != nil {
		privateKey, err := makeHostKey()
		if err != nil {
			return "", errors.Wrap(err, "generate host key")
		}

		err = os.WriteFile(hostKeyFile, []byte(privateKey), 0600)
		if err != nil {
			return "", errors.Wrap(err, "write host key")
		}
	}

	// read public key
	out, err := os.ReadFile(hostKeyFile)
	if err != nil {
		return "", errors.Wrap(err, "read host ssh key")
	}

	return base64.StdEncoding.EncodeToString(out), nil
}

func GetPublicKeyBase(dir string) (string, error) {
	keyLock.Lock()
	defer keyLock.Unlock()

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return "", err
	}

	// check if key pair exists
	privateKeyFile := filepath.Join(dir, DevSpaceSSHPrivateKeyFile)
	publicKeyFile := filepath.Join(dir, DevSpaceSSHPublicKeyFile)
	_, err = os.Stat(privateKeyFile)
	if err != nil {
		pubKey, privateKey, err := makeSSHKeyPair()
		if err != nil {
			return "", errors.Wrap(err, "generate key pair")
		}

		err = os.WriteFile(publicKeyFile, []byte(pubKey), 0644)
		if err != nil {
			return "", errors.Wrap(err, "write public ssh key")
		}

		err = os.WriteFile(privateKeyFile, []byte(privateKey), 0600)
		if err != nil {
			return "", errors.Wrap(err, "write private ssh key")
		}
	}

	// read public key
	out, err := os.ReadFile(publicKeyFile)
	if err != nil {
		return "", errors.Wrap(err, "read public ssh key")
	}

	return base64.StdEncoding.EncodeToString(out), nil
}

func GetPublicKey(context, workspaceID string) (string, error) {
	workspaceDir, err := provider.GetWorkspaceDir(context, workspaceID)
	if err != nil {
		return "", err
	}

	return GetPublicKeyBase(workspaceDir)
}
