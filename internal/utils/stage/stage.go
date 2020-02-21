package stage

import "os"

const (
	StageEnv     = "STAGE"
	Production   = "production"
	Staging      = "staging"
	Development  = "development"
	DefaultStage = Development
)

func IsProd() bool {
	return Get() == Production
}

func IsStaging() bool {
	return Get() == Staging
}

func Get() string {
	st := os.Getenv(StageEnv)
	if len(st) == 0 {
		return DefaultStage
	}
	return st
}
