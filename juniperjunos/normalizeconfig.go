package juniperjunos

import "strings"

// NormalizeConfig "normalizes" the configuration provided -- generally this means it replaces any
// kind of "header" pattern(s), but may vary from platform to platform.
func (p *Platform) NormalizeConfig(config string) string {
	config = p.patterns.outputHeader.ReplaceAllString(config, "")
	config = p.patterns.edit.ReplaceAllString(config, "")
	config = strings.ReplaceAll(config, "\n\n", "\n")

	return config
}
