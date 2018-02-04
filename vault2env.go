package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

const (
	version = "0.0.1-alpha"
)

var (
	vaultAddr   = flag.String("address", "http://localhost:8200", "Vault URL")
	vaultToken  = flag.String("token", "", "Vault token")
	vaultPrefix = flag.String("prefix", "", "Vault prefix")
	vaultSecret = flag.String("secret", "", "Variable or path of variables")
	isRecursive = flag.Bool("recursive", false, "Exports, or not, all variables inside Vault secret")
	envPrefix   = flag.String("env-prefix", "VAULT_", "Sets a prefix to exported variables")

	vaultClient *vault.Client
)

func main() {
	var err error

	flag.Parse()

	if len(strings.TrimSpace(*vaultToken)) == 0 {
		fmt.Print("No valid token.")
		return
	}

	*vaultSecret = setSecretPath()

	config := &vault.Config{
		Address: *vaultAddr,
	}

	if vaultClient, err = vault.NewClient(config); err != nil {
		fmt.Print("Error creating Vault client...")
		return
	}

	vaultClient.SetToken(*vaultToken)

	_, err = readFromVault(*vaultSecret, *isRecursive)
}

// readFromVault function in progress...
func readFromVault(path string, recursive bool) (bool, error) {
	vault := vaultClient.Logical()
	secrets, err := vault.Read(path)

	if err != nil {
		log.Fatal(err)
	}

	if secrets == nil {
		secrets, err = vault.List(path)
	}

	log.Printf("%#v", *secrets)

	return true, nil
}

func setSecretPath() string {
	var secretPath string

	if len(strings.Trim(*vaultPrefix, " ")) > 0 {
		secretPath = "/secret" + *vaultPrefix + "/" + *vaultSecret
	} else {
		secretPath = "/secret" + "/" + *vaultSecret
	}

	return secretPath
}

// convertVaultToEnv converts a Vault variable to an environment variable
func convertVaultToEnv(variable string) string {
	return strings.Trim(strings.ToUpper(strings.Replace(variable, "/", "_", -1)), "_")
}

// exportToEnvironment exports a given variable to environment
func exportToEnvironment(envName string, envValue string) (bool, error) {
	fmt.Printf("Setting up environment variable '%s'", envName)

	if err := os.Setenv(envName, envValue); err != nil {
		return false, err
	}

	return true, nil
}
