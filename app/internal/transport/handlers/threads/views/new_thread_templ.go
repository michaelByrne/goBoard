// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.663
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import commonviews "goBoard/internal/transport/handlers/common/views"
import "fmt"

func NewThreadTitleGroup() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<title>Create New Thread</title><h3>Create New Thread</h3>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func NewThreadForm(username string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"hr\"><hr></div><span class=\"response_form\" hx-ext=\"response-targets\"><div class=\"clear\"></div><div id=\"response_form\"></div><form name=\"form\" id=\"form\" class=\"coreform\" hx-boost=\"true\" hx-post=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/thread/create"))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `app/internal/transport/handlers/threads/views/new_thread.templ`, Line: 23, Col: 42}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" hx-target-400=\"find .error_hanger\"><input type=\"hidden\" id=\"idx\" name=\"idx\" value=\"1\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = commonviews.Account(username).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<fieldset><legend>Thread Details</legend><ol><li><label id=\"label_subject\" for=\"subject\">Subject:</label> <input type=\"text\" name=\"subject\" id=\"subject\" value=\"\" style=\"width:400px;\" maxlength=\"200\" notnull=\"Please enter a subject.\" class=\"validate_form\"><div class=\"clear\"></div></li><li><label id=\"label_body\" for=\"body\">Body: </label><div contentEditable=\"true\" id=\"body\" name=\"body\" style=\"float: left; height: 100px; width: 500px;\" _=\"on htmx:afterSwap or input put my innerHTML into #hidden_body&#39;s value\"></div><textarea id=\"hidden_body\" hidden name=\"hidden_body\"></textarea><div class=\"clear\"></div></li></ol></fieldset><input type=\"submit\" class=\"submit\"> <input type=\"button\" name=\"preview\" id=\"preview\" value=\"preview\" hx-post=\"/thread/preview\" hx-target=\"#response_form\" hx-swap=\"afterbegin\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = commonviews.Uploader().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<sup><a style=\"cursor:pointer;color:white\" _=\"on click toggle @hidden on #bbcode\">[help]</a></sup> <span class=\"error_hanger\"></span></form><div id=\"bbcode\" class=\"view\" style=\"font-size: 0.85em\" hidden><pre><h4>TAGS:</h4>http://www.google.com/ &lt;-- automatic link<br>[url]http://www.google.com/[/url]<br>[url=http://www.google.com/]with my own link text[/url]<br>[img]http://www.google.com/intl/en_ALL/images/logo.gif[/img]<br>[u]underline[/u]<br>[strong]bold[/strong]<br>[b]bold[/b]<br>[i]italic[/i]<br>[em]italic[/em]<br>[strike]strikethrough[/strike]<br>[code]like pre[/code]<br>[sub]subscript[/sub]<br>[sup]superscript[/sup]<br>[soundcloud]http://soundcloud.com/goingslowly/047-railroad-lullabye[/soundcloud]<br>[youtube]http://youtube.com/watch?v=WAwLYJYsa0A[/youtube] or [youtube]http://youtu.be/L8xXb-P4wZY[/youtube]<br>[vimeo]http://vimeo.com/2467457[/vimeo]<br>[tweet]https://twitter.com/dril/status/134787490526658561[/tweet]<br>[quote]quote[/quote]<br>[spoiler]spoiler[/spoiler]<br>[trigger]trigger[/trigger]<br></pre><div class=\"clear\"></div></div></span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
