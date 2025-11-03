package app

import (
	"fmt"

	"github.com/jefrryss/Bank/internal/handmaker"
)

func (m *app) prepareForm() {
	switch m.manualEntity {
	case "account":
		if m.manualAction == "create" {
			m.formFields = handmaker.AccountCreateFields()
		} else {
			m.formFields = handmaker.AccountUpdateFields()
		}
	case "category":
		if m.manualAction == "create" {
			m.formFields = handmaker.CategoryCreateFields()
		} else {
			m.formFields = handmaker.CategoryUpdateFields()
		}
	case "operation":
		if m.manualAction == "create" {
			m.formFields = handmaker.OperationCreateFields()
		} else {
			m.formFields = handmaker.OperationUpdateFields()
		}
	}
	m.formValues = make(map[string]string, len(m.formFields))
	m.formStep = 0
	m.textInput.Placeholder = m.formFields[0].Label
	m.textInput.SetValue("")
}

func (m *app) submitForm() error {
	switch m.manualEntity {
	case "account":
		if m.manualAction == "create" {
			return m.ctrl.CreateAccount(
				m.formValues["id"],
				m.formValues["name"],
				m.formValues["balance"],
			)
		}
		return m.ctrl.UpdateAccount(
			m.formValues["id"],
			m.formValues["name"],
			m.formValues["balance"],
		)

	case "category":
		if m.manualAction == "create" {
			return m.ctrl.CreateCategory(
				m.formValues["id"],
				m.formValues["name"],
				m.formValues["type"],
			)
		}
		return m.ctrl.UpdateCategory(
			m.formValues["id"],
			m.formValues["name"],
			m.formValues["type"],
		)

	case "operation":
		if m.manualAction == "create" {
			return m.ctrl.CreateOperation(
				m.formValues["id"],
				m.formValues["type"],
				m.formValues["account_id"],
				m.formValues["category_id"],
				m.formValues["amount"],
				m.formValues["date"],
				m.formValues["description"],
			)
		}
		return m.ctrl.UpdateOperation(
			m.formValues["id"],
			m.formValues["type"],
			m.formValues["account_id"],
			m.formValues["category_id"],
			m.formValues["amount"],
			m.formValues["date"],
			m.formValues["description"],
		)
	}
	return fmt.Errorf("неизвестная сущность")
}
