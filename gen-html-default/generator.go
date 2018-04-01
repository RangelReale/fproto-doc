package fproto_doc_html_default

import (
	"fmt"
	"io"

	"github.com/RangelReale/fdep"
	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto-doc"
	"github.com/gosimple/slug"
)

type Generator struct {
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Generate(dep *fdep.Dep, w io.Writer) error {
	layout := &Layout{w: w}

	helper := fproto_doc.NewHelper(dep)

	//
	// HEADER
	//
	layout.WriteHeader()

	type litem struct {
		layoutItem layoutItem
		list       []*fdep.DepType
	}

	llist := []*litem{
		{layoutItem: li_service, list: helper.GetServiceList(fproto_doc.NewGetFilter(fproto_doc.ST_ALIAS_NAME, fproto_doc.DT_OWN))},
		{layoutItem: li_enum, list: helper.GetEnumList(fproto_doc.NewGetFilter(fproto_doc.ST_ALIAS_NAME, fproto_doc.DT_OWN))},
		{layoutItem: li_message, list: helper.GetMessageList(fproto_doc.NewGetFilter(fproto_doc.ST_ALIAS_NAME, fproto_doc.DT_OWN))},
	}

	last_alias := ""
	slug_ns := ""

	//
	// NAV
	//
	layout.WriteNav(LS_BEGIN)

	for _, li := range llist {
		layout.WriteNavItem(LS_BEGIN, li.layoutItem.Title(), fmt.Sprintf("content-%s", li.layoutItem.String()))

		last_alias = ""
		slug_ns = ""
		for _, e := range li.list {
			if e.Alias != last_alias {
				if last_alias != "" {
					layout.WriteNavNs(LS_END, last_alias, "")
				}

				slug_ns = slug.Make(e.Alias)

				layout.WriteNavNs(LS_BEGIN, e.Alias, fmt.Sprintf("content-%s-%s", li.layoutItem.String(), slug_ns))
				last_alias = e.Alias
			}

			slug_nsitem := slug.Make(e.Name)

			layout.WriteNavNsItem(LS_BEGIN, e.Name, fmt.Sprintf("content-%s-%s-%s", li.layoutItem.String(), slug_ns, slug_nsitem))
			layout.WriteNavNsItem(LS_END, e.Name, "")
		}
		if last_alias != "" {
			layout.WriteNavNs(LS_END, last_alias, "")
		}

		layout.WriteNavItem(LS_END, li.layoutItem.String(), "")
	}

	layout.WriteNav(LS_END)

	//
	// CONTENT
	//
	layout.WriteContent(LS_BEGIN)

	for _, li := range llist {
		layout.WriteContentItem(LS_BEGIN, li.layoutItem.Title(), fmt.Sprintf("content-%s", li.layoutItem.String()))

		last_alias = ""
		slug_ns = ""
		for _, e := range li.list {
			if e.Alias != last_alias {
				if last_alias != "" {
					layout.WriteContentNs(LS_END, last_alias, "")
				}

				slug_ns = slug.Make(e.Alias)

				layout.WriteContentNs(LS_BEGIN, e.Alias, fmt.Sprintf("content-%s-%s", li.layoutItem.String(), slug_ns))
				last_alias = e.Alias
			}

			slug_nsitem := slug.Make(e.Name)
			fn := ""
			if e.DepFile != nil {
				fn = e.DepFile.FilePath
			}

			layout.WriteContentNsItem(LS_BEGIN, e.Name, fmt.Sprintf("content-%s-%s-%s", li.layoutItem.String(), slug_ns, slug_nsitem), fn, e.Alias)

			switch li.layoutItem {
			case li_service:
				layout.WriteContentService(e)
			case li_enum:
				layout.WriteContentEnum(e)
			case li_message:
				layout.WriteContentMessage(e)
				layout.WriteContentOneofFields(e, helper.GetOneOfFieldList(e.Item.(*fproto.MessageElement).Fields))
			}

			layout.WriteContentNsItem(LS_END, e.Name, "", "", "")
		}
		if last_alias != "" {
			layout.WriteContentNs(LS_END, last_alias, "")
		}

		layout.WriteContentItem(LS_END, li.layoutItem.Title(), "")
	}

	layout.WriteContent(LS_END)

	//
	// FOOTER
	//
	layout.WriteFooter()

	return layout.Err()
}

type layoutItem int

const (
	li_service layoutItem = iota
	li_enum
	li_message
)

func (li layoutItem) String() string {
	switch li {
	case li_service:
		return "Service"
	case li_enum:
		return "Enum"
	case li_message:
		return "Message"
	}
	return "Unknown"
}

func (li layoutItem) Title() string {
	switch li {
	case li_service:
		return "Services"
	case li_enum:
		return "Enums"
	case li_message:
		return "Messages"
	}
	return "Unknown"
}
