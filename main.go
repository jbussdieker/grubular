package main

import "io/ioutil"
import "path/filepath"

import rubex "rubex/lib"
import "web.go"

func serveAjax(ctx *web.Context) string {
	defer func() {
		recover()
	}()

	test := ctx.Request.FormValue("test")
	regex := ctx.Request.FormValue("regex")
	message_id := ctx.Request.FormValue("message_id")
	ctx.ContentType(".js")


	result := "if (Rubular.lastMessageSent == " + message_id + ") {\n"
	result += "Rubular.recordLastMessageReceived(" + message_id + ");\n"
	result += "Element.update(\"result\", \"  "
	result += "<span class=\\\"result_label\\\">Match result:</span>\\n"
	//result += "<span class=\\"result_label\\">Match result:</span>\\n  <div id=\\"match_string\\" class=\\"\\"><span id=\\"match_string_inner\\"><span class=\\"match\\">asdfaa</span><span class=\\"match\\"></span></span></div>\\n\\n\");\n"

	x := rubex.MustCompile(regex)
	r := x.FindStringSubmatch(test)
	for _, match := range r {
		result += "<div id=\\\"match_string\\\">\\n"
		result += "<span id=\\\"match_string_inner\\\">"
		result += "<span class=\\\"match\\\">"
		result += match
		result += "</span>"
		result += "</span>"
		result += "</div>"
	}

	result += "\");\n"
	result += "Rubular.handleParseSuccess();\n"
	result += "}"

	return result
}

func serveAsset(ctx *web.Context, message string) string {
	println(message)
	if message == "" {
		// Read the file contents
		data, err := ioutil.ReadFile("template.html")
		if err != nil {
			ctx.Server.Logger.Warn("ERROR: Failed to read asset: " + err.String())
			ctx.NotFound(err.String())
		} else {
			ctx.ContentType(".html")
		}

		return string(data)
	} else {
		// Read the file contents
		data, err := ioutil.ReadFile(message)
		if err != nil {
			ctx.Server.Logger.Warn("ERROR: Failed to read asset: " + err.String())
			ctx.NotFound(err.String())
		} else {
			ctx.ContentType(filepath.Ext(message))
		}

		return string(data)
	}
	return ""
}

func main() {
	web.Get("/(.*)", serveAsset)
	web.Post("/regex/do_test", serveAjax)
	web.Run("0.0.0.0:6060")
}

