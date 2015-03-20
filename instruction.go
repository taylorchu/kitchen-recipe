package main

type Instruction []Step

func (instr Instruction) Thing() []string {
	m := make(map[string]struct{})
	for _, s := range instr {
		for _, t := range s.Input() {
			m[t] = struct{}{}
		}
		for _, t := range s.Output() {
			m[t] = struct{}{}
		}
	}
	var ts []string
	for t := range m {
		ts = append(ts, t)
	}
	return ts
}

func (instr Instruction) Ref() [][]int {
	d := make([][]int, len(instr))
	for i := 0; i < len(instr); i++ {
		m := make(map[string]struct{})
		for _, t := range instr[i].Input() {
			m[t] = struct{}{}
		}
		for j := i - 1; j >= 0; j-- {
			prev := instr[j]
			var add bool
			for _, t := range prev.Input() {
				if _, ok := m[t]; ok {
					delete(m, t)
					add = true
				}
			}
			for _, t := range prev.Output() {
				if _, ok := m[t]; ok {
					delete(m, t)
					add = true
				}
			}
			if add {
				d[i] = append(d[i], j)
			}
		}
	}
	return d
}

func (instr Instruction) AllRef() [][]int {
	d := make([][]int, len(instr))
	r := instr.Ref()
	var findRef func(int, int) bool
	findRef = func(i, key int) bool {
		for _, s := range r[i] {
			if s == key {
				return true
			}
			if findRef(s, key) {
				return true
			}
		}
		return false
	}
	for i := 0; i < len(instr); i++ {
		for j := i - 1; j >= 0; j-- {
			if findRef(i, j) {
				d[i] = append(d[i], j)
			}
		}
	}
	return d
}
