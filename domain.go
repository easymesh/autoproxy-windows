package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/astaxie/beego/logs"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

type DomainItem struct {
	Index  int
	Domain string

	checked bool
}

type DomainModel struct {
	sync.RWMutex

	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder

	items []*DomainItem
}

func (n *DomainModel) RowCount() int {
	return len(n.items)
}

func (n *DomainModel) Value(row, col int) interface{} {
	item := n.items[row]
	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Domain
	}
	panic("unexpected col")
}

func (n *DomainModel) Checked(row int) bool {
	return n.items[row].checked
}

func (n *DomainModel) SetChecked(row int, checked bool) error {
	n.items[row].checked = checked
	return nil
}

func (m *DomainModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]
		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}
		switch m.sortColumn {
		case 0:
			return c(a.Index < b.Index)
		case 1:
			return c(a.Domain < b.Domain)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

var domainTable *DomainModel

func DomainTableUpdate(find string) {
	item := make([]*DomainItem, 0)
	domainList := ConfigGet().DomainList
	for i, value := range domainList {
		if strings.Index(value, find) == -1 {
			continue
		}
		item = append(item, &DomainItem{
			Index: i, Domain: value,
		})
	}
	domainTable.items = item
	domainTable.PublishRowsReset()
	domainTable.Sort(domainTable.sortColumn, domainTable.sortOrder)
}

func DomainAdd(domain string) error {
	domainList := ConfigGet().DomainList
	for _, v := range domainList {
		if v == domain {
			return fmt.Errorf("Domain %s already exists", domain)
		}
	}
	domainList = append(domainList, domain)
	sort.Strings(domainList)
	return DomainListSave(domainList)
}

func DomainDelete(owner *walk.Dialog) error {
	var remainderList []string
	var deleteList []string
	for _, v := range domainTable.items {
		if !v.checked {
			remainderList = append(remainderList, v.Domain)
		} else {
			deleteList = append(deleteList, v.Domain)
		}
	}
	if len(deleteList) == 0 {
		return fmt.Errorf("No choice any domain")
	}

	err := DomainListSave(remainderList)
	if err != nil {
		ErrorBoxAction(owner, fmt.Sprintf("%v %s", deleteList, "Delete Fail"))
	} else {
		InfoBoxAction(owner, fmt.Sprintf("%v %s", deleteList, "Delete Success"))
	}
	return err
}

func RemodeEdit() {
	domainTable = new(DomainModel)
	domainTable.items = make([]*DomainItem, 0)

	DomainTableUpdate("")

	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var findPB, addPB *walk.PushButton
	var addLine, findLine *walk.LineEdit

	_, err := Dialog{
		AssignTo:      &dlg,
		Title:         "Forward Domain",
		Icon:          walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		Size:          Size{300, 450},
		MinSize:       Size{300, 300},
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 3, MarginsZero: true},
				Children: []Widget{
					Label{
						Text: "Domain: ",
					},
					LineEdit{
						AssignTo: &addLine,
						Text:     "",
					},
					PushButton{
						AssignTo: &addPB,
						Text:     "Add",
						OnClicked: func() {
							addDomain := addLine.Text()

							if addDomain == "" {
								ErrorBoxAction(dlg, "Input Domain")
								return
							}

							err := DomainAdd(addDomain)
							if err != nil {
								ErrorBoxAction(dlg, err.Error())
								return
							}

							addLine.SetText("")
							findLine.SetText("")
							DomainTableUpdate("")
							RouteUpdate()

							InfoBoxAction(dlg, addDomain+" add success")
							return
						},
					},
					Label{
						Text: "Find Key: ",
					},
					LineEdit{
						AssignTo: &findLine,
						Text:     "",
					},
					PushButton{
						AssignTo: &findPB,
						Text:     "Find",
						OnClicked: func() {
							DomainTableUpdate(findLine.Text())
						},
					},
				},
			},
			TableView{
				AlternatingRowBG: true,
				ColumnsOrderable: true,
				CheckBoxes:       true,
				Columns: []TableViewColumn{
					{Title: "#", Width: 60},
					{Title: "Domain", Width: 160},
				},
				StyleCell: func(style *walk.CellStyle) {
					if style.Row()%2 == 0 {
						style.BackgroundColor = walk.RGB(248, 248, 255)
					} else {
						style.BackgroundColor = walk.RGB(220, 220, 220)
					}
				},
				Model: domainTable,
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     "Delete",
						OnClicked: func() {
							err := DomainDelete(dlg)
							if err != nil {
								logs.Error(err.Error())
								ErrorBoxAction(dlg, err.Error())
							} else {
								DomainTableUpdate(findLine.Text())
								RouteUpdate()
							}
						},
					},
					PushButton{
						AssignTo: &cancelPB,
						Text:     "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		logs.Error(err.Error())
	}
}
