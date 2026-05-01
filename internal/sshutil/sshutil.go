package sshutil

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

// KeyDir returns the directory where SSH keys for the given profile are stored:
// ~/.ssh/gixy/<profileName>/
func KeyDir(profileName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".ssh", "gixy", profileName), nil
}

// sshDir returns ~/.ssh
func sshDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	return filepath.Join(home, ".ssh"), nil
}

// GenerateKeypair generates an ed25519 SSH keypair for the given profile.
// Keys are written to ~/.ssh/gixy/<profileName>/id_ed25519{,.pub}.
// If keys already exist, the function returns without overwriting them.
func GenerateKeypair(profileName, comment string) error {
	dir, err := KeyDir(profileName)
	if err != nil {
		return err
	}

	privPath := filepath.Join(dir, "id_ed25519")
	pubPath := filepath.Join(dir, "id_ed25519.pub")

	// Idempotent: skip if keypair already exists
	if _, err := os.Stat(privPath); err == nil {
		return nil
	}

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("create key dir: %w", err)
	}

	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("generate ed25519 key: %w", err)
	}

	// Marshal private key to OpenSSH PEM format
	privPEM, err := ssh.MarshalPrivateKey(privKey, comment)
	if err != nil {
		return fmt.Errorf("marshal private key: %w", err)
	}
	privBytes := pem.EncodeToMemory(privPEM)

	if err := os.WriteFile(privPath, privBytes, 0o600); err != nil {
		return fmt.Errorf("write private key: %w", err)
	}

	// Marshal public key to OpenSSH authorized_keys format
	sshPubKey, err := ssh.NewPublicKey(pubKey)
	if err != nil {
		return fmt.Errorf("marshal public key: %w", err)
	}
	pubBytes := ssh.MarshalAuthorizedKey(sshPubKey)
	// Append comment to the public key line
	pubLine := strings.TrimRight(string(pubBytes), "\n") + " " + comment + "\n"

	if err := os.WriteFile(pubPath, []byte(pubLine), 0o644); err != nil {
		return fmt.Errorf("write public key: %w", err)
	}

	return nil
}

// ActivateKeys symlinks ~/.ssh/id_ed25519{,.pub} to the profile's keypair.
// Switching between gixy-managed profiles is always silent.
// Only prompts when an unrecognized real file or external symlink is in the way.
func ActivateKeys(profileName string) error {
	base, err := sshDir()
	if err != nil {
		return err
	}

	dir, err := KeyDir(profileName)
	if err != nil {
		return err
	}

	gixyBase, err := func() (string, error) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".ssh", "gixy"), nil
	}()
	if err != nil {
		return err
	}

	files := []string{"id_ed25519", "id_ed25519.pub"}
	for _, f := range files {
		linkPath := filepath.Join(base, f)
		targetPath := filepath.Join(dir, f)

		if err := ensureSymlink(linkPath, targetPath, gixyBase); err != nil {
			return err
		}
	}
	return nil
}

// ensureSymlink creates or updates the symlink at linkPath → targetPath.
// gixyBase is the ~/.ssh/gixy prefix used to detect gixy-managed paths.
func ensureSymlink(linkPath, targetPath, gixyBase string) error {
	info, err := os.Lstat(linkPath)
	if os.IsNotExist(err) {
		// Nothing there — create symlink directly
		return os.Symlink(targetPath, linkPath)
	}
	if err != nil {
		return fmt.Errorf("stat %s: %w", linkPath, err)
	}

	if info.Mode()&os.ModeSymlink != 0 {
		// It's a symlink — check where it points
		existing, err := os.Readlink(linkPath)
		if err != nil {
			return fmt.Errorf("readlink %s: %w", linkPath, err)
		}
		if existing == targetPath {
			// Already correct — nothing to do
			return nil
		}
		if strings.HasPrefix(existing, gixyBase+string(filepath.Separator)) || existing == gixyBase {
			// Managed by gixy — silently switch
			if err := os.Remove(linkPath); err != nil {
				return fmt.Errorf("remove old symlink %s: %w", linkPath, err)
			}
			return os.Symlink(targetPath, linkPath)
		}
		// External symlink — prompt
		if !promptOverride(fmt.Sprintf("%s is linked to an external location (%s). Override?", filepath.Base(linkPath), existing)) {
			fmt.Printf("Skipped updating %s\n", filepath.Base(linkPath))
			return nil
		}
	} else {
		// Real file — prompt
		if !promptOverride(fmt.Sprintf("%s already exists as a file. Override?", filepath.Base(linkPath))) {
			fmt.Printf("Skipped updating %s\n", filepath.Base(linkPath))
			return nil
		}
	}

	if err := os.Remove(linkPath); err != nil {
		return fmt.Errorf("remove %s: %w", linkPath, err)
	}
	return os.Symlink(targetPath, linkPath)
}

// promptOverride prints the message with a [y/N] prompt and returns true if the user confirms.
func promptOverride(message string) bool {
	fmt.Printf("%s [y/N]: ", message)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimSpace(strings.ToLower(answer))
	return answer == "y" || answer == "yes"
}
