package sponsorblock

import (
	"runtime/debug"
	"strings"
)

//goland:noinspection GoUnusedConst
const (
	GitHub = "https://github.com/DisgoOrg/sponsorblock-plugin"
	Name   = "sponsorblock-plugin"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	Version = getVersion()
)

func getVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, dep := range bi.Deps {
			if strings.Contains(GitHub, dep.Path) {
				return dep.Version
			}
		}
	}
	return "unknown"
}
