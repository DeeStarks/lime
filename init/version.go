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