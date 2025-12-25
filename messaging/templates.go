package messaging

import (
	"fmt"
	"strings"

	"github.com/automation-poc/browser-automation/storage"
)

type Template struct {
	Name    string
	Content string
	Variables []string
}

var DefaultTemplates = []Template{
	{
		Name:    "follow_up_general",
		Content: "Hi {{name}}, thanks for connecting! I'd love to learn more about your work at {{company}}. Are you available for a quick chat sometime?",
		Variables: []string{"name", "company"},
	},
	{
		Name:    "follow_up_tech",
		Content: "Hello {{name}}, great to connect! I'm really interested in {{title}} roles and would appreciate any insights you might share about your experience at {{company}}.",
		Variables: []string{"name", "title", "company"},
	},
	{
		Name:    "follow_up_collaboration",
		Content: "Hi {{name}}, thanks for accepting my connection request! I noticed we have similar professional interests. I'd love to explore potential collaboration opportunities.",
		Variables: []string{"name"},
	},
	{
		Name:    "follow_up_learning",
		Content: "Hello {{name}}, I appreciate you connecting! Your background in {{title}} is impressive. I'm currently exploring this field and would value any advice you're willing to share.",
		Variables: []string{"name", "title"},
	},
}

type TemplateEngine struct {
	templates map[string]Template
}

func NewTemplateEngine() *TemplateEngine {
	engine := &TemplateEngine{
		templates: make(map[string]Template),
	}

	for _, template := range DefaultTemplates {
		engine.templates[template.Name] = template
	}

	return engine
}

func (te *TemplateEngine) AddTemplate(template Template) {
	te.templates[template.Name] = template
}

func (te *TemplateEngine) GetTemplate(name string) (Template, error) {
	template, exists := te.templates[name]
	if !exists {
		return Template{}, fmt.Errorf("template not found: %s", name)
	}

	return template, nil
}

func (te *TemplateEngine) RenderTemplate(templateName string, variables map[string]string) (string, error) {
	template, err := te.GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	content := template.Content

	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		content = strings.ReplaceAll(content, placeholder, value)
	}

	if strings.Contains(content, "{{") {
		return content, fmt.Errorf("warning: unresolved variables in template")
	}

	return content, nil
}

func (te *TemplateEngine) RenderForConnection(conn *storage.ConnectionRequest, templateName string) (string, error) {
	variables := map[string]string{
		"name":    extractFirstName(conn.Name),
		"title":   conn.Title,
		"company": conn.Company,
	}

	return te.RenderTemplate(templateName, variables)
}

func extractFirstName(fullName string) string {
	if fullName == "" {
		return "there"
	}

	parts := strings.Split(fullName, " ")
	if len(parts) > 0 {
		return parts[0]
	}

	return fullName
}

func (te *TemplateEngine) SelectBestTemplate(conn *storage.ConnectionRequest) string {
	if conn.Title != "" && conn.Company != "" {
		return "follow_up_tech"
	}

	if conn.Company != "" {
		return "follow_up_general"
	}

	if conn.Title != "" {
		return "follow_up_learning"
	}

	return "follow_up_collaboration"
}

func (te *TemplateEngine) ListTemplates() []string {
	var names []string
	for name := range te.templates {
		names = append(names, name)
	}
	return names
}

func (te *TemplateEngine) ValidateVariables(templateName string, variables map[string]string) error {
	template, err := te.GetTemplate(templateName)
	if err != nil {
		return err
	}

	missingVars := []string{}
	for _, requiredVar := range template.Variables {
		if _, exists := variables[requiredVar]; !exists {
			missingVars = append(missingVars, requiredVar)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("missing required variables: %v", missingVars)
	}

	return nil
}
