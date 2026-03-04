package src

type UnitValueInt struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type UnitValueFloat struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
}
