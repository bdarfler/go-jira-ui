package main

import (
	ui "github.com/gizak/termui"
	"github.com/op/go-logging"
	"os"
)

const (
	ticketQuery = iota
	ticketList  = iota
	labelList   = iota
	ticketShow  = iota
)

const (
	default_list_template = `{{ range .issues }}{{ .key | printf "%-20s"}}  {{ dateFormat "2006-01-02" .fields.created }}/{{ dateFormat "2006-01-02T15:04" .fields.updated }}  {{ .fields.summary | printf "%-75s"}} -- labels({{ join "," .fields.labels }})
{{ end }}`
)

var exitNow = false

type GoBacker interface {
	GoBack()
}

type ItemSelecter interface {
	SelectItem()
}

type TicketEditer interface {
	EditTicket()
}

type PagePager interface {
	NextLine(int)
	PreviousLine(int)
	NextPage()
	PreviousPage()
	Update()
}

type Navigable interface {
	Create(...interface{})
	Update()
	PreviousLine(int)
	NextLine(int)
	PreviousPage()
	NextPage()
}

var currentPage Navigable
var previousPage Navigable

var ticketQueryPage QueryPage
var ticketListPage TicketListPage
var labelListPage LabelListPage
var ticketShowPage TicketShowPage

func changePage(opts ...interface{}) {
	newopts := make(map[string]string)
	if len(opts) > 0 {
		newopts = opts[0].(map[string]string)
	}
	switch currentPage.(type) {
	case *QueryPage:
		currentPage.Create(newopts)
	case *TicketListPage:
		currentPage.Create(newopts)
	case *LabelListPage:
		currentPage.Create(newopts)
	case *TicketShowPage:
		currentPage.Create(newopts)
	}
}

var (
	log    = logging.MustGetLogger("jira")
	format = "%{color}%{time:2006-01-02T15:04:05.000Z07:00} %{level:-5s} [%{shortfile}]%{color:reset} %{message}"
)

func main() {

	var err error
	logging.SetLevel(logging.NOTICE, "")

	err = ensureLoggedIntoJira()
	if err != nil {
		log.Error("Login failed. Aborting")
		os.Exit(2)
	}

	err = ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()

	registerKeyboardHandlers()

	currentPage = &ticketQueryPage

	for exitNow != true {

		currentPage.Create()
		ui.Loop()

	}

}