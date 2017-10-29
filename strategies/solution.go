package strategies

type SolutionStep struct {
	Solutions []Solution
}

type Solution struct {
	Steps   []SolutionStep
	Outcome float32
	Cost    float32
	Risk    uint
}
