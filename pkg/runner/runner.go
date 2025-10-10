package runner



type Runner interface {
	Run(input string) (any, error)
}
