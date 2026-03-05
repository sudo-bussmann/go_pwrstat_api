package src

import (
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
	val, err := strconv.Atoi(kv[1])
	if err != nil {
		log.Printf("Failed to parse assign unit value: %s", s)
		return UnitValueInt{}
	}
	return UnitValueInt{
		Unit:  kv[1],
		Value: val,
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

func parseStatusStdOut(input string) (StatusCommand, error) {
	const _delimiter = ". "
	const _redelimiter = "_$"

	var properties = Properties{}
	var currentStatus = CurrentStatus{}

	input = strings.Replace(input, _delimiter, _redelimiter, -1)

	for _, line := range strings.Split(input, "\n") {
		// skip line if not content relevant
		if !strings.Contains(line, _redelimiter) {
			continue
		}
		switch strings.ToLower(line) {
		case "model name":
			properties.ModelName = valAsString(line, _redelimiter)
		case "firmware number":
			properties.Firmware = valAsString(line, _redelimiter)
		case "rating voltage":
			properties.RatingVoltage = splitAssignUnitValue(
				valAsString(line, _redelimiter),
			)
		case "rating power":
			// reduced to "900 Watt(1500 VA)
			v1, v2 := splitCompoundUnitValue(
				valAsString(line, _redelimiter),
			)
			properties.RatingPower = make([]UnitValueInt, 2)
			properties.RatingPower[0] = splitAssignUnitValue(v1)
			properties.RatingPower[1] = splitAssignUnitValue(v2)
		case "state":
			currentStatus.State = valAsString(line, _redelimiter)
		case "power supply by":
			currentStatus.PowerSuppliedBy = valAsString(line, _redelimiter)
		case "utility voltage":
			currentStatus.UtilityVoltage = splitAssignUnitValue(
				valAsString(line, _redelimiter),
			)
		case "output voltage":
			currentStatus.OutputVoltage = splitAssignUnitValue(
				valAsString(line, _redelimiter),
			)
		case "battery capacity":
			currentStatus.BatteryCapacity = splitAssignUnitValue(
				valAsString(line, _redelimiter),
			)
		case "load":
			// reduced to "252 Watt(28 %)"
			v1, v2 := splitCompoundUnitValue(
				valAsString(line, _redelimiter),
			)
			currentStatus.Load = make([]UnitValueInt, 2)
			currentStatus.Load[0] = splitAssignUnitValue(v1)
			currentStatus.Load[1] = splitAssignUnitValue(v2)
		case "line interaction":
			currentStatus.LineInteraction = valAsString(line, _redelimiter)
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
			currentStatus.LastPowerEvent = valAsString(line, _redelimiter)
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
