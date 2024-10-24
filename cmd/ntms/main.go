// A command line interface for the Agent to use as a Docker image entrypoint.
package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// A cli maintains the state of the CLI.
// It contains all the logic and data that is common to all commands.
// It needs to hold the configuration so it can create the components, and
// it needs to hold a reference to the agent so it can shutdown gracefully.
type cli struct {
	cfg Config
}

// A Config contains the configuration needed to setup the components for services.
type Config struct {
	// DataDir stores the log and raft data.
	DataDir string
	// BindAddr is the address serf runs on.
	BindAddr string
	// NodeName is the name of the node.
	NodeName string
	// StartJoinAddrs is the list of addresses of the other nodes in the cluster.
	StartJoinAddrs []string
	// Bootstrap should be set to true when starting the first node of the cluster.
	Bootstrap bool
}

func main() {
	cli := &cli{}

	cmd := &cobra.Command{
		Use:     "ntms",
		PreRunE: cli.setupConfig,
		RunE:    cli.run,
		PostRunE: func(cmd *cobra.Command, args []string) error {
			shutdown := make(chan os.Signal, 1)
			signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
			// signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
			<-shutdown
			return nil
		},
	}
	log.Println("before setup flags")
	if err := setupFlags(cmd); err != nil {
		log.Println("failed to initiate flags", err)
		log.Fatal(err)
	}
	log.Println("after setup flags")

	log.Println("before cmd execute")
	if err := cmd.Execute(); err != nil {
		log.Println("failed to execute command", err)
		log.Fatal(err)
	}
	log.Println("after cmd execute")

}

// setupFlags registers flags for all commands, necessary too run the agent.
func setupFlags(cmd *cobra.Command) error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Flags().String("config-file", "", "Path to config file.")

	dataDir := path.Join(os.TempDir(), "ntms")
	cmd.Flags().String("data-dir",
		dataDir,
		"Directory to store log and Raft data.")
	cmd.Flags().String("node-name", hostname, "Unique server ID.")

	cmd.Flags().String("bind-addr",
		"127.0.0.1:8401",
		"Address to bind Serf on.")
	cmd.Flags().Bool("bootstrap", false, "Bootstrap the cluster.")

	return viper.BindPFlags(cmd.Flags())
}

// setupConfig creates the agent configuration from the CLI flags.
func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {
	var err error

	log.Println("start PreRunE ")

	// viper.SetConfigFile(configFile)
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return err
	}
	viper.SetConfigFile(configFile)

	if err = viper.ReadInConfig(); err != nil {
		// it's ok if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	c.cfg.DataDir = viper.GetString("data-dir")
	c.cfg.NodeName = viper.GetString("node-name")
	c.cfg.BindAddr = viper.GetString("bind-addr")
	c.cfg.StartJoinAddrs = viper.GetStringSlice("start-join-addrs")
	c.cfg.Bootstrap = viper.GetBool("bootstrap")

	log.Println("end PreRunE ")

	return nil
}

// run starts the agent, and terminates the agent when it detects either a SIGINT or SIGTERM.
// SIGINT occurs when the user hits Ctrl+C.
// SIGTERM occurs when the Docker container manager stops the container.
func (c *cli) run(cmd *cobra.Command, args []string) error {

	log.Println("start RunE ")

	log.Println("c.cfg.DataDir", c.cfg.DataDir)
	log.Println("c.cfg.NodeName", c.cfg.NodeName)
	log.Println("c.cfg.BindAddr", c.cfg.BindAddr)
	for _, addr := range c.cfg.StartJoinAddrs {
		log.Println("Join Address", addr)
	}
	log.Println("c.cfg.Bootstrap", c.cfg.Bootstrap)

	log.Println("end RunE ")

	return nil
}
