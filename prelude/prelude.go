// prelude takes the most common usage from the library and set them up accordingly in simple steps.

package prelude

type GlobalOption struct {
	// Service name.
	Service string
	// Scope for this service. When the application have multiple commands,
	// filling the scope value helps points out the origin of tracing from which command.
	Scope string
	// The environment name.
	Environment string
}

func InitGlobal(opts GlobalOption) {}
