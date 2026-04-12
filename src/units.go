package src

import (
	"log"
	"strconv"
	"strings"
)

type UnitValueInt struct {
	Value int    `json:"value"`
	Unit  string `json:"unit"`
}

type UnitValueFloat struct {
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
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
