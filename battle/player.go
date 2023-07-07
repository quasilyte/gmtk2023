package battle

type player interface {
	Update(scaledDelta, delta float64)
}
