package checklist

// Checklist is a module for creating generic checklist implementations
// See 'Todo' for an implementation example
type Checklist struct {
	Items []*ChecklistItem

	checkedIcon   string
	selected      int
	uncheckedIcon string
}

func NewChecklist(checkedIcon, uncheckedIcon string) Checklist {
	list := Checklist{
		checkedIcon:   checkedIcon,
		selected:      -1,
		uncheckedIcon: uncheckedIcon,
	}

	return list
}

/* -------------------- Exported Functions -------------------- */

// Add creates a new item in the checklist
func (list *Checklist) Add(checked bool, text string) {
	item := NewChecklistItem(
		checked,
		text,
		list.checkedIcon,
		list.uncheckedIcon,
	)

	list.Items = append([]*ChecklistItem{item}, list.Items...)
}

// CheckedItems returns a slice of all the checked items
func (list *Checklist) CheckedItems() []*ChecklistItem {
	items := []*ChecklistItem{}

	for _, item := range list.Items {
		if item.Checked {
			items = append(items, item)
		}
	}

	return items
}

// Delete removes the selected item from the checklist
func (list *Checklist) Delete(selectedIndex int) {
	list.Items = append(list.Items[:selectedIndex], list.Items[selectedIndex+1:]...)
}

// IsSelectable returns true if the checklist has selectable items, false if it does not
func (list *Checklist) IsSelectable() bool {
	return list.selected >= 0 && list.selected < len(list.Items)
}

// IsUnselectable returns true if the checklist has no selectable items, false if it does
func (list *Checklist) IsUnselectable() bool {
	return !list.IsSelectable()
}

// LongestLine returns the length of the longest checklist item's text
func (list *Checklist) LongestLine() int {
	maxLen := 0

	for _, item := range list.Items {
		if len(item.Text) > maxLen {
			maxLen = len(item.Text)
		}
	}

	return maxLen
}

// IndexByItem returns the index of a giving item if found ,otherwise returns 0 with ok set to false
func (list *Checklist) IndexByItem(selectableItem *ChecklistItem) (index int, ok bool) {
	for idx, item := range list.Items {
		if item == selectableItem {
			return idx, true
		}
	}
	return 0, false
}

// UncheckedItems returns a slice of all the unchecked items
func (list *Checklist) UncheckedItems() []*ChecklistItem {
	items := []*ChecklistItem{}

	for _, item := range list.Items {
		if !item.Checked {
			items = append(items, item)
		}
	}

	return items
}

// Unselect removes the current select such that no item is selected
func (list *Checklist) Unselect() {
	list.selected = -1
}

/* -------------------- Sort Interface -------------------- */

func (list *Checklist) Len() int {
	return len(list.Items)
}

func (list *Checklist) Less(i, j int) bool {
	return list.Items[i].Text < list.Items[j].Text
}

func (list *Checklist) Swap(i, j int) {
	list.Items[i], list.Items[j] = list.Items[j], list.Items[i]
}
