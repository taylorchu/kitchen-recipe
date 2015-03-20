package nl

import "unicode"

const (
	Prep = "Prep"

	Conj = "Conj"

	Unit = "Unit"

	Noun = "Noun"
	Adj  = "Adj"

	Verb = "Verb"
	Adv  = "Adv"

	Comma  = "Comma"
	Period = "Period"
)

var (
	dict = map[string][]string{
		// exhaustive
		".": {Period},
		",": {Comma},

		// exhaustive
		"minute":   {Unit},
		"F":        {Unit},
		"C":        {Unit},
		"degree":   {Unit},
		"cup":      {Unit},
		"teaspoon": {Unit},
		"inch":     {Unit},
		"pound":    {Unit},

		// exhaustive
		"in":      {Prep},
		"to":      {Prep},
		"for":     {Prep},
		"with":    {Prep},
		"on":      {Prep},
		"at":      {Prep},
		"from":    {Prep},
		"by":      {Prep},
		"about":   {Prep},
		"as":      {Prep},
		"into":    {Prep},
		"through": {Prep},
		"over":    {Prep},
		"between": {Prep},
		"without": {Prep},
		"against": {Prep},
		"during":  {Prep},
		"under":   {Prep},
		"around":  {Prep},
		"among":   {Prep},

		"large":         {Adj},
		"stock":         {Noun},
		"pot":           {Noun},
		"bring":         {Verb},
		"water":         {Noun},
		"carrot":        {Noun},
		"potato":        {Noun},
		"onion":         {Noun},
		"salsa":         {Noun},
		"bouillon":      {Noun},
		"cube":          {Noun},
		"boil":          {Noun, Verb},
		"reduce":        {Verb},
		"medium":        {Adj},
		"simmer":        {Noun, Verb},
		"stirring":      {Verb, Adj},
		"occasionally":  {Adv},
		"approximately": {Adv},
		"mix":           {Verb, Noun},
		"beef":          {Noun},
		"breadcrumb":    {Noun},
		"milk":          {Noun},
		"together":      {Adv},
		"bowl":          {Noun},
		"form":          {Verb, Noun},
		"meatball":      {Noun},
		"and":           {Conj},
		"drop":          {Verb},
		"boiling":       {Verb, Adj},
		"broth":         {Noun},
		"once":          {Conj, Adv},
		"soup":          {Noun},
		"return":        {Verb},
		"heat":          {Noun},
		"medium_low":    {Noun},
		"cover":         {Verb, Noun},
		"cook":          {Verb, Noun},
		"or":            {Conj},
		"until":         {Prep, Conj},
		"are":           {Verb},
		"no":            {Adv},
		"longer":        {Adv, Adj},
		"pink":          {Adj},
		"center":        {Noun},
		"vegetable":     {Noun},
		"tender":        {Adj},
		"serve":         {Verb},
		"sprinkled":     {Adj},
		"cilantro":      {Noun},
		"garnish":       {Verb, Noun},
	}
)

func containsDigit(s string) bool {
	for _, r := range s {
		if unicode.IsDigit(r) {
			return true
		}
	}
	return false
}

func TagPOS(s string) []string {
	if containsDigit(s) {
		return []string{Unit}
	}
	return dict[s]
}
