package fproto_doc_html_default

import (
	"fmt"
	"io"
	"strings"

	"github.com/RangelReale/fproto"
	"github.com/RangelReale/fproto/fdep"
	"github.com/gosimple/slug"
)

type LayoutState int

const (
	LS_BEGIN LayoutState = iota
	LS_END
)

type Layout struct {
	w   io.Writer
	err error
}

func (l *Layout) Err() error {
	return l.err
}

func (l *Layout) WriteHeader() {
	if l.err != nil {
		return
	}

	_, l.err = fmt.Fprint(l.w, layout_header)
}

func (l *Layout) WriteFooter() {
	if l.err != nil {
		return
	}

	_, l.err = fmt.Fprint(l.w, layout_footer)
}

func (l *Layout) WriteContent(layoutState LayoutState) {
	if l.err != nil {
		return
	}

	switch layoutState {
	case LS_BEGIN:
		_, l.err = fmt.Fprint(l.w, content_begin)
	case LS_END:
		_, l.err = fmt.Fprint(l.w, content_end)
	}

}

func (l *Layout) WriteContentItem(layoutState LayoutState, itemName string, link string) {
	if l.err != nil {
		return
	}

	switch layoutState {
	case LS_BEGIN:
		_, l.err = fmt.Fprintf(l.w, `
        <div class="item">
            <a name="%s">%s</a>
        </div>
		`, link, itemName)
	case LS_END:
	}
}

func (l *Layout) WriteContentNs(layoutState LayoutState, nsName string, link string) {
	if l.err != nil {
		return
	}

	switch layoutState {
	case LS_BEGIN:
		_, l.err = fmt.Fprintf(l.w, `
        <div class="ns">
            <a name="%s">%s</a>
        </div>
		`, link, nsName)
	case LS_END:
	}
}

func (l *Layout) WriteContentNsItem(layoutState LayoutState, nsName string, link string, fileName string) {
	if l.err != nil {
		return
	}

	switch layoutState {
	case LS_BEGIN:
		_, l.err = fmt.Fprintf(l.w, `
        <div class="ns-item">
            <a name="%s">%s</a>
		`, link, nsName)

		if fileName != "" {
			fmt.Fprintf(l.w, `<span class="filename">[%s]</span>`, fileName)
		}

		_, l.err = fmt.Fprint(l.w, `
        </div>`)

	case LS_END:
	}
}

func (l *Layout) WriteNav(layoutState LayoutState) {
	if l.err != nil {
		return
	}

	switch layoutState {
	case LS_BEGIN:
		_, l.err = fmt.Fprintf(l.w, nav_begin)
	case LS_END:
		_, l.err = fmt.Fprintf(l.w, nav_end)
	}
}

func (l *Layout) WriteNavItem(layoutState LayoutState, itemName string, link string) {
	if l.err != nil {
		return
	}

	switch layoutState {
	case LS_BEGIN:
		_, l.err = fmt.Fprintf(l.w, `
        <div class="item">
            <a href="#%s">%s</a>
        </div>
		`, link, itemName)
	case LS_END:
	}
}

func (l *Layout) WriteNavNs(layoutState LayoutState, nsName string, link string) {
	if l.err != nil {
		return
	}

	switch layoutState {
	case LS_BEGIN:
		_, l.err = fmt.Fprintf(l.w, `
        <div class="ns">
            <a href="#%s">%s</a>
        </div>
		`, link, nsName)
	case LS_END:
	}
}

func (l *Layout) WriteNavNsItem(layoutState LayoutState, nsName string, link string) {
	if l.err != nil {
		return
	}

	switch layoutState {
	case LS_BEGIN:
		_, l.err = fmt.Fprintf(l.w, `
        <div class="ns-item">
            <a href="#%s">%s</a>
        </div>
		`, link, nsName)
	case LS_END:
	}
}

//
// Data
//

func (l *Layout) WriteContentService(dt *fdep.DepType) {
	if l.err != nil {
		return
	}

	element := dt.Item.(*fproto.ServiceElement)

	svc_comment := l.concatComment(element.Comment)

	fmt.Fprint(l.w, `<div class="definition service">`)
	if svc_comment != "" {
		fmt.Fprint(l.w, `<div class="description"><p>`)
		fmt.Fprintf(l.w, `%s`, svc_comment)
		fmt.Fprint(l.w, `</p></div>`)
	}

	fmt.Fprint(l.w, `<div class="list">
		<table>
			<tr>
				<th>Method name</th><th>Request Type</th><th>Response Type</th><th>Description</th>
			</tr>`)

	for _, rpc := range element.RPCs {
		rpc_comment := l.concatComment(rpc.Comment)

		fmt.Fprintf(l.w, `
		<tr>
			<td class="fld-svc-method">%s</td>
			<td class="fld-svc-req"><a href="#%s">%s</a></td>
			<td  class="fld-svc-ret"><a href="#%s">%s</a></td>
			<td class="fld-svc-doc">%s</td>
		</tr>`,
			rpc.Name, "cc", rpc.RequestType, "cc", rpc.ResponseType, rpc_comment)
	}

	_, l.err = fmt.Fprint(l.w, `</table>
	</div>
	</div>`)
}

func (l *Layout) WriteContentMessage(dt *fdep.DepType) {
	if l.err != nil {
		return
	}

	element := dt.Item.(*fproto.MessageElement)

	msg_comment := l.concatComment(element.Comment)

	fmt.Fprint(l.w, `<div class="definition message">`)
	if msg_comment != "" {
		fmt.Fprint(l.w, `<div class="description"><p>`)
		fmt.Fprintf(l.w, `%s`, msg_comment)
		fmt.Fprint(l.w, `</p></div>`)
	}

	l.writeFields(dt, element.Fields, "")

	_, l.err = fmt.Fprint(l.w, `</div>`)
}

func (l *Layout) WriteContentOneofFields(dt *fdep.DepType, fields []fproto.FieldElementTag) {
	if l.err != nil {
		return
	}

	for _, fld := range fields {
		switch xfld := fld.(type) {
		case *fproto.OneofFieldElement:
			fmt.Fprintf(l.w, `<div class="ns-itemsub">
				<a name="%s">Oneof %s.%s</a>
			</div>`, fmt.Sprintf("content-Oneof-%s-%s", slug.Make(dt.FullOriginalName()), slug.Make(xfld.Name)), dt.Name, xfld.Name)

			fmt.Fprint(l.w, `<div class="definition message">`)

			l.writeFields(dt, xfld.Fields, "oneof")

			fmt.Fprint(l.w, `</div>`)
		}
	}
}

func (l *Layout) writeFields(dt *fdep.DepType, fields []fproto.FieldElementTag, tableClass string) {
	if l.err != nil {
		return
	}

	if tableClass != "" {
		tableClass = fmt.Sprintf(" class=\"%s\"", tableClass)
	}

	fmt.Fprintf(l.w, `<div class="list">
		<table%s>
			<tr>
				<th>Fieldname</th><th>Type</th><th>Flags</th><th>Description</th>
			</tr>`, tableClass)

	for _, fld := range fields {
		var fld_comment string
		var fld_type string
		var fld_type_link string
		var fld_type_extra string
		var fld_opt []string

		switch xfld := fld.(type) {
		case *fproto.FieldElement:
			fld_comment = l.concatComment(xfld.Comment)

			// load field type
			ft, err := dt.GetType(xfld.Type)
			if err != nil {
				l.err = err
				return
			}

			if ft != nil {
				fld_type = ft.FullOriginalName()
				if !ft.IsScalar() && ft.FileDep.DepType == fdep.DepType_Own {
					switch ft.Item.(type) {
					case *fproto.EnumElement:
						fld_type_link = fmt.Sprintf("content-Enum-%s", slug.Make(fld_type))
					default:
						fld_type_link = fmt.Sprintf("content-Message-%s", slug.Make(fld_type))
					}
				}
			} else {
				fld_type = xfld.Type
			}

			if xfld.Required {
				fld_opt = append(fld_opt, "required")
			}
			if xfld.Repeated {
				fld_opt = append(fld_opt, "repeated")
			}
			if xfld.Optional {
				fld_opt = append(fld_opt, "optional")
			}
		case *fproto.MapFieldElement:
			fld_comment = l.concatComment(xfld.Comment)

			// load key and field type
			ftkey, err := dt.GetType(xfld.KeyType)
			if err != nil {
				l.err = err
				return
			}

			ft, err := dt.GetType(xfld.Type)
			if err != nil {
				l.err = err
				return
			}

			var f_key string
			var f_value string
			if ftkey != nil {
				f_key = ftkey.FullOriginalName()
			} else {
				f_key = xfld.KeyType
			}
			if ft != nil {
				f_value = ft.FullOriginalName()
			} else {
				f_value = xfld.Type
			}

			fld_type = fmt.Sprintf("map[%s]%s", f_key, f_value)
		case *fproto.OneofFieldElement:
			fld_type = fmt.Sprintf("oneof %s.%s", dt.Name, xfld.Name)
			fld_type_link = fmt.Sprintf("content-Oneof-%s-%s", slug.Make(dt.FullOriginalName()), slug.Make(xfld.Name))

			var fextra []string
			for _, oofld := range xfld.Fields {
				fextra = append(fextra, oofld.FieldName())
			}

			fld_type_extra = strings.Join(fextra, ", ")
		}

		ftlink := fld_type
		if fld_type_link != "" {
			ftlink = fmt.Sprintf(`<a href="#%s">%s</a>`, fld_type_link, fld_type)
		}

		ftextra := ""
		if fld_type_extra != "" {
			ftextra = fmt.Sprintf("<br/>(%s)", fld_type_extra)
		}

		fmt.Fprintf(l.w, `
			<tr>
				<td class="fld-msg-fieldname">%s</td>
				<td class="fld-msg-type">%s%s</td>
				<td  class="fld-msg-opt">%s</td>
				<td class="fld-msg-doc">%s</td>
			</tr>`,
			fld.FieldName(), ftlink, ftextra, strings.Join(fld_opt, ","), fld_comment)
	}

	_, l.err = fmt.Fprint(l.w, `</table>
	</div>`)
}

func (l *Layout) concatComment(comment *fproto.Comment) string {
	var ret string

	if comment != nil && len(comment.Lines) > 0 {
		for _, cl := range comment.Lines {
			ln := strings.TrimSpace(cl)
			if len(ln) > 0 {
				if ret != "" {
					ret += " "
				}
				ret += ln
			}
		}
	}
	return ret
}

// layout strings
var (
	layout_header = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
    <style type="text/css">
        /* RESET BEGIN */
        html, body, div, span, applet, object, iframe,
        h1, h2, h3, h4, h5, h6, p, blockquote, pre,
        a, abbr, acronym, address, big, cite, code,
        del, dfn, em, img, ins, kbd, q, s, samp,
        small, strike, strong, sub, sup, tt, var,
        b, u, i, center,
        dl, dt, dd, ol, ul, li,
        fieldset, form, label, legend,
        table, caption, tbody, tfoot, thead, tr, th, td,
        article, aside, canvas, details, embed,
        figure, figcaption, footer, header, hgroup,
        menu, nav, output, ruby, section, summary,
        time, mark, audio, video {
            margin: 0;
            padding: 0;
            border: 0;
            font-size: 100%;
            font: inherit;
            vertical-align: baseline;
        }
        /* HTML5 display-role reset for older browsers */
        article, aside, details, figcaption, figure,
        footer, header, hgroup, menu, nav, section {
            display: block;
        }
        body {
            line-height: 1;
        }
        ol, ul {
            list-style: none;
        }
        blockquote, q {
            quotes: none;
        }
        blockquote:before, blockquote:after,
        q:before, q:after {
            content: '';
            content: none;
        }
        table {
            border-collapse: collapse;
            border-spacing: 0;
        }
        /* RESET END*/

        body
        {
            font-family: "Lucida Sans", "Lucida Grande", Verdana, Arial, sans-serif;
            font-size: 13px;
            width: 100%;
            margin: 0;
            padding: 0;
            background-color: #FFFFFF;

            display: flex;
            flex-direction: column;
            min-height: 100vh;
        }

        a, a:visited {
            text-decoration: none;
            color: #05a;
        }

        a[name] {
            color: black;
        }

        .header{
            width: 100%;
            height: 60px;
        }

        .body {
            flex: 1 0 auto;
            display: flex;
        }

        .body .content{
            flex: 1 0 auto;
            line-height: 1.5145em;
            font-size: 15px;

            padding: 1.2em;
            padding-top: 0.2em;
        }

        .body .content .content-header {
            margin-bottom: 10px;
        }

        .body .content .item{
            font-weight: bold;
            border-bottom: solid 1px black;
            font-size: 1.2em;
            margin: 6px 0;
        }

        .body .content .ns{
            font-style: italic;
            background-color: #fafafa;
			margin-top: 6px;
        }

        .body .content .ns-item{
            margin-left: 20px;
            font-size: 1.1em;
			padding: 8px 0;
        }

        .body .content .ns-itemsub{
            margin-left: 30px;
            font-size: 1.0em;
			padding: 8px 0;
        }

        .body .content .ns-item .filename{
            color: #a0a0a0;
			font-size: 0.8em;
			margin-left: 10px;
        }

        .body .nav{
            width: 20%;
            text-align: left;
            order: -1;
            margin: 0;
            font-size: 1.1em;
            background: #fdfdfd;
        }

        .body .nav .menu{
        }

        .body .nav .menu div {
            padding: 5px 5px 5px 0;
            background: #fafafa;
        }

        .body .nav .menu div:nth-child(odd){
            background-color: white;
        }

        .body .nav .menu .item{
            padding-left: 10px;
            font-weight: bold;
        }

        .body .nav .menu .ns{
            padding-left: 30px;
            font-style: italic;
        }

        .body .nav .menu .ns-item{
            padding-left: 50px;
        }

        .body .nav .nav-header {
            background: #ffffff;
            padding: 1em;
        }

        .footer{
            width: 100%;
            height: 60px;
        }

        @media (max-width: 700px) {
            .body {
                flex-direction: column;
            }

            .body .nav{
                width: 100%;
            }
        }

        h1 {
            font-size: 1.4em;
        }

        .body .content .definition {
            margin-left: 40px;
            /*background-color: #fcfcfc;*/
        }

        .body .content .definition .description {
            font-family: "Helvetica Neue-Light", "Helvetica Neue Light", "Helvetica Neue", Helvetica, Arial, "Lucida Grande", sans-serif;
            background-color: white;
            padding: 0px 8px 4px;
            font-size: 1em;
        }

        .body .content .definition .list table {
            width: 75%;
        }

        .body .content .definition .list table th {
            background-color: #e0e0e0;
        }

        .body .content .definition .list table td, .body .content .definition .list table th {
            border: solid 1px black;
            font-size: 0.9em;
            padding: 1px 4px;
        }

        .body .content .definition .list table tbody tr:nth-child(odd){
            background-color: #fcfcfc;
        }

        .body .content .definition .list td.fld-svc-method {
            width: 15%;
        }

        .body .content .definition .list td.fld-svc-req {
            width: 20%;
        }

        .body .content .definition .list td.fld-svc-ret {
            width: 20%;
        }

        .body .content .definition .list td.fld-msg-fieldname {
            width: 15%;
        }

        .body .content .definition .list td.fld-msg-type {
            width: 30%;
        }

        .body .content .definition .list td.fld-msg-opt {
            width: 20%;
        }

        .body .content .definition .list table.oneof{
            margin-top: 6px;
        }

        .body .content .definition .list table.oneof th.table-title {
            font-weight: bold;
        }


    </style>
</head>
<body>

<header class="header">Header</header>
<div class="body">
`

	layout_footer = `
</div>
<footer class="footer">Footer</footer>

</body>
</html>
`

	content_begin = `
    <div class="content">
        <h1 class="content-header">Documentation</h1>
`

	content_end = `
	</div>
`

	nav_begin = `
    <div class="nav">
        <div class="nav-header">
            <h1>Table of Contents</h1>
        </div>
		<div class="menu">

`

	nav_end = `
		</div>
	</div>
`
)
