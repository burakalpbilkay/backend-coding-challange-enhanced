package constants

type ActionType string

const (
	ActionWelcome      ActionType = "WELCOME"
	ActionConnectCRM   ActionType = "CONNECT_CRM"
	ActionEditContract ActionType = "EDIT_CONTACT"
	ActionAddContract  ActionType = "ADD_CONTACT"
	ActionViewContract ActionType = "VIEW_CONTACTS"
	ActionReferUser    ActionType = "REFER_USER"
)

// Set of valid action types for validation
var ValidActionTypes = map[ActionType]bool{
	ActionWelcome:      true,
	ActionConnectCRM:   true,
	ActionEditContract: true,
	ActionAddContract:  true,
	ActionViewContract: true,
	ActionReferUser:    true,
}
