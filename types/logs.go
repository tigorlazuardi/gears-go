package types

// An interface similar to Fielder, but this one is intended to also hold more confidential values
//
// The user of this interface must be aware of data breaches. So ensure the data never goes out to wrong party.
type Logs interface {
	Logs() Fields
}
