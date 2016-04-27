package main

type valueString struct {
    v string
}

func (vs *valueString) getType() valueType {
    return vtString
}

func(vs *valueString) result() []byte {
    return []byte(vs.v)
}