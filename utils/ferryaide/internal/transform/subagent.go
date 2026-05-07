package transform

import (
	"fmt"
	"strconv"
	"strings"
)

type SubAgent struct {
	Name        string
	Description string
	Prompt      string
}

func MarkdownSubAgentToTOML(input []byte) ([]byte, error) {
	agent, err := parseMarkdownSubAgent(string(input))
	if err != nil {
		return nil, err
	}
	var builder strings.Builder
	builder.WriteString("name = ")
	builder.WriteString(strconv.Quote(agent.Name))
	builder.WriteString("\n")
	builder.WriteString("description = ")
	builder.WriteString(strconv.Quote(agent.Description))
	builder.WriteString("\n")
	builder.WriteString("prompt = ")
	builder.WriteString(strconv.Quote(agent.Prompt))
	builder.WriteString("\n")
	return []byte(builder.String()), nil
}

func parseMarkdownSubAgent(input string) (SubAgent, error) {
	input = strings.TrimLeft(input, "\ufeff\r\n\t ")
	if !strings.HasPrefix(input, "---") {
		return SubAgent{}, fmt.Errorf("missing front matter")
	}
	parts := strings.SplitN(input, "---", 3)
	if len(parts) < 3 {
		return SubAgent{}, fmt.Errorf("unterminated front matter")
	}
	fields := parseFrontMatter(parts[1])
	agent := SubAgent{
		Name:        fields["name"],
		Description: fields["description"],
		Prompt:      strings.TrimSpace(parts[2]),
	}
	if agent.Name == "" {
		return SubAgent{}, fmt.Errorf("front matter field name is required")
	}
	return agent, nil
}

func parseFrontMatter(input string) map[string]string {
	fields := map[string]string{}
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		fields[strings.TrimSpace(key)] = strings.Trim(strings.TrimSpace(value), `"'`)
	}
	return fields
}
