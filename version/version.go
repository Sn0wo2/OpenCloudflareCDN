package version

import (
	"fmt"
	"time"
)

var (
	version = "v0.0.0"
	commit  = "dev"
	date    = "1970-01-01T00:00:00Z"
)

func GetVersion() string {
	return version
}

func GetCommit() string {
	return commit
}

func GetShortCommit() string {
	if len(commit) > 7 {
		return commit[:7]
	}

	return commit
}

func GetDate() string {
	return date
}

func GetDateTime() time.Time {
	t, _ := time.Parse(time.RFC3339, date)

	return t
}

func GetFormatVersion() string {
	return fmt.Sprintf("%s-%s(%s)", GetVersion(), GetShortCommit(), GetDateTime().Format("060102150405"))
}

func SetVersion(ver string) {
	version = ver
}

func SetCommit(comm string) {
	commit = comm
}

func SetDate(dat string) {
	date = dat
}
