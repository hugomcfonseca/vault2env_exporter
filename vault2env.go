package vault2env

import (
	"flag"
	"fmt"
	"os"
	"strings"

	vault "github.com/hashicorp/vault/api"
)

const (
	version = "0.0.1-alpha"
)

var (
	vaultAddr   = flag.String("address", "http://vault:8200", "Vault URL")
	vaultToken  = flag.String("token", "", "Vault token")
	vaultPrefix = flag.String("prefix", "/prefix", "Vault prefix")
	vaultSecret = flag.String("secret", "", "Variable or path of variables")
	isRecursive = flag.Bool("recursive", false, "Exports, or not, all variables inside Vault secret")
	envPrefix   = flag.String("env-prefix", "VAULT_", "Sets a prefix to exported variables")

	vaultClient *vault.Client
)

func main() {
	var err error

	flag.Parse()

	config := &vault.Config{
		Address: *vaultAddr,
	}

	if vaultClient, err = vault.NewClient(config); err != nil {
		fmt.Printf("Error creating Vault client...")
		return
	}

	vaultClient.SetToken(*vaultToken)
}

// readFromVault ...
func readFromVault() {

}

// convertVaultToEnv Convert Vault characters
func convertVaultToEnv(path string) string {
	return strings.Trim(strings.ToUpper(strings.Replace(path, "/", "_", -1)), "_")
}

// exportToEnvironment exports a given variable to environment
func exportToEnvironment(envName string, envValue string) (bool, error) {
	fmt.Printf("Setting up environment variable '%s'", envName)

	if err := os.Setenv(envName, envValue); err != nil {
		return false, err
	}

	return true, nil
}
