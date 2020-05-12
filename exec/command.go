package exec

func CreateSchemaCommand(database, username, password string) []string {
	return []string{"mysql", "--user=" + username, "--password=" + password, "-e", "CREATE SCHEMA " + database + " DEFAULT CHARACTER SET utf8mb4"}
}
