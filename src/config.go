package src

import (
	"fmt"
	"strings"
)

type DaemonConfig struct {
	Alarm     string `json:"alarm"`
	Hibernate string `json:"hibernate"`
	Cloud     string `json:"cloud"`
}

type PowerFailureAction struct {
	DelayTimeSincePowerFailure UnitValueInt `json:"delay_time_since_power_failure"`
	RunScriptCommandPF         string       `json:"run_script_command_pf"`
	PathOfScriptCommand        string       `json:"path_of_script_command"`
	DurationOfCommandRunningPF UnitValueInt `json:"duration_of_command_running_pf"`
	EnableShutdownSystemPF     string       `json:"enable_shutdown_system_pf"`
}

type BatteryLowAction struct {
	RemainingRuntimeThreshold  UnitValueInt `json:"remaining_runtime_threshold"`
	BatteryCapacityThreshold   UnitValueInt `json:"battery_capacity_threshold"`
	RunScriptCommandBL         string       `json:"run_script_command_bl"`
	PathOfCommand              string       `json:"path_of_command"`
	DurationOfCommandRunningBL UnitValueInt `json:"duration_of_command_running_bl"`
	EnableShutdownSystemBL     string       `json:"enable_shutdown_system_bl"`
}

type ConfigurationCLIResp struct {
	DaemonConfig DaemonConfig       `json:"daemon_config"`
	PowerFailure PowerFailureAction `json:"power_failure"`
	BatteryLow   BatteryLowAction   `json:"battery_low"`
}

func ParseConfigStdOut(input string) (ConfigurationCLIResp, error) {
	const _delimiter = ". "
	const _redelimiter = "_$"

	var daemonConf DaemonConfig
	var powerFailAction PowerFailureAction
	var batteryLowAction BatteryLowAction

	// create unique key value delimiter
	input = strings.Replace(input, _delimiter, _redelimiter, -1)

	// tmp change of file extentions to guard next part breaking file exts
	input = strings.Replace(input, ".sh", "_sh", -1)
	// clear all the annoying "..."
	input = strings.Replace(input, ".", "", -1)
	// restore file extensions
	input = strings.Replace(input, "_sh", ".sh", -1)

	// kinda annoying, but it correctly assigns to the Action for Power Failure values
	input = strings.Replace(input, "Run script command", "aRun script command", 1)
	input = strings.Replace(input, "Duration of command running", "aDuration of command running", 1)
	input = strings.Replace(input, "Enable shutdown system", "aEnable shutdown system", 1)
	for _, line := range strings.Split(input, "\n") {
		// skip line if not content relevant
		if !strings.Contains(line, _redelimiter) {
			continue
		}
		key := strings.TrimSpace(strings.Split(line, _redelimiter)[0])
		value := strings.TrimSpace(strings.Split(line, _redelimiter)[1])
		fmt.Println(key, "->", value)
		switch strings.ToLower(key) {
		case "alarm":
			daemonConf.Alarm = value
		case "hibernate":
			daemonConf.Hibernate = value
		case "cloud":
			daemonConf.Cloud = value
		case "delay time since power failure":
			powerFailAction.DelayTimeSincePowerFailure = splitAssignUnitValue(value)
		case "arun script command":
			powerFailAction.RunScriptCommandPF = value
		case "path of script command":
			powerFailAction.PathOfScriptCommand = value
		case "aduration of command running":
			powerFailAction.DurationOfCommandRunningPF = splitAssignUnitValue(value)
		case "aenable shutdown system":
			powerFailAction.EnableShutdownSystemPF = value
		case "remaining runtime threshold":
			batteryLowAction.RemainingRuntimeThreshold = splitAssignUnitValue(value)
		case "battery capacity threshold":
			batteryLowAction.BatteryCapacityThreshold = splitAssignUnitValue(value)
		case "run script command":
			batteryLowAction.RunScriptCommandBL = value
		case "path of command":
			batteryLowAction.PathOfCommand = value
		case "duration of command running":
			batteryLowAction.DurationOfCommandRunningBL = splitAssignUnitValue(value)
		case "enable shutdown system":
			batteryLowAction.EnableShutdownSystemBL = value
		default:
			err := fmt.Errorf("Unknown configuration key: %s", key)
			return ConfigurationCLIResp{}, err
		}
	}
	return ConfigurationCLIResp{
		daemonConf,
		powerFailAction,
		batteryLowAction,
	}, nil
}
