package siren

type Status int

func (s Status) String() string {
	switch s {
	case ok:
		return "OK"
	case warn:
		return "WARN"
	default:
		return "FAIL"
	}
}

const (
	ok Status = iota
	warn
	fail
)
