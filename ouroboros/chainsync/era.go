package chainsync

type Era struct {
	name string
}

var (
	Byron   = Era{name: "byron"}
	Shelley = Era{name: "shelley"}
	Allegra = Era{name: "allegra"}
	Mary    = Era{name: "mary"}
	Alonzo  = Era{name: "alonzo"}
	Babbage = Era{name: "babbage"}
)

var Eras = [...]Era{Byron, Shelley, Allegra, Mary, Alonzo, Babbage}

func (e Era) String() string {
	return e.name
}

func (e Era) AlonzoOrGreater() bool {
	alonzoIdx := -1
	for idx, era := range Eras {
		if era == Alonzo {
			alonzoIdx = idx
		}
	}

	for idx, era := range Eras {
		if e == era {
			return idx >= alonzoIdx
		}
	}

	panic("new era unaccounted for")
}

func (r RollForwardBlock) Era() Era {
	switch {
	case r.Byron != nil:
		return Byron
	case r.Allegra != nil:
		return Allegra
	case r.Alonzo != nil:
		return Alonzo
	case r.Mary != nil:
		return Mary
	case r.Shelley != nil:
		return Shelley
	case r.Babbage != nil:
		return Babbage
	default:
		return Era{}
	}
}

func (r RollForwardBlock) AlonzoOrGreaterBlock() *Block {
	if !r.Era().AlonzoOrGreater() {
		return nil
	}

	if r.Alonzo != nil {
		return r.Alonzo
	}

	if r.Babbage != nil {
		return r.Babbage
	}

	return nil
}
