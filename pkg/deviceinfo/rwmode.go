package deviceinfo

type RWMode byte

const (
	R RWMode = 1 << iota
	W
)

const RW = R | W

func (t RWMode) HasFlag(flag RWMode) bool {
	return t&flag != 0
}
