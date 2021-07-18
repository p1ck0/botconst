package models

type Dialog struct {
	IsMain  bool
	Trigger string
	Actions []Action
}
