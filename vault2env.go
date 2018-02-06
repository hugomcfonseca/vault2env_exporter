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

type responseAPI struct {
	RequestID string `json:",omitempty"`
}

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

	data, err := readFromVault(*vaultSecret)

	_ = setEnvironmentVariable(data)
}

// readFromVault function in progress...
func readFromVault(path string) (*vault.Secret, error) {
	var err error
	secret := new(vault.Secret)

	vaultRequest := vaultClient.Logical()
	secret, err = vaultRequest.Read(path)

	if err != nil {
		log.Fatal(err)
	}

	if secret == nil {
		secret, err = vaultRequest.List(path)
	}

	log.Printf("%#v", secret)

	return secret, nil
}

func setEnvironmentVariable(secret *vault.Secret) bool {
	for key, value := range secret.Data {
		switch val := value.(type) {
		case string:
			key = convertVaultToEnv(key)
			_, _ = exportToEnvironment(key, fmt.Sprintf("%v", value))
			fmt.Println(key, "is string", fmt.Sprintf("%v", value))
		case []interface{}:
			for index, value2 := range val {
				fmt.Println(index, value2)
			}
		default:
			fmt.Println(key, "is of unknown type")
		}
	}

	return true
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
	if len(strings.TrimSpace(*envPrefix)) > 0 {
		envName = *envPrefix + "_" + envName
	}

	fmt.Printf("Setting up environment variable '%s'", envName)

	if err := os.Setenv(envName, envValue); err != nil {
		return false, err
	}

	// @todo: set up a shell script here appending all environment variables and
	// source it
	// check it here http://craigwickesser.com/2015/02/golang-cmd-with-custom-environment/

	return true, nil
}
