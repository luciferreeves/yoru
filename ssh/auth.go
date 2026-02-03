package ssh

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
	"yoru/models"
	"yoru/repository"
	"yoru/types"

	"golang.org/x/crypto/ssh"
)

// LoadCredential loads the credential for a host from the database
func LoadCredential(host *models.Host) (any, error) {
	if host.CredentialID == 0 {
		return nil, errors.New("no credential configured for this host")
	}

	switch host.CredentialType {
	case types.CredentialIdentity:
		return repository.GetIdentityByID(host.CredentialID)
	case types.CredentialKey:
		return repository.GetKeyByID(host.CredentialID)
	default:
		return nil, errors.New("unknown credential type")
	}
}

// BuildSSHConfig creates an SSH client configuration from a credential
func BuildSSHConfig(credential any) (*ssh.ClientConfig, error) {
	config := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // We handle host key verification separately
		Timeout:         30 * time.Second,
	}

	switch cred := credential.(type) {
	case *models.Identity:
		config.User = cred.Username
		config.Auth = []ssh.AuthMethod{
			ssh.Password(cred.Password),
		}
		return config, nil

	case *models.Key:
		if cred.Username == "" {
			return nil, errors.New("username is required for SSH key authentication")
		}

		// Parse private key
		signer, err := ssh.ParsePrivateKey([]byte(cred.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}

		config.User = cred.Username
		config.Auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}

		// If certificate is present, add it
		if cred.Certificate != "" {
			pcert, _, _, _, err := ssh.ParseAuthorizedKey([]byte(cred.Certificate))
			if err == nil {
				if cert, ok := pcert.(*ssh.Certificate); ok {
					certSigner, err := ssh.NewCertSigner(cert, signer)
					if err == nil {
						config.Auth = append(config.Auth, ssh.PublicKeys(certSigner))
					}
				}
			}
		}

		return config, nil

	default:
		return nil, errors.New("unsupported credential type")
	}
}

// GetFingerprint calculates the SSH fingerprint from a public key
func GetFingerprint(key ssh.PublicKey) string {
	hash := sha256.Sum256(key.Marshal())
	return "SHA256:" + base64.StdEncoding.EncodeToString(hash[:])
}

// GetMD5Fingerprint calculates the MD5 fingerprint (legacy format)
func GetMD5Fingerprint(key ssh.PublicKey) string {
	hash := md5.Sum(key.Marshal())
	hexStr := hex.EncodeToString(hash[:])

	// Format as xx:xx:xx:xx:...
	formatted := ""
	for i := 0; i < len(hexStr); i += 2 {
		if i > 0 {
			formatted += ":"
		}
		formatted += hexStr[i : i+2]
	}
	return formatted
}

// VerifyHostKey checks if a host key matches a known host in the database
func VerifyHostKey(hostname string, port int, key ssh.PublicKey) (*models.KnownHost, error) {
	fingerprint := GetFingerprint(key)

	knownHost, err := repository.GetKnownHostByFingerprint(fingerprint)
	if err != nil {
		// Host not known
		return nil, err
	}

	// Verify hostname and port match
	if knownHost.Hostname != hostname || knownHost.Port != port {
		return nil, fmt.Errorf("host key fingerprint matches but hostname/port differs")
	}

	return knownHost, nil
}

// SaveHostKey saves a verified host key to the database
func SaveHostKey(hostname string, port int, key ssh.PublicKey) error {
	fingerprint := GetFingerprint(key)
	keyType := key.Type()

	knownHost := &models.KnownHost{
		Hostname:    hostname,
		Port:        port,
		KeyType:     keyType,
		Fingerprint: fingerprint,
	}

	return repository.CreateKnownHost(knownHost)
}