package application

// TemplateNamesForTest exposes the package-internal templateNames
// helper to external _test packages. The `_test.go` suffix means
// the symbol only exists in the test binary; production callers
// cannot reach it.
func TemplateNamesForTest() ([]string, error) {
	return templateNames()
}

// RenderTemplateForTest exposes the package-internal renderTemplate
// helper to external _test packages so the error path
// (template-not-found) is reachable.
func RenderTemplateForTest(name, projectName string) ([]byte, error) {
	return renderTemplate(name, templateData{Name: projectName})
}
