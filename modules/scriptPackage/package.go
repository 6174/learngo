package scriptPackage


type ScriptPackage struct {
	Name        string `yaml: "name"`
	Description string `yaml: "description"`
	Version     string `yaml: "version"`
	Keywords  []string `yaml: "keywords"`
	Catalogs  []string `yaml: "catalog"`
	Repositry   string `yaml: "repositry"`
	License     string `yaml: "license"`
	Homepage    string `yaml: "homepage"`
	Issues        string `yaml: "issues"`
	Author      string `yaml: "author"`
	Email       string `yaml: "email"`
}