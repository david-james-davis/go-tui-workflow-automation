package bubble

// DefaultChoice represents main menu choices
type DefaultChoice string

const (
	Choice1 DefaultChoice = "Choice 1"
	Choice2 DefaultChoice = "Choice 2"
	Choice3 DefaultChoice = "Choice 3"
)

// AllDefaultChoices returns all default menu choices
func AllDefaultChoices() []string {
	return []string{
		string(Choice1),
		string(Choice2),
		string(Choice3),
	}
}

// GoChoice represents Go template choices
type GoChoice string

const (
	GoTemplate1 GoChoice = "Go Template 1"
	GoTemplate2 GoChoice = "Go Template 2"
	GoTemplate3 GoChoice = "Go Template 3"
)

// AllGoChoices returns all Go template choices
func AllGoChoices() []string {
	return []string{
		string(GoTemplate1),
		string(GoTemplate2),
		string(GoTemplate3),
	}
}
