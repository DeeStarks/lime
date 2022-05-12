package init

import "github.com/DeeStarks/lime/configs"

type Version struct {
	configs.LimeVersion
}

func NewVersion(version configs.LimeVersion) *Version {
	return &Version{
		version,
	}
}

func (v *Version) GetNumber() int {
	return v.Number
}

func (v *Version) GetLogo() []string {
	return v.Logo
}

func (v *Version) GetInfo() string {
	return v.InfoText
}

func (v *Version) GetAuthor() string {
	return v.Author
}
