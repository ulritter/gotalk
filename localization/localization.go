package localization

import (
	"gotalk/constants"
	"gotalk/models"
	"gotalk/utils"
	"log"

	language "github.com/moemoe89/go-localization"
)

// Initialize localization environment, if localization file is not present, create one by downloading it from github
func InitLocalization(a *models.Application) error {
	if !utils.FileExists(constants.LANGFILE) {
		fileUrl := constants.RAWFILE
		err := utils.GetFileFromGithub(constants.LANGFILE, fileUrl)
		if err != nil {
			log.Fatal(err)
		}
	}
	var err error
	cfg := language.New()
	cfg.BindPath(constants.LANGFILE)
	cfg.BindMainLocale("en")
	a.Lang, err = cfg.Init()
	if err != nil {
		log.Fatal(err)
	}
	return err
}
