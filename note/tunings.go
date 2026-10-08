package note

type noteLabel struct {
	Name   string
	Octave int
	Label  string
}

type Tuning struct {
	Name  string
	Notes []noteLabel
}

var Standard = Tuning{
	Name: "Standard",
	Notes: []noteLabel{
		{"E", 2, "E"}, {"A", 2, "A"}, {"D", 3, "D"},
		{"G", 3, "G"}, {"B", 3, "B"}, {"E", 4, "e"},
	},
}

var DropD = Tuning{
	Name: "Drop D",
	Notes: []noteLabel{
		{"D", 2, "D"}, {"A", 2, "A"}, {"D", 3, "D"},
		{"G", 3, "G"}, {"B", 3, "B"}, {"E", 4, "e"},
	},
}

var Eb = Tuning{
	Name: "Eb Standard",
	Notes: []noteLabel{
		{"Eb", 2, "Eb"}, {"Ab", 2, "Ab"}, {"Db", 3, "Db"},
		{"Gb", 3, "Gb"}, {"Bb", 3, "Bb"}, {"Eb", 4, "eb"},
	},
}

var DropC = Tuning{
	Name: "Drop C",
	Notes: []noteLabel{
		{"C", 2, "C"}, {"G", 2, "G"}, {"C", 3, "C"},
		{"F", 3, "F"}, {"A", 3, "A"}, {"D", 4, "d"},
	},
}

var DADGAD = Tuning{
	Name: "DADGAD",
	Notes: []noteLabel{
		{"D", 2, "D"}, {"A", 2, "A"}, {"D", 3, "D"},
		{"G", 3, "G"}, {"A", 3, "A"}, {"D", 4, "d"},
	},
}

var OpenG = Tuning{
	Name: "Open G",
	Notes: []noteLabel{
		{"D", 2, "D"}, {"G", 2, "G"}, {"D", 3, "D"},
		{"G", 3, "G"}, {"B", 3, "B"}, {"D", 4, "d"},
	},
}

var OpenD = Tuning{
	Name: "Open D",
	Notes: []noteLabel{
		{"D", 2, "D"}, {"A", 2, "A"}, {"D", 3, "D"},
		{"F#", 3, "F#"}, {"A", 3, "A"}, {"D", 4, "d"},
	},
}

var DStandard = Tuning{
	Name: "D Standard",
	Notes: []noteLabel{
		{"D", 2, "D"}, {"G", 2, "G"}, {"C", 3, "C"},
		{"F", 3, "F"}, {"A", 3, "A"}, {"D", 4, "d"},
	},
}

var CStandard = Tuning{
	Name: "C Standard",
	Notes: []noteLabel{
		{"C", 2, "C"}, {"F", 2, "F"}, {"Bb", 2, "Bb"},
		{"Eb", 3, "Eb"}, {"G", 3, "G"}, {"C", 4, "c"},
	},
}

var BStandard = Tuning{
	Name: "B Standard",
	Notes: []noteLabel{
		{"B", 1, "B"}, {"E", 2, "E"}, {"A", 2, "A"},
		{"D", 3, "D"}, {"F#", 3, "F#"}, {"B", 3, "b"},
	},
}
var Tunings = [10]Tuning{ Standard, DropD, Eb, DropC, DADGAD, OpenG, OpenD, DStandard, CStandard, BStandard, }