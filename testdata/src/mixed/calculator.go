package mixed // want package:"testalign source order"

type Calculator struct{}

func (c *Calculator) Add(a, b int) int { return a + b }

func (c *Calculator) Subtract(a, b int) int { return a - b }

func (c *Calculator) Multiply(a, b int) int { return a * b }
