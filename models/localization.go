package models

import (
	"gotalk/utils"

	language "github.com/moemoe89/go-localization"
)

// Initialize localization environment, if localization file is not present, create one by downloading it from github
func (a *Application) InitLocalization() error {
	if !utils.FileExists(LANGFILE) {
		fileUrl := RAWFILE
		err := utils.GetFileFromGithub(LANGFILE, fileUrl)
		if err != nil {
			panic(err)
		}
	}
	var err error
	cfg := language.New()
	cfg.BindPath(LANGFILE)
	cfg.BindMainLocale("en")
	a.Lang, err = cfg.Init()
	if err != nil {
		panic(err)
	}
	return err
}
