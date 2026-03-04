package src

type Properties struct {
	ModelName     string         `json:"model_name"`
	Firmware      string         `json:"firmware"`
	RatingVoltage UnitValueInt   `json:"rating_voltage"`
	RatingPower   []UnitValueInt `json:"rating_power"`
}

type CurrentStatus struct {
	State            string         `json:"state"`
	PowerSuppliedBy  string         `json:"power_supplied_by"`
	UtilityVoltage   UnitValueInt   `json:"utility_voltage"`
	OutputVoltage    UnitValueInt   `json:"output_voltage"`
	BatteryCapacity  UnitValueInt   `json:"battery_capacity"`
	RemainingRuntime UnitValueInt   `json:"remaining_runtime"`
	Load             []UnitValueInt `json:"load"`
	LineInteraction  string         `json:"line_interaction"`
	TestResults      string         `json:"test_results"`
	LastPowerEvent   string         `json:"last_power_event"`
}

type StatusCommand struct {
	Properties
	CurrentStatus
}
