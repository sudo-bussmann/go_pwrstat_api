package src

type DaemonConfig struct {
	Alarm     string `json:"alarm"`
	Hibernate string `json:"hibernate"`
	Cloud     string `json:"cloud"`
}

type PowerFailureAction struct {
	DelayTimeSincePowerFailure UnitValueInt `json:"delay_time_since_power_failure"`
	RunScriptCommandPF         string       `json:"run_script_command"`
	PathOfScriptCommand        string       `json:"path_of_script_command"`
	DurationOfCommandRunningPF UnitValueInt `json:"duration_of_command_running"`
	EnableShutdownSystem       string       `json:"enable_shutdown_system"`
}

type BatteryLowAction struct {
	RemainingRuntimeThreshold  string       `json:"remaining_runtime_threshold"`
	BatteryCapacityThreshold   UnitValueInt `json:"battery_capacity_threshold"`
	RunScriptCommandBL         string       `json:"run_script_command"`
	PathOfCommand              string       `json:"path_of_command"`
	DurationOfCommandRunningBL UnitValueInt `json:"duration_of_command_running"`
	EnableShutdownSystem       string       `json:"enable_shutdown_system"`
}

type ConfigurationCLIResp struct {
	DaemonConfig
	PowerFailureAction
	BatteryLowAction
}

func NewConfigurationResp() *ConfigurationCLIResp {
	return &ConfigurationCLIResp{
		DaemonConfig:       DaemonConfig{},
		PowerFailureAction: PowerFailureAction{},
		BatteryLowAction:   BatteryLowAction{},
	}
}
