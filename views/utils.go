package views

import (
	"fmt"
	"todo/model"
)

func deleteUrl(todo model.Todo) string {
	return fmt.Sprintf("/delete?item=%s", todo.Item)
}

func updateUrl(todo model.Todo) string {
	return fmt.Sprintf("/update?item=%s", todo.Item)
}
