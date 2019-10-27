package view

import "html/template"

var ListTmpl = template.Must(template.New("list").Parse(
	`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8">
			<title>宝くじ</title>
		</head>
		<body>
			<ul>
			{{range .}}
				<li>{{.Name}}
					<a href="/purchase_page?id={{.ID}}">Purchase</a>
					<a href="/result?id={{.ID}}">Result</a>
				</li>
			{{end}}
			</ul>
		</body>
	</html>
	`),
)
