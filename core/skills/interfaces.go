package skills

// SkillDefinition is the canonical contract exposed to the engine and LLM.
type SkillDefinition struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Actions     []ActionDefinition `json:"actions"`
}

// ActionDefinition describes one callable action for a skill.
type ActionDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema string `json:"input_schema"`
}
