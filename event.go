//TODO: Вынести в отдельный package controller
package main

type eventType string

const (
	addEcho             eventType = "addEcho"
)

type event struct {
	eventType      eventType
	oldObj, newObj interface{}
}