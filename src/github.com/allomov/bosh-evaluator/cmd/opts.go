package cmd

import (
	boshcmd "github.com/cloudfoundry/bosh-cli/cmd"
)

type BoshOpts struct {
	// -----> Global options
	VersionOpt    func() error `long:"version" short:"v" description:"Show CLI version"`
	DeploymentOpt string       `long:"name" short:"n" description:"Environment or Foundation name" env:"FOUNDATION_NAME"`
	VaultAddrOpt  string       `long:"vault-addr" description:"Vault Address" env:"VAULT_ADDR"`
	VaultTokenOpt string       `long:"vault-token" description:"Vault Address" env:"VAULT_TOKEN"`
	Help          HelpOpts     `command:"help" description:"Show this help message"`
	Read          ReadOpts     `command:"fetch"  alias:"f"  description:"Fetch variables"`
	Write         WriteOpts    `command:"fetch"  alias:"f"  description:"Fetch variables"`
}

type cmd struct{}

type HelpOpts struct {
	cmd
}

type ReadOpts struct {
	Args ReadArgs `positional-args:"true" required:"true"`
	boshcmd.VarFlags
	// ConfigPath string `long:"config" value-name:"CONFIG" description:"Config with values description"`
	cmd
}

type WriteOpts struct {
	Args ReadArgs `positional-args:"true" required:"true"`
	cmd
}

type ReadArgs struct {
	ConfigPath boshcmd.FileBytesWithPathArg `positional-arg-name:"CONFIG_PATH" description:"Path to a config file"`
}
