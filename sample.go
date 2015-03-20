package main

import "time"

var (
	Input1 = []string{
		"In a large stock pot, bring water, carrots, potatoes, onion, salsa, and bouillon cubes to a boil.",
		"Reduce to a medium simmer with 200 degrees F, stirring occasionally, approximately 10 minutes.",
		"Mix the beef, breadcrumbs, and milk together in a bowl.",
		"Form into 1-inch meatballs, and drop into boiling broth.",
		"Once soup returns to a boil, reduce heat to medium-low.",
		"Cover and cook 20 minutes, or until meatballs are no longer pink in center, and vegetables are tender.",
		"Serve with sprinkled cilantro for garnish.",
	}

	// Uppercase for combined thing
	// http://allrecipes.com/Recipe/Albondigas/Detail.aspx?soid=carousel_0_rotd&prop24=rotd
	Instr1 = Instruction{
		// In a large stock pot, bring water, carrots, potatoes, onion, salsa, and bouillon cubes to a boil.
		&BasicStep{
			Do:        ActionCombine,
			FromThing: []string{"water", "carrot", "potatoe", "onion", "salsa", "bouillon_cube"},
			ToThing:   []string{"POT_FOOD_1"},
			WithTool:  []string{ToolContainer},
		},
		&ContinueStep{
			BasicStep: BasicStep{
				Do:        ActionHeatHigh,
				FromThing: []string{"POT_FOOD_1"},
				WithTool:  []string{ToolContainer},
			},
			Temp: 100,
		},
		// Reduce to a medium simmer, stirring occasionally, approximately 10 minutes.
		&ContinueStep{
			BasicStep: BasicStep{
				Do:        ActionHeatMedium,
				FromThing: []string{"POT_FOOD_1"},
				WithTool:  []string{ToolContainer, ToolMixer},
			},
			Duration: 10 * time.Minute,
		},
		// Mix the beef, breadcrumbs, and milk together in a bowl.
		&BasicStep{
			Do:        ActionCombine,
			FromThing: []string{"beef", "breadcrumb", "milk"},
			ToThing:   []string{"BOWL_FOOD_1"},
			WithTool:  []string{ToolMixer, ToolContainer},
		},
		// Form into 1-inch meatballs, and drop into boiling broth.
		&BasicStep{
			Do:        ActionToChunk,
			FromThing: []string{"BOWL_FOOD_1"},
		},
		&BasicStep{
			Do:        ActionCombine,
			FromThing: []string{"POT_FOOD_1", "BOWL_FOOD_1"},
			ToThing:   []string{"POT_FOOD_2"},
			WithTool:  []string{ToolContainer},
		},
		// Once soup returns to a boil, reduce heat to medium-low.
		&ContinueStep{
			BasicStep: BasicStep{
				Do:        ActionHeatMedium,
				FromThing: []string{"POT_FOOD_2"},
				WithTool:  []string{ToolContainer},
			},
			Temp: 100,
		},
		&BasicStep{
			Do:        ActionHeatMediumLow,
			FromThing: []string{"POT_FOOD_2"},
			WithTool:  []string{ToolContainer},
		},
		// Cover and cook 20 minutes, or until meatballs are no longer pink in center, and vegetables are tender.
		&ContinueStep{
			BasicStep: BasicStep{
				Do:        ActionHeatMedium,
				FromThing: []string{"POT_FOOD_2"},
				WithTool:  []string{ToolContainer},
			},
			Duration: 20 * time.Minute,
		},
		// Serve with sprinkled cilantro for garnish.
		&BasicStep{
			Do:        ActionToParticle,
			FromThing: []string{"cilantro"},
			Optional:  true,
		},
		&BasicStep{
			Do:        ActionPlaceSurface,
			FromThing: []string{"cilantro", "POT_FOOD_2"},
			Optional:  true,
		},
	}
)
