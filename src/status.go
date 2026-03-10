package src

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type TestResult struct {
	Status string    `json:"status"`
	Time   time.Time `json:"time"`
}

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
	TestResults      TestResult     `json:"test_results"`
	LastPowerEvent   string         `json:"last_power_event"`
}

func splitAssignUnitValue(s string) UnitValueInt {
	s = strings.TrimSpace(s)
	kv := strings.Split(s, " ")
	if len(kv) != 2 {
		log.Printf("Failed to split assign unit value: %s", s)
		return UnitValueInt{}
	}
	kVal, err := strconv.Atoi(kv[0])
	if err != nil {
		log.Printf("Failed to parse assign unit value: %s", s)
		return UnitValueInt{}
	}
	return UnitValueInt{
		Unit:  kv[1],
		Value: kVal,
	}
}

func splitCompoundUnitValue(s string) (string, string) {
	nLine := strings.Split(s, "(")
	val1 := strings.TrimSpace(nLine[0])
	val2 := strings.Trim(nLine[1], "()")
	return val1, val2
}

// valAsString splits the lines and returns clean end values without whitespace.
// Some values with compound data might require further processing.
func valAsString(s string, delimiter string) string {
	return strings.TrimSpace(strings.Split(s, delimiter)[1])
}

func ParseStatusStdOut(input string) (StatusCommand, error) {
	const _delimiter = ". "
	const _redelimiter = "_$"

	var properties = Properties{}
	var currentStatus = CurrentStatus{}

	// create unique key value delimiter
	input = strings.Replace(input, _delimiter, _redelimiter, -1)
	// clear all the annoying "..."
	input = strings.Replace(input, ".", "", -1)

	for _, line := range strings.Split(input, "\n") {
		// skip line if not content relevant
		if !strings.Contains(line, _redelimiter) {
			continue
		}
		key := strings.TrimSpace(strings.Split(line, _redelimiter)[0])
		value := strings.TrimSpace(strings.Split(line, _redelimiter)[1])
		fmt.Println(key, "->", value)

		switch strings.ToLower(key) {
		case "model name":
			properties.ModelName = value
		case "firmware number":
			properties.Firmware = value
		case "rating voltage":
			properties.RatingVoltage = splitAssignUnitValue(value)
		case "rating power":
			// reduced to "900 Watt(1500 VA)
			v1, v2 := splitCompoundUnitValue(value)
			properties.RatingPower = make([]UnitValueInt, 2)
			properties.RatingPower[0] = splitAssignUnitValue(v1)
			properties.RatingPower[1] = splitAssignUnitValue(v2)
		case "state":
			currentStatus.State = value
		case "power supply by":
			currentStatus.PowerSuppliedBy = value
		case "utility voltage":
			currentStatus.UtilityVoltage = splitAssignUnitValue(value)
		case "output voltage":
			currentStatus.OutputVoltage = splitAssignUnitValue(value)
		case "battery capacity":
			currentStatus.BatteryCapacity = splitAssignUnitValue(value)
		case "load":
			// reduced to "252 Watt(28 %)"
			v1, v2 := splitCompoundUnitValue(value)
			currentStatus.Load = make([]UnitValueInt, 2)
			currentStatus.Load[0] = splitAssignUnitValue(v1)
			currentStatus.Load[1] = splitAssignUnitValue(v2)
		case "line interaction":
			currentStatus.LineInteraction = value

			// TODO: Do better man.
		case "test results":
			const cursed = " at "
			// kind of cursed but I do not know what
			if !strings.Contains(line, cursed) {
				currentStatus.TestResults = TestResult{
					Status: "NA",
					Time:   time.Now(),
				}
			}
			tLine := strings.Split(valAsString(line, _redelimiter), cursed)
			pTime, err := time.Parse("2026/05/08 00:00:00", tLine[1])
			if err != nil {
				log.Printf("Failed to parse test results time: %s", tLine[1])
			}
			currentStatus.TestResults = TestResult{
				Status: strings.ToLower(tLine[0]),
				Time:   pTime,
			}
		case "last_power_event":
			currentStatus.LastPowerEvent = value
		}
	}
	return StatusCommand{
		Properties:    properties,
		CurrentStatus: currentStatus,
	}, nil
}

type StatusCommand struct {
	Properties
	CurrentStatus
}
