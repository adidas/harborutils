package main

import (
	"main/server"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run a server exposing some options of the cli",
	Long:  `Run a server exposing some options of the cli, more info in http://localhost:8080/swagger/index.html`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Execute(server.ServerConfig{
			ClientId: clientId,
			TenantId: tenant,
			Host:     harborServer,
		})
	},
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&clientId, "oidcClient", "", "", "Oidc client id for authentication")
	serverCmd.MarkPersistentFlagRequired("oidcClient")
	serverCmd.PersistentFlags().StringVarP(&tenant, "tenant", "", "", "Azure tenant for oidc authentication")
	serverCmd.MarkPersistentFlagRequired("tenant")
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
