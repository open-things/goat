package goat_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	check "gopkg.in/check.v1"

	"github.com/open-things/goat"
)

const (
	base string = "" +
		`{{define "main"}}<!DOCTYPE html>
<html>
	<head>
		<title>{{template "title" .}}</title>
		{{template "headers" .}}
	</head>
	<body>
		{{template "content" .}}
	</body>
</html>
{{end}}`
	adminBase string = "" +
		`{{define "title"}}ADMIN: {{value "title"}}{{end}}
{{define "headers"}}
	<script src="{{value "adminScriptUrl"}}"></script>
	{{template "custom_headers" .}}
{{end}}
{{define "content"}}
	<p>{{value "adminNotice"}}</p>
	{{template "custom_content" .}}
{{end}}`
	adminUsers string = "" +
		`{{define "custom_headers"}}{{end}}
{{define "custom_content"}}
	<!-- USERS -->
	<table>
		<thead><tr><th>username</th><th>last seen</th><th>actions</th></tr></thead>
		<tbody>{{range value "users"}}
			<tr><td>{{.Name}}</td><td>{{.LastSeen}}</td><td>{{template "admin_user_actions" .}}</td></tr>
		{{end}}</tbody>
	</table>
	<!-- /USERS -->
{{end}}
{{define "admin_users_main"}}{{template "main" .}}{{end}}`
	adminUsersActions string = "" +
		`{{define "admin_user_actions"}}
	<a href="{{actionUrl .Id "contact"}}"><img src="{{value "icon_contact"}}" title="Contact user '{{.Name}}'" /></a>
	<a href="{{actionUrl .Id "ban"}}"><img src="{{value "icon_ban"}}" title="Ban user '{{.Name}}'" /></a>
{{end}}`
	adminInfo string = "" +
		`{{define "custom_headers"}}{{end}}
{{define "custom_content"}}
	<table>{{range value "adminInfoStats"}}
		<tr><th>{{.Title}}</th><td>{{.Value}}</td></tr>
	{{end}}</table>
{{end}}
{{define "admin_info_main"}}{{template "main" .}}{{end}}`
)

type GoatFlowSuite struct {
	root   string
	loader *goat.Loader
}

var _ = check.Suite(&GoatFlowSuite{})

func (this *GoatFlowSuite) createFile(c *check.C, name, content string) {
	path := filepath.Join(this.root, name)
	err := ioutil.WriteFile(path, []byte(content), 0777)
	c.Assert(err, check.IsNil)
}

func (this *GoatFlowSuite) createTemplateFiles(c *check.C) {
	this.root = c.MkDir()
	c.Logf("Created temporary directory '%s' for templates", this.root)

	this.createFile(c, "base.html", base)
	this.createFile(c, "admin_base.html", adminBase)
	this.createFile(c, "admin_users.html", adminUsers)
	this.createFile(c, "admin_users_actions.html", adminUsersActions)
	this.createFile(c, "admin_info.html", adminInfo)
}

func (this *GoatFlowSuite) render(c *check.C, t *goat.Template) string {
	buf := bytes.Buffer{}
	err := t.Execute(&buf, nil)
	c.Assert(err, check.IsNil)
	return buf.String()
}

func (this *GoatFlowSuite) SetUpSuite(c *check.C) {
	this.loader = goat.NewLoader("")
	this.loader.SetExecNamePattern(":group:_:name:_main")
	// Global values
	this.loader.SetValue("title", "Goat Tests")
	this.loader.SetValue("adminScriptUrl", "/scripts/admin.js")
	this.loader.SetValue("adminNotice", "Notice: Don't ban random people!")
}

func (this *GoatFlowSuite) SetUpTest(c *check.C) {
	this.createTemplateFiles(c)
	this.loader.SetTemplateRoot(this.root)
}

func (this *GoatFlowSuite) TestAdminInfoCanBeExecuted(c *check.C) {
	template, err := this.loader.Get("admin", "info")
	c.Assert(err, check.IsNil)
	// overload global value
	template.SetValue("adminNotice", "Notice: Don't randomly restart servers!")
	// page-only value
	type adminStat struct{ Title, Value string }
	dummyStats := []adminStat{adminStat{Title: "Uptime", Value: "less than a day"}, adminStat{Title: "Free HDD space", Value: "7.8 GB"}}
	template.SetValue("adminInfoStats", dummyStats)

	// let's see what we get
	rendered := this.render(c, template)
	c.Logf("========  RENDERED  ========\n%s======== /RENDERED  ========\n", rendered)

	c.Assert(strings.Contains(rendered, "<title>ADMIN: Goat Tests</title>"), check.Equals, true)
	c.Assert(strings.Contains(rendered, "<p>Notice: Don&#39;t randomly restart servers!</p>"), check.Equals, true)
	c.Assert(strings.Contains(rendered, "<tr><th>Uptime</th><td>less than a day</td></tr>"), check.Equals, true)
	c.Assert(strings.Contains(rendered, "<tr><th>Free HDD space</th><td>7.8 GB</td></tr>"), check.Equals, true)
}

func (this *GoatFlowSuite) TestAdminUsersCanBeExecuted(c *check.C) {
	template, err := this.loader.Get("admin", "users")
	c.Assert(err, check.IsNil)
	// page-only values
	type adminUser struct {
		Id             int
		Name, LastSeen string
	}
	dummyUsers := []adminUser{
		adminUser{Id: 1, Name: "John Doe", LastSeen: "less than an hour ago"},
		adminUser{Id: 2, Name: "Jane Doe", LastSeen: "less than two hours ago"},
		adminUser{Id: 3, Name: "Jack Smith", LastSeen: "three weeks ago"},
	}
	template.SetValue("users", dummyUsers)
	template.SetValue("icon_contact", "/img/contact.png")
	template.SetValue("icon_ban", "/img/ban.png")
	// page-only func
	funcs := map[string]interface{}{
		"actionUrl": func(userId int, action string) string { return fmt.Sprintf("/action/%v?user=%v", action, userId) },
	}
	template.Funcs(funcs)

	// let's see what we get
	rendered := this.render(c, template)
	c.Logf("========  RENDERED  ========\n%s======== /RENDERED  ========\n", rendered)

	c.Assert(strings.Contains(rendered, "<title>ADMIN: Goat Tests</title>"), check.Equals, true)
	c.Assert(strings.Contains(rendered, "<p>Notice: Don&#39;t ban random people!</p>"), check.Equals, true)
	c.Assert(strings.Contains(rendered, "<td>John Doe</td><td>less than an hour ago</td>"), check.Equals, true)
	c.Assert(strings.Contains(rendered, `<a href="/action/contact?user=2"><img src="/img/contact.png" title="Contact user 'Jane Doe'" /></a>`), check.Equals, true)
}
