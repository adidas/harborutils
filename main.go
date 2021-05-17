package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

//TODO: Use viper to allow reading from config files

var harborServer, harborServerTarget string
var harborUser, harborUserTarget string
var harborPassword, harborPasswordTarget string
var harborAPIVersion, harborAPIVersionTarget string

// database
var dbHostSource, dbHostTarget string
var dbUserSource, dbUserTarget string
var dbPasswordSource, dbPasswordTarget string
var dbPortSource, dbPortTarget int
var verbose bool

var rootCmd = &cobra.Command{
	Use:   "harborutils",
	Short: "Interacts with harbor registry API",
}

var getProjectsGroupsCmd = &cobra.Command{
	Use:   "getProjects",
	Short: "Get projects from Harbor",
	Long:  `Get projects from Harbor.`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		harborPassword = promptForPassword(harborServer, harborPassword)
		if harborAPIVersion != "" {
			harborAPIVersion = harborAPIVersion + "/"
		}
		for _, project := range listProjects(harborServer, harborUser, harborPassword, harborAPIVersion) {
			fmt.Println(project.Name)
		}
	},
}

var getGroupsCmd = &cobra.Command{
	Use:   "getGroups <prefix>",
	Short: "Get groups from Harbor",
	Long: `Get groups from Harbor.
	
If provided, only look for groups starting by <prefix>	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		harborPassword = promptForPassword(harborServer, harborPassword)
		prefix := ""
		if len(args) > 0 {
			prefix = args[0]
		}
		if harborAPIVersion != "" {
			harborAPIVersion = harborAPIVersion + "/"
		}
		groups, err := getGroupsFromPrefix(harborServer, harborUser, harborPassword, harborAPIVersion, prefix)
		if err == nil {
			for _, group := range groups {
				fmt.Println(group.GroupName)
				fmt.Println(group.LdapGroupDN)
				fmt.Println(group.GroupType)
			}
		} else {
			log.Fatal(err)
		}
	},
}

var deleteGroupsCmd = &cobra.Command{
	Use:   "deleteGroups <prefix>",
	Short: "Delete groups from Harbor",
	Long: `Delete groups from Harbor.
	
If provided, only look for groups starting by <prefix>	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		harborPassword = promptForPassword(harborServer, harborPassword)
		prefix := ""
		if len(args) > 0 {
			prefix = args[0]
		}
		if harborAPIVersion != "" {
			harborAPIVersion = harborAPIVersion + "/"
		}
		groups, err := getGroupsFromPrefix(harborServer, harborUser, harborPassword, harborAPIVersion, prefix)
		if err == nil {
			deleteGroups(harborServer, harborUser, harborPassword, harborAPIVersion, groups)
		}
	},
}

var syncGrantsCmd = &cobra.Command{
	Use:   "syncGrants <projectList>",
	Short: "Propagate grants from primary harbor to secondary",
	Long: `Propagate grants form primary server to secondary

Users and groups in registry will be checked against secondary registry and projects updated
if found.
- projectList should contain the project names (separated by comma)`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		harborPassword = promptForPassword(harborServer, harborPassword)
		harborPasswordTarget = promptForPassword(harborServerTarget, harborPasswordTarget)
		var projects ProjectListResponse
		if harborAPIVersion != "" {
			harborAPIVersion = harborAPIVersion + "/"
		}
		if harborAPIVersionTarget != "" {
			harborAPIVersionTarget = harborAPIVersionTarget + "/"
		}
		if len(args) > 0 {
			project_list := strings.Split(args[0], ",")
			for _, prj := range project_list {
				project, err := getProject(harborServer, harborUser, harborPassword, prj, harborAPIVersion)
				if err == nil {
					projects = append(projects, project)
				} else {
					log.Println(err)
				}
			}

		} else {
			projects = listProjects(harborServer, harborUser, harborPassword, harborAPIVersion)
		}
		for _, project := range projects {
			fmt.Println("Sync project ", project.Name)
			syncProjectGrants(project)
		}

	},
}

var syncLabelsCmd = &cobra.Command{
	Use:   "syncLabels <projectList>",
	Short: "Propagate project labels from primary harbor to secondary",
	Long: `Propagate project labels form primary server to secondary

- projectList should contain the project names (separated by comma)`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		harborPassword = promptForPassword(harborServer, harborPassword)
		harborPasswordTarget = promptForPassword(harborServerTarget, harborPasswordTarget)
		var projects ProjectListResponse
		if harborAPIVersion != "" {
			harborAPIVersion = harborAPIVersion + "/"
		}
		if harborAPIVersionTarget != "" {
			harborAPIVersionTarget = harborAPIVersionTarget + "/"
		}
		if len(args) > 0 {
			project_list := strings.Split(args[0], ",")
			for _, prj := range project_list {
				project, err := getProject(harborServer, harborUser, harborPassword, prj, harborAPIVersion)
				if err == nil {
					projects = append(projects, project)
				} else {
					log.Println(err)
				}
			}

		} else {
			projects = listProjects(harborServer, harborUser, harborPassword, harborAPIVersion)
		}
		for _, project := range projects {
			fmt.Println("Sync project ", project.Name)
			syncProjectLabels(project)
		}

	},
}

var importLdapUsersCmd = &cobra.Command{
	Use:   "importLdapUsers <harbor>",
	Short: "Propagate users from primary harbor to secondary",
	Long:  `Propagate users from primary harbor to secondary`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		harborPassword = promptForPassword(harborServer, harborPassword)
		harborPasswordTarget = promptForPassword(harborServerTarget, harborPasswordTarget)

		if harborAPIVersion != "" {
			harborAPIVersion = harborAPIVersion + "/"
		}
		if harborAPIVersionTarget != "" {
			harborAPIVersionTarget = harborAPIVersionTarget + "/"
		}
		users, err := getSourceUsers(harborServer, harborUser, harborPassword, harborAPIVersion)
		if err == nil {
			usernames := make([]string, len(users))
			for i, user := range users {
				usernames[i] = user.Username
			}

			newUserImport := UserImport{}
			newUserImport.LdapUIDList = usernames

			importLdapUser(harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget, newUserImport)
		} else {
			log.Println(err)
		}
	},
}

var importLdapGroupsCmd = &cobra.Command{
	Use:   "importLdapGroups <harbor>",
	Short: "Propagate groups from primary harbor to secondary",
	Long:  `Propagate groups from primary harbor to secondary`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		harborPassword = promptForPassword(harborServer, harborPassword)
		harborPasswordTarget = promptForPassword(harborServerTarget, harborPasswordTarget)

		if harborAPIVersion != "" {
			harborAPIVersion = harborAPIVersion + "/"
		}
		if harborAPIVersionTarget != "" {
			harborAPIVersionTarget = harborAPIVersionTarget + "/"
		}
		groups, err := listGroups(harborServer, harborUser, harborPassword, harborAPIVersion)
		if err == nil {
			for _, g := range groups {
				fmt.Println(g.GroupName)
			}
			importGroups(harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget, groups)
		}

	},
}

// Short: "Propagate project labels from primary harbor to secondary",
// Long: `Propagate project labels form primary server to secondary

var syncRobotAccountCmd = &cobra.Command{
	Use:   "syncRobotAccount",
	Short: "Propagate robot account from primary harbor to secundary",
	Long:  "Propagate robot account from primary harbor to secundary",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// robots := getRobots(harborServer, harborUser, harborPassword, harborAPIVersion)
		// for _, robot := range robots {
		// 	fmt.Printf("%+v", robot)
		// }
		syncRobots(harborServer, harborUser, harborPassword, harborAPIVersion, harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget,
			clientDb(dbHostSource, dbUserSource, dbPasswordSource, dbPortSource, verbose),
			clientDb(dbHostTarget, dbUserTarget, dbPasswordTarget, dbPortTarget, verbose))

	},
}
var fixEmptyEmailsDbCmd = &cobra.Command{
	Use:   "fixEmptyEmails",
	Short: "Fix empty emails in database",
	Long:  `Fix empty emails in database`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		FixEmptyEmails(clientDb(dbHostSource, dbUserSource, dbPasswordSource, dbPortSource, verbose))
	},
}

var syncUsersDbCmd = &cobra.Command{
	Use:   "syncUsersDb",
	Short: "Sync useres between harbor primarty and harbor secundary",
	Long:  "Sync useres between harbor primarty and harbor secundary",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SyncUsersDatabase(clientDb(dbHostSource, dbUserSource, dbPasswordSource, dbPortSource, verbose),
			clientDb(dbHostTarget, dbUserTarget, dbPasswordTarget, dbPortTarget, verbose))
	},
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&harborServer, "harbor", "s", "", "Harbor Server address")
	rootCmd.MarkPersistentFlagRequired("harbor")
	rootCmd.PersistentFlags().StringVarP(&harborUser, "user", "u", "", "Username Harbor")
	rootCmd.MarkPersistentFlagRequired("user")
	rootCmd.PersistentFlags().StringVarP(&harborPassword, "password", "p", "", "Password")
	rootCmd.PersistentFlags().StringVarP(&harborAPIVersion, "apiVersion", "v", "", "APIVersion (ie v2.0)")

	syncGrantsCmd.PersistentFlags().StringVarP(&harborServerTarget, "harbor2", "", "", "Harbor Secondary Server address")
	syncGrantsCmd.MarkPersistentFlagRequired("harbor2")
	syncGrantsCmd.PersistentFlags().StringVarP(&harborUserTarget, "user2", "", "", "Username Secondary Harbor")
	syncGrantsCmd.MarkPersistentFlagRequired("user2")
	syncGrantsCmd.PersistentFlags().StringVarP(&harborPasswordTarget, "password2", "", "", "Password Secondary Harbor")
	syncGrantsCmd.PersistentFlags().StringVarP(&harborAPIVersionTarget, "apiVersion2", "", "", "API Version Secondary Harbor (ie v2.0) ")

	syncLabelsCmd.PersistentFlags().StringVarP(&harborServerTarget, "harbor2", "", "", "Harbor Secondary Server address")
	syncLabelsCmd.MarkPersistentFlagRequired("harbor2")
	syncLabelsCmd.PersistentFlags().StringVarP(&harborUserTarget, "user2", "", "", "Username Secondary Harbor")
	syncLabelsCmd.MarkPersistentFlagRequired("user2")
	syncLabelsCmd.PersistentFlags().StringVarP(&harborPasswordTarget, "password2", "", "", "Password Secondary Harbor")
	syncLabelsCmd.PersistentFlags().StringVarP(&harborAPIVersionTarget, "apiVersion2", "", "", "API Version Secondary Harbor (ie v2.0)")

	importLdapUsersCmd.PersistentFlags().StringVarP(&harborServerTarget, "harbor2", "", "", "Harbor Secondary Server address")
	importLdapUsersCmd.MarkPersistentFlagRequired("harbor2")
	importLdapUsersCmd.PersistentFlags().StringVarP(&harborUserTarget, "user2", "", "", "Username Secondary Harbor")
	importLdapUsersCmd.MarkPersistentFlagRequired("user2")
	importLdapUsersCmd.PersistentFlags().StringVarP(&harborPasswordTarget, "password2", "", "", "Password Secondary Harbor")
	importLdapUsersCmd.PersistentFlags().StringVarP(&harborAPIVersionTarget, "apiVersion2", "", "", "API Version Secondary Harbor (ie v2.0)")

	importLdapGroupsCmd.PersistentFlags().StringVarP(&harborServerTarget, "harbor2", "", "", "Harbor Secondary Server address")
	importLdapGroupsCmd.MarkPersistentFlagRequired("harbor2")
	importLdapGroupsCmd.PersistentFlags().StringVarP(&harborUserTarget, "user2", "", "", "Username Secondary Harbor")
	importLdapGroupsCmd.MarkPersistentFlagRequired("user2")
	importLdapGroupsCmd.PersistentFlags().StringVarP(&harborPasswordTarget, "password2", "", "", "Password Secondary Harbor")
	importLdapGroupsCmd.PersistentFlags().StringVarP(&harborAPIVersionTarget, "apiVersion2", "", "", "API Version Secondary Harbor (ie v2.0)")

	fixEmptyEmailsDbCmd.PersistentFlags().StringVarP(&dbHostSource, "dbHostSource", "", "", "source database host")
	fixEmptyEmailsDbCmd.MarkPersistentFlagRequired("dbHostSource")
	fixEmptyEmailsDbCmd.PersistentFlags().StringVarP(&dbUserSource, "dbUserSource", "", "", "source database user")
	fixEmptyEmailsDbCmd.MarkPersistentFlagRequired("dbUserSource")
	fixEmptyEmailsDbCmd.PersistentFlags().StringVarP(&dbPasswordSource, "dbPasswordSource", "", "", "source database password")
	fixEmptyEmailsDbCmd.MarkPersistentFlagRequired("dbPasswordSource")
	fixEmptyEmailsDbCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "verbose output, shows all sql operations")
	fixEmptyEmailsDbCmd.MarkPersistentFlagRequired("")

	syncUsersDbCmd.PersistentFlags().StringVarP(&dbHostSource, "dbHostSource", "", "", "source database host")
	syncUsersDbCmd.MarkPersistentFlagRequired("dbHostSource")
	syncUsersDbCmd.PersistentFlags().StringVarP(&dbUserSource, "dbUserSource", "", "", "source database user")
	syncUsersDbCmd.MarkPersistentFlagRequired("dbUserSource")
	syncUsersDbCmd.PersistentFlags().StringVarP(&dbPasswordSource, "dbPasswordSource", "", "", "source database password")
	syncUsersDbCmd.MarkPersistentFlagRequired("dbPasswordSource")
	syncUsersDbCmd.PersistentFlags().StringVarP(&dbHostTarget, "dbHostTarget", "", "", "source database host")
	syncUsersDbCmd.MarkPersistentFlagRequired("dbHostTarget")
	syncUsersDbCmd.PersistentFlags().StringVarP(&dbUserTarget, "dbUserTarget", "", "", "source database user")
	syncUsersDbCmd.MarkPersistentFlagRequired("dbUserTarget")
	syncUsersDbCmd.PersistentFlags().StringVarP(&dbPasswordTarget, "dbPasswordTarget", "", "", "source database password")
	syncUsersDbCmd.MarkPersistentFlagRequired("dbPasswordTarget")
	syncUsersDbCmd.PersistentFlags().IntVarP(&dbPortSource, "dbPortSource", "", 5432, "source database port, defualt 5432")
	syncUsersDbCmd.PersistentFlags().IntVarP(&dbPortTarget, "dbPortTarget", "", 5432, "target database port, default 5432")
	syncUsersDbCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "verbose output, shows all sql operations")
	syncUsersDbCmd.MarkPersistentFlagRequired("")

	syncRobotAccountCmd.PersistentFlags().StringVarP(&harborServerTarget, "harbor2", "", "", "Harbor Secondary Server address")
	syncRobotAccountCmd.MarkPersistentFlagRequired("harbor2")
	syncRobotAccountCmd.PersistentFlags().StringVarP(&harborUserTarget, "user2", "", "", "Username Secondary Harbor")
	syncRobotAccountCmd.MarkPersistentFlagRequired("user2")
	syncRobotAccountCmd.PersistentFlags().StringVarP(&harborPasswordTarget, "password2", "", "", "Password Secondary Harbor")
	syncRobotAccountCmd.PersistentFlags().StringVarP(&harborAPIVersionTarget, "apiVersion2", "", "", "API Version Secondary Harbor (ie v2.0)")
	syncRobotAccountCmd.PersistentFlags().StringVarP(&dbHostSource, "dbHostSource", "", "", "source database host")
	syncRobotAccountCmd.MarkPersistentFlagRequired("dbHostSource")
	syncRobotAccountCmd.PersistentFlags().StringVarP(&dbUserSource, "dbUserSource", "", "", "source database user")
	syncRobotAccountCmd.MarkPersistentFlagRequired("dbUserSource")
	syncRobotAccountCmd.PersistentFlags().StringVarP(&dbPasswordSource, "dbPasswordSource", "", "", "source database password")
	syncRobotAccountCmd.MarkPersistentFlagRequired("dbPasswordSource")
	syncRobotAccountCmd.PersistentFlags().StringVarP(&dbHostTarget, "dbHostTarget", "", "", "source database host")
	syncRobotAccountCmd.MarkPersistentFlagRequired("dbHostTarget")
	syncRobotAccountCmd.PersistentFlags().StringVarP(&dbUserTarget, "dbUserTarget", "", "", "source database user")
	syncRobotAccountCmd.MarkPersistentFlagRequired("dbUserTarget")
	syncRobotAccountCmd.PersistentFlags().StringVarP(&dbPasswordTarget, "dbPasswordTarget", "", "", "source database password")
	syncRobotAccountCmd.MarkPersistentFlagRequired("dbPasswordTarget")
	syncRobotAccountCmd.PersistentFlags().IntVarP(&dbPortSource, "dbPortSource", "", 5432, "source database port, defualt 5432")
	syncRobotAccountCmd.PersistentFlags().IntVarP(&dbPortTarget, "dbPortTarget", "", 5432, "target database port, default 5432")
	syncRobotAccountCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "", false, "verbose output, shows all sql operations")
	syncRobotAccountCmd.MarkPersistentFlagRequired("")

	rootCmd.AddCommand(getProjectsGroupsCmd)
	rootCmd.AddCommand(getGroupsCmd)
	rootCmd.AddCommand(deleteGroupsCmd)
	rootCmd.AddCommand(syncGrantsCmd)
	rootCmd.AddCommand(syncLabelsCmd)
	rootCmd.AddCommand(importLdapUsersCmd)
	rootCmd.AddCommand(importLdapGroupsCmd)
	/*rootCmd.AddCommand(rmUsersFromGroupsCmd)
	rootCmd.AddCommand(listUsersInGroupsCmd)
	rootCmd.AddCommand(listManagedCmd)
	rootCmd.AddCommand(listGroupsCmd)
	rootCmd.AddCommand(whoManagesCmd)*/
	rootCmd.AddCommand(syncRobotAccountCmd)

	rootCmd.AddCommand(fixEmptyEmailsDbCmd)
	rootCmd.AddCommand(syncUsersDbCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func promptForPassword(field, pass string) string {
	var password string
	if pass == "" {
		print("Password for ", field, ": ")
		pass, _ := gopass.GetPasswd()
		password = string(pass)
	} else {
		password = pass
	}
	return password

}
