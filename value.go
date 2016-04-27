package main


type valueType int
const (
    _ valueType = iota
    vtString
    vtMap
    vtList
    vtSet 
)

type value interface {
    getType() valueType
    result() []byte
}
