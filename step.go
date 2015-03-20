package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	ModifierTime = "ModifierTime"
	ModifierTemp = "ModifierTemp"
)

const (
	ToolContainer = "ToolContainer"
	ToolHeater    = "ToolHeater"
	ToolCooler    = "ToolCooler"
	ToolMixer     = "ToolMixer"
	ToolCutter    = "ToolCutter"
	ToolFilter    = "ToolFilter"
	ToolCoater    = "ToolCoater"
)

const (
	// default tool: Heater
	ActionHeatHigh       = "ActionHeatHigh"
	ActionHeatMediumHigh = "ActionHeatMediumHigh"
	ActionHeatMedium     = "ActionHeatMedium"
	ActionHeatMediumLow  = "ActionHeatMediumLow"
	ActionHeatLow        = "ActionHeatLow"
	// default tool: Cooler
	ActionCool = "ActionCool"
	// default tool: Hand
	ActionPlaceIn      = "ActionPlaceIn"
	ActionPlaceSide    = "ActionPlaceSide"
	ActionPlaceSurface = "ActionPlaceSurface"
	// default tool: Cutter
	ActionToChunk    = "ActionToChunk"    // bigger chunk
	ActionToSlice    = "ActionToSlice"    // flat
	ActionToParticle = "ActionToParticle" // tiny particle
	ActionToBar      = "ActionToBar"      // bar
	// default tool: Hand
	ActionCombine  = "ActionCombine"  // into 1 object
	ActionSeparate = "ActionSeparate" // into 2 objects
)

const (
	ThingProduced       = "ThingProduced" // produced thing
	ThingProducedSolid  = "ThingProducedSolid"
	ThingProducedLiquid = "ThingProducedLiquid"

	ThingLiquid = "ThingLiquid"

	ThingSolid         = "ThingSolid"
	ThingSolidChunk    = "ThingSolidChunk"
	ThingSolidSlice    = "ThingSolidSlice"
	ThingSolidParticle = "ThingSolidParticle"
	ThingSolidBar      = "ThingSolidBar"
)

type Step interface {
	Action() string
	Input() []string
	Output() []string
	Tool() []string
	Skip() bool
	String() string
}

// object has no state, state is inferred from step history
// (From,To) is either one-to-many or many-to-one
type BasicStep struct {
	Do        string
	FromThing []string // operands, thing
	ToThing   []string // returned, thing
	WithTool  []string // used to perform step
	Optional  bool
}

func (s *BasicStep) dump() []string {
	var word []string
	if s.Optional {
		word = append(word, "OPTIONALLY")
	}
	switch s.Do {
	case ActionHeatHigh:
		word = append(word, "COOK", "input", "AT HIGH")
	case ActionHeatMediumHigh:
		word = append(word, "COOK", "input", "AT MEDIUM HIGH")
	case ActionHeatMedium:
		word = append(word, "COOK", "input")
	case ActionHeatMediumLow:
		word = append(word, "COOK", "input", "AT MEDIUM LOW")
	case ActionHeatLow:
		word = append(word, "COOK", "input", "AT LOW")
	case ActionCool:
		word = append(word, "COOL", "input")
	case ActionPlaceIn:
		word = append(word, "PUT", "input1", "IN", "input2")
	case ActionPlaceSide:
		word = append(word, "PUT", "input1", "NEXT TO", "input2")
	case ActionPlaceSurface:
		word = append(word, "PUT", "input1", "ON THE SURFACE OF", "input2")
	case ActionToSlice:
		word = append(word, "CUT", "input", "INTO SLICES")
	case ActionToChunk:
		word = append(word, "CUT", "input", "INTO CHUNKS")
	case ActionToParticle:
		word = append(word, "GRIND", "input")
	case ActionToBar:
		word = append(word, "CUT", "input", "INTO BARS")
	case ActionCombine:
		word = append(word, "MIX", "input", "INTO", "output")
	case ActionSeparate:
		word = append(word, "SEPARATE", "input", "INTO", "output")
	}

	for i := range word {
		switch word[i] {
		case "input1":
			if len(s.FromThing) > 0 {
				word[i] = s.FromThing[0]
			}
		case "input2":
			if len(s.FromThing) > 1 {
				word[i] = s.FromThing[1]
			}
		case "input":
			if len(s.FromThing) > 0 {
				word[i] = strings.Join(s.FromThing, ",")
			}
		case "output":
			if len(s.ToThing) > 0 {
				word[i] = strings.Join(s.ToThing, ",")
			}
		}
	}

	if len(s.WithTool) > 0 {
		word = append(word, "WITH", strings.Join(s.WithTool, ","))
	}
	word = append(word, "("+strings.Join(s.ToThing, ",")+")")
	return word
}

func (s *BasicStep) String() string {
	return strings.Join(s.dump(), " ") + "."
}

func (s *BasicStep) Action() string {
	return s.Do
}

func (s *BasicStep) Input() []string {
	return s.FromThing
}

func (s *BasicStep) Output() []string {
	return s.ToThing
}

func (s *BasicStep) Tool() []string {
	return s.WithTool
}

func (s *BasicStep) Skip() bool {
	return s.Optional
}

func (s *ContinueStep) dump() []string {
	word := s.BasicStep.dump()
	if s.Duration != 0 {
		word = append(word, "FOR", s.Duration.String())
	}
	if s.Temp > 0 {
		word = append(word, "TIL", "TEMPERATURE", "REACHES", fmt.Sprintf("%d°F", s.Temp))
	}
	return word
}

func (s *ContinueStep) String() string {
	return strings.Join(s.dump(), " ") + "."
}

// do until stop condition
type ContinueStep struct {
	BasicStep

	// stop after time
	Duration time.Duration
	// stop after reaching temperature
	Temp int
}
