//go:generate go run golang.org/x/tools/cmd/stringer -type=SkillLevelId
package jamulusprotocol

type SkillLevelId uint8

const (
	SkillNone         SkillLevelId = 0
	SkillBeginner     SkillLevelId = 1
	SkillIntermediate SkillLevelId = 2
	SkillExpert       SkillLevelId = 3
)
