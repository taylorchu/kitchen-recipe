package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/taylorchu/blueprint/nl"
)

var (
	foodCount = 0
)

func newFood() string {
	foodCount++
	return fmt.Sprintf("FOOD_%d", foodCount)
}

func ParseInstruction(steps []*ContinueStep) (instr Instruction) {
	for i, step := range steps {
		// autocomplete FromThing
		if i > 0 {
			switch len(step.FromThing) {
			case 0:
				// prev output = curr input
				step.FromThing = steps[i-1].ToThing
			case 1:
				// placement needs 2 inputs; get the second from prev output
				switch step.Do {
				case ActionPlaceIn:
					fallthrough
				case ActionPlaceSide:
					fallthrough
				case ActionPlaceSurface:
					step.FromThing = append(step.FromThing, steps[i-1].ToThing[0])
				}
			}
		}

		// autocomplete ToThing
		if len(step.ToThing) == 0 {
			// the second input of placement action can be output
			if len(step.FromThing) > 1 {
				switch step.Do {
				case ActionPlaceIn:
					fallthrough
				case ActionPlaceSide:
					fallthrough
				case ActionPlaceSurface:
					step.ToThing = []string{step.FromThing[1]}
				}
			}
		}

		if len(step.ToThing) == 0 {
			// if there is only a produced object in curr input, use it as curr output
			if len(step.FromThing) == 1 &&
				(strings.HasPrefix(Pub[step.FromThing[0]], ThingProduced) || strings.HasPrefix(step.FromThing[0], "FOOD_")) {
				step.ToThing = []string{step.FromThing[0]}
			}
		}

		if len(step.ToThing) == 0 {
			// look ahead! curr output is used in next input
		LOOK_AHEAD:
			for _, next := range steps[i+1:] {
				for _, t := range next.FromThing {
					if !strings.HasPrefix(Pub[t], ThingProduced) && !strings.HasPrefix(t, "FOOD_") {
						continue
					}
					for _, sub := range Sub[t] {
						if sub != step.Do {
							continue
						}
						step.ToThing = []string{t}
						break LOOK_AHEAD
					}
				}
			}
		}

		if len(step.ToThing) == 0 {
			// look behind! curr input is already assigned to some output
			track := make(map[string]struct{})
			for _, t := range step.FromThing {
				track[t] = struct{}{}
			}
		LOOK_BEHIND:
			for j := i - 1; j >= 0; j-- {
				prev := steps[j]
				for _, t := range prev.FromThing {
					if _, ok := track[t]; !ok {
						continue
					}
					if !strings.HasPrefix(Pub[prev.ToThing[0]], ThingProduced) && !strings.HasPrefix(prev.ToThing[0], "FOOD_") {
						continue
					}
					for _, sub := range Sub[prev.ToThing[0]] {
						if sub != step.Do {
							continue
						}
						step.ToThing = []string{prev.ToThing[0]}
						break LOOK_BEHIND
					}
				}
			}
		}

		if len(step.ToThing) == 0 {
			// cannot find any; use new symbol
			step.ToThing = []string{newFood()}
		}

		if step.Temp == 0 && step.Duration == 0 {
			instr = append(instr, &step.BasicStep)
		} else {
			instr = append(instr, step)
		}
	}
	return
}

func ParseStep(chunks [][]string) (steps []*ContinueStep) {
	published := make(map[string][]string)
	var actions []string
	for _, chunk := range chunks {
		var tokens []string
		var unitstr string
		for _, token := range chunk {
			parts := strings.SplitN(token, "/", 2)
			switch parts[1] {
			case nl.Prep:
				fallthrough
			case nl.Adv:
			case nl.Unit:
				switch parts[0] {
				case "hour":
					unitstr += "h"
				case "minute":
					unitstr += "m"
				case "second":
					unitstr += "s"
				case "degree":
				default:
					unitstr += parts[0]
				}
			default:
				tokens = append(tokens, parts[0])
			}
		}

		if unitstr != "" {
			if strings.HasSuffix(unitstr, "F") {
				published[ModifierTemp] = append(published[ModifierTemp], unitstr)
			} else {
				published[ModifierTime] = append(published[ModifierTime], unitstr)
			}
		}

		phrase := strings.Join(tokens, "_")
		if pub, ok := Pub[phrase]; ok {
			published[pub] = append(published[pub], phrase)
			if strings.HasPrefix(pub, "Action") {
				actions = append(actions, phrase)
			}
		}
	}
	for _, action := range actions {
		var s ContinueStep
		s.Do = Pub[action]
		for _, label := range Sub[s.Do] {
			tries := []string{label}
			switch label {
			case ThingSolid:
				tries = append(tries, ThingSolidBar, ThingSolidParticle, ThingSolidChunk, ThingSolidSlice)
			case ThingProduced:
				tries = append(tries, ThingProducedSolid, ThingProducedLiquid)
			}
			for _, try := range tries {
				for _, publisher := range published[try] {
					if strings.HasPrefix(try, "Tool") {
						s.WithTool = append(s.WithTool, publisher)
					}
					if strings.HasPrefix(try, "Thing") {
						s.FromThing = append(s.FromThing, publisher)
					}

					if strings.HasSuffix(publisher, "F") {
						if temp, err := strconv.Atoi(strings.TrimSuffix(publisher, "F")); err == nil {
							s.Temp = temp
						}
					}
					if duration, err := time.ParseDuration(publisher); err == nil {
						s.Duration = duration
					}
				}
			}
		}
		steps = append(steps, &s)
	}
	return
}

var (
	Pub = map[string]string{
		"large_stock_pot":    ToolContainer,
		"water":              ThingLiquid,
		"carrot":             ThingSolidBar,
		"potato":             ThingSolidChunk,
		"onion":              ThingSolidChunk,
		"bouillon_cube":      ThingSolidChunk,
		"boil":               ActionHeatHigh,
		"medium_simmer":      ActionHeatMedium,
		"stirring":           ActionCombine,
		"mix":                ActionCombine,
		"beef":               ThingSolidChunk,
		"breadcrumb":         ThingSolidParticle,
		"milk":               ThingLiquid,
		"bowl":               ToolContainer,
		"form":               ActionToChunk,
		"meatball":           ThingProducedSolid,
		"drop":               ActionCombine,
		"boiling_broth":      ThingProducedLiquid,
		"soup":               ThingProducedLiquid,
		"medium_low":         ActionHeatMediumLow,
		"cover":              ActionPlaceSurface,
		"cook":               ActionHeatMedium,
		"vegetable":          ThingSolid,
		"serve":              ActionPlaceSurface,
		"sprinkled_cilantro": ThingSolidParticle,
	}
	Sub = map[string][]string{
		ActionHeatHigh:       {ModifierTemp, ModifierTime, ToolHeater, ToolContainer, ThingProduced, ThingSolid, ThingLiquid},
		ActionHeatMediumHigh: {ModifierTemp, ModifierTime, ToolHeater, ToolContainer, ThingProduced, ThingSolid, ThingLiquid},
		ActionHeatMedium:     {ModifierTemp, ModifierTime, ToolHeater, ToolContainer, ThingProduced, ThingSolid, ThingLiquid},
		ActionHeatMediumLow:  {ModifierTemp, ModifierTime, ToolHeater, ToolContainer, ThingProduced, ThingSolid, ThingLiquid},
		ActionHeatLow:        {ModifierTemp, ModifierTime, ToolHeater, ToolContainer, ThingProduced, ThingSolid, ThingLiquid},

		ActionCool: {ModifierTemp, ModifierTime, ToolCooler, ToolContainer, ThingProduced, ThingSolid, ThingLiquid},

		ActionPlaceIn:      {ThingProducedSolid, ThingSolid},
		ActionPlaceSide:    {ThingProducedSolid, ThingSolid},
		ActionPlaceSurface: {ThingProducedSolid, ThingSolid, ThingLiquid},

		ActionToChunk:    {ToolCutter, ThingProducedSolid, ThingSolidChunk},
		ActionToSlice:    {ToolCutter, ThingProducedSolid, ThingSolidChunk},
		ActionToParticle: {ToolCutter, ThingProducedSolid, ThingSolidChunk, ThingSolidBar, ThingSolidSlice},
		ActionToBar:      {ToolCutter, ThingProducedSolid, ThingSolidChunk, ThingSolidSlice},

		ActionCombine:  {ToolMixer, ToolContainer, ThingProduced, ThingSolid, ThingLiquid},
		ActionSeparate: {ToolFilter, ThingProduced, ThingSolid, ThingLiquid},

		"meatball":      {ActionToChunk, ActionCombine},
		"boiling_broth": {ActionHeatHigh, ActionCombine},
		"soup":          {ActionHeatMedium, ActionCombine},
	}
)
