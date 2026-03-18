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
	input = strings.Replace(input, "Path of script command", "aPath	of script command", 1)
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
			// finish implementation
		default:

		}
	}
	return ConfigurationCLIResp{daemonConf, powerFailAction, batteryLowAction}, nil
}
