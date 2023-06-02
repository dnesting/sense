// Package sensecli just contains some helpers used to set up binaries that need Sense credentials.
//
// It permits users to configure the Sense client using flags (--sense-email), environment variables
// (SENSE_EMAIL), or a YAML configuration file.
//
// The easiest thing to do is:
//
//	func main() {
//		configFile, flagCreds := sensecli.SetupStandardFlags()
//		flag.Parse()
//		clients, err := sensecli.CreateClients(context.Background(), configFile, flagCreds)
//		...
//	}
//
// By default, the --sense-config flag and SENSE_CONFIG environment variable overrides everything
// else by specifying the name of a YAML configuration file that contains the following format:
//
//	accounts:
//	- credentials:
//	    email:          # specify the e-mail address directly
//	    email-from:		# read the e-mail address from a file
//	    password:		# specify the password directly
//	    password-from:	# read the password from a file
//	    mfa-from:		# read the MFA code from a file
//	    mfa-command:	# read the MFA code from a command
//
// Multiple accounts can be configured in the same file.  If you specify a configuration file, flags
// and environment variables will be ignored.
//
// If no configuration file is specified, the following flags will be used first to configure a single Sense client:
//
//	--sense-email
//	--sense-email-from
//	--sense-password
//	--sense-password-from
//	--sense-mfa-from
//	--sense-mfa-command
//
// These environment variables will be used for any flag that is not set:
//
//	SENSE_EMAIL
//	SENSE_EMAIL_FROM
//	SENSE_PASSWORD
//	SENSE_PASSWORD_FROM
//	SENSE_MFA_FROM
//	SENSE_MFA_COMMAND
//
// These will be blended into one set of credentials that will be used to authenticate the client.
//
// If no flags or environment variables are set, a single unauthenticated Sense client will be created.
package sensecli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/dnesting/sense"
	"gopkg.in/yaml.v3"
)

// PasswordCredentials defines the various ways to specify Sense credentials.
type PasswordCredentials struct {
	Email        string `json:"email,omitempty" yaml:"email,omitempty"`
	EmailFrom    string `json:"email-from,omitempty" yaml:"email-from,omitempty"`
	Password     string `json:"password,omitempty" yaml:"password,omitempty"`
	PasswordFrom string `json:"password-from,omitempty" yaml:"password-from,omitempty"`
	MfaFrom      string `json:"mfa-from,omitempty" yaml:"mfa-from,omitempty"`
	MfaCommand   string `json:"mfa-command,omitempty" yaml:"mfa-command,omitempty"`
}

// Account in the ConfigFile contains credentials for a single Sense account.
type Account struct {
	Credentials *PasswordCredentials `json:"credentials,omitempty" yaml:"credentials,omitempty"`
}

// ConfigFile defines a general Sense client configuration. At the moment this
// is simply a list of Sense accounts.
type ConfigFile struct {
	Accounts []Account `json:"accounts,omitempty" yaml:"accounts,omitempty"`
}

type varNames struct {
	ConfigFile   string
	Email        string
	EmailFrom    string
	Password     string
	PasswordFrom string
	MfaFrom      string
	MfaCommand   string
}

var StandardEnvVars = varNames{
	ConfigFile:   "SENSE_CONFIG",
	Email:        "SENSE_EMAIL",
	EmailFrom:    "SENSE_EMAIL_FROM",
	Password:     "SENSE_PASSWORD",
	PasswordFrom: "SENSE_PASSWORD_FROM",
	MfaFrom:      "SENSE_MFA_FROM",
	MfaCommand:   "SENSE_MFA_COMMAND",
}

var StandardFlagVars = varNames{
	ConfigFile:   "sense-config",
	Email:        "sense-email",
	EmailFrom:    "sense-email-from",
	Password:     "sense-password",
	PasswordFrom: "sense-password-from",
	MfaFrom:      "sense-mfa-from",
	MfaCommand:   "sense-mfa-command",
}

func CredentialsFromStandardEnv() *PasswordCredentials {
	return CredentialsFromEnv(StandardEnvVars)
}

func CredentialsFromEnv(vars varNames) *PasswordCredentials {
	return &PasswordCredentials{
		Email:        os.Getenv(vars.Email),
		EmailFrom:    os.Getenv(vars.EmailFrom),
		Password:     os.Getenv(vars.Password),
		PasswordFrom: os.Getenv(vars.PasswordFrom),
		MfaFrom:      os.Getenv(vars.MfaFrom),
		MfaCommand:   os.Getenv(vars.MfaCommand),
	}
}

func BlendCredentials(a, b PasswordCredentials) *PasswordCredentials {
	if a.Email == "" {
		a.Email = b.Email
	}
	if a.EmailFrom == "" {
		a.EmailFrom = b.EmailFrom
	}
	if a.Password == "" {
		a.Password = b.Password
	}
	if a.PasswordFrom == "" {
		a.PasswordFrom = b.PasswordFrom
	}
	if a.MfaFrom == "" {
		a.MfaFrom = b.MfaFrom
	}
	if a.MfaCommand == "" {
		a.MfaCommand = b.MfaCommand
	}
	return &a
}

func SetupCredentialsFromStandardFlags() *PasswordCredentials {
	return SetupCredentialsFromFlags("", StandardFlagVars, "", "")
}

func SetupCredentialsFromFlags(flagPrefix string, names varNames, flagSuffix, descSuffix string) *PasswordCredentials {
	creds := &PasswordCredentials{}
	flag.StringVar(&creds.Email, flagPrefix+names.Email+flagSuffix, "", "Sense email address"+descSuffix)
	flag.StringVar(&creds.EmailFrom, flagPrefix+names.EmailFrom+flagSuffix, "", "Read Sense email address from this file"+descSuffix)
	flag.StringVar(&creds.Password, flagPrefix+names.Password+flagSuffix, "", "Sense password"+descSuffix)
	flag.StringVar(&creds.PasswordFrom, flagPrefix+names.PasswordFrom+flagSuffix, "", "Read Sense password from this file"+descSuffix)
	flag.StringVar(&creds.MfaFrom, flagPrefix+names.MfaFrom+flagSuffix, "", "Sense MFA code"+descSuffix)
	flag.StringVar(&creds.MfaCommand, flagPrefix+names.MfaCommand+flagSuffix, "", "Read Sense MFA code from the output of this command"+descSuffix)
	return creds
}

func ReadConfigYAML(r io.Reader) (*ConfigFile, error) {
	dec := yaml.NewDecoder(r)
	var cfg ConfigFile
	err := dec.Decode(&cfg)
	return &cfg, err
}

func ReadConfigFileYAML(path string) (*ConfigFile, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ReadConfigYAML(f)
}

func SetupStandardFlags() (*string, *PasswordCredentials) {
	configFile := flag.String(StandardFlagVars.ConfigFile, "", "Path to config file")
	creds := SetupCredentialsFromStandardFlags()
	return configFile, creds
}

func GetConfig(configFile *string, accounts []*PasswordCredentials) (*ConfigFile, error) {
	// use --sense-config or SENSE_CONFIG preferentially
	filename := *configFile
	if filename == "" {
		filename = os.Getenv(StandardEnvVars.ConfigFile)
	}
	if filename != "" {
		return ReadConfigFileYAML(filename)
	}

	// otherwise, synthesize a config
	cfg := &ConfigFile{
		Accounts: make([]Account, len(accounts)),
	}
	for i, cred := range accounts {
		cfg.Accounts[i].Credentials = cred
	}
	return cfg, nil
}

func ConfigFromFileOrFlags(configFile *string, flagCreds *PasswordCredentials) (*ConfigFile, error) {
	envConfig := CredentialsFromStandardEnv()
	blended := BlendCredentials(*flagCreds, *envConfig)
	return GetConfig(configFile, []*PasswordCredentials{blended})
}

func readFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("%s: %w", filename, err)
	}
	return strings.TrimSpace(string(data)), nil
}

func generateSenseCreds(creds *PasswordCredentials) (*sense.PasswordCredentials, error) {
	if creds == nil {
		return nil, nil // empty credentials = unauthenticated
	}

	var userpass sense.PasswordCredentials
	var err error

	if creds.Email != "" {
		userpass.Email = creds.Email
	} else if creds.EmailFrom != "" {
		userpass.Email, err = readFile(creds.EmailFrom)
		if err != nil {
			return nil, err
		}
	}

	if creds.Password != "" {
		userpass.Password = creds.Password
	} else if creds.PasswordFrom != "" {
		userpass.Password, err = readFile(creds.PasswordFrom)
		if err != nil {
			return nil, err
		}
	}

	if creds.MfaFrom != "" {
		code, err := readFile(creds.MfaFrom)
		if err != nil {
			return nil, err
		}
		userpass.MfaFn = func(_ context.Context) (string, error) {
			return code, nil
		}
	} else if creds.MfaCommand != "" {
		userpass.MfaFn = func(_ context.Context) (string, error) {
			cmd := exec.Command(creds.MfaCommand)
			out, err := cmd.Output()
			if err != nil {
				return "", fmt.Errorf("exec %q: %w", creds.MfaCommand, err)
			}
			return strings.TrimSpace(string(out)), nil
		}
	}
	return &userpass, nil
}

func clientFromAccount(ctx context.Context, acct Account, opts ...sense.Option) (*sense.Client, error) {
	senseCreds, err := generateSenseCreds(acct.Credentials)
	if err != nil {
		return nil, err
	}
	if senseCreds == nil || senseCreds.Email == "" {
		return sense.New(opts...), nil
	}
	return sense.Connect(ctx, senseCreds, opts...)
}

func clientsFromConfig(ctx context.Context, config ConfigFile, opts ...sense.Option) (clients []*sense.Client, err error) {
	for _, acct := range config.Accounts {
		client, err := clientFromAccount(ctx, acct, opts...)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func CreateClients(ctx context.Context, configFile *string, flagCreds *PasswordCredentials, opts ...sense.Option) ([]*sense.Client, error) {
	cfg, err := ConfigFromFileOrFlags(configFile, flagCreds)
	if err != nil {
		return nil, err
	}
	return clientsFromConfig(ctx, *cfg, opts...)
}
