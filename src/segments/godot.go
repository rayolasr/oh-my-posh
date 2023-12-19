package segments

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jandedobbeleer/oh-my-posh/src/platform"
	"github.com/jandedobbeleer/oh-my-posh/src/properties"
)

type Godot struct {
	language
	version string
}

func (g *Godot) Template() string {
	// Personaliza este template según tus necesidades
	return fmt.Sprintf(`{{if .Godot}}Godot %s{{end}}`, g.version)
}

func (g *Godot) Init(props properties.Properties, env platform.Environment) {
	godotFiles := []string{
		"project.godot",
		"*.tscn",
		"*.gd",
		"*.import",
		"*.tres",
	}

	g.language = language{
		env:        env,
		props:      props,
		extensions: godotFiles,
	}
}

func (g *Godot) Enabled() bool {
	// Check if any Godot-specific files are present in the current directory
	for _, pattern := range g.language.extensions {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			fmt.Println("Error:", err)
			return false
		}
		if len(matches) > 0 {
			// If Godot files are found, attempt to read version information
			g.version = g.readGodotVersion("project.godot")
			return g.version != ""
		}
	}

	return false
}

func (g *Godot) readGodotVersion(filePath string) string {
	// Attempt to read the content of the "project.godot" file
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading Godot project file:", err)
		return ""
	}

	// Use a regular expression to extract the Godot version
	versionRegex := regexp.MustCompile(`^ *\[gd_scene load_steps=.* format=2.*\] *$.*^ *version.*=(?P<version>.*)$`)
	matches := versionRegex.FindStringSubmatch(string(content))
	if len(matches) >= 2 {
		return matches[1]
	}

	return ""
}
