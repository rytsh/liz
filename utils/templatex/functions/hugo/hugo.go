// This package functions, scraped from Hugo's functions.
package hugo

import (
	"os"
	"path/filepath"
)

func FuncMapFn(workdir string) func() map[string]interface{} {
	if !filepath.IsAbs(workdir) {
		wd, err := os.Getwd()
		if err == nil {
			workdir = filepath.Join(wd, workdir)
		}
	}

	return func() map[string]interface{} {
		ns := New(workdir)

		fMap := map[string]interface{}{
			"readFile":   ns.ReadFile,
			"readDir":    ns.ReadDir,
			"stat":       ns.Stat,
			"fileExists": ns.FileExists,
		}

		return fMap
	}
}
