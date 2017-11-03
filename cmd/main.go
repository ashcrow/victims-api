/*
* victi.ms API microservice
* Copyright (C) 2017 The victi.ms team
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU Affero General Public License as published
* by the Free Software Foundation, either version 3 of the License, or
* (at your option) any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU Affero General Public License for more details.
*
* You should have received a copy of the GNU Affero General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/victims/victims-api/api"
	"github.com/victims/victims-common/db"
	"github.com/victims/victims-common/log"
	//	"github.com/victims/victims-api/middleware"
)

var version string
var commitHash string
var buildTime string

// versionCmd returns version information for the current binary
func versionCmd() *cobra.Command {

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Version of the binary",
		Long:  `Returns the version and build information for this binary`,

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Version: %s\n", version)
			cmd.Printf("Go Version: %s\n", runtime.Version())
			cmd.Printf("Build Time: %s\n", buildTime)
			cmd.Printf("Commit: %s\n", commitHash)
		},
	}
	return versionCmd
}

// runCmd creates the proper cobra.Command for exeucting the server
func runCmd() *cobra.Command {
	var bind string
	var certFile string
	var certKey string
	var mongoDBHost string
	var mongoDBPassword string
	var mongoDBUser string
	var mongoDBDatabase string
	var mongoDBUseEnv bool

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Run the server",
		Long:  `Run the victims-public-api server`,

		Run: func(cmd *cobra.Command, args []string) {
			router := gin.New()
			//router.Use(api.AuthenticationMiddleware())
			//router.Use(middleware.Limit())
			mount(router)
			// If the user requested us to pull from the environment then do so
			if mongoDBUseEnv == true {
				mongoDBUser, mongoDBPassword, mongoDBDatabase = db.GetMongoDBAuthFromEnv()
				log.Logger.Debugf("Getting mongodb variables from the environment")
			}
			log.Logger.Warnf("%s %s", mongoDBHost, mongoDBDatabase)
			_, err := db.New(mongoDBHost, mongoDBDatabase)
			if err != nil {
				log.Logger.Fatalf("Unable to connect to the databse: %s", err)
			}
			if certFile != "" && certKey != "" {
				log.Logger.Info("Serving with TLS")
				log.Logger.Debugf("certFile=%s, certKey=%s", certFile, certKey)
				router.RunTLS(bind, certFile, certKey)
			} else {
				log.Logger.Info(
					"Serving without TLS as both certKey and certFile were not provided")
				router.Run(bind)
			}
		},
	}

	// Add flags
	runCmd.PersistentFlags().StringVar(
		&bind, "bind", "0.0.0.0:8080", "Bind address")
	runCmd.PersistentFlags().StringVar(
		&certFile, "cert-file", "", "Path to certificate file")
	runCmd.PersistentFlags().StringVar(
		&certKey, "cert-key", "", "Path to certificate key")
	runCmd.PersistentFlags().StringVar(
		&mongoDBHost, "mongodb-host", "127.0.0.1", "MongoDB host to connect to")
	runCmd.PersistentFlags().StringVar(
		&mongoDBUser, "mongodb-user", "", "MongoDB user to auth with")
	runCmd.PersistentFlags().StringVar(
		&mongoDBPassword, "mongodb-password", "", "MongoDB password to auth with")
	runCmd.PersistentFlags().StringVar(
		&mongoDBDatabase, "mongodb-database", "", "MongoDB database to auth to")
	runCmd.PersistentFlags().BoolVar(
		&mongoDBUseEnv, "mongodb-use-env",
		false, "If set MONGODB_USER/PASSWORD/DATABASE environment variables will be used.")
	return runCmd
}

// main is the main entry point for running the victims-public-api service.
func main() {
	rootCmd := &cobra.Command{}
	rootCmd.AddCommand(runCmd())
	rootCmd.AddCommand(versionCmd())
	rootCmd.Execute()
}

// mount takes care of mounting all endpoints on to the router
func mount(router *gin.Engine) {
	log.Logger.Debug("Mounting all endpoints")
	api.HashMounts(router)
	log.Logger.Debug("Finished mounting endpoints")
}
