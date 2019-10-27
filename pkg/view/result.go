package view

import "html/template"

var ResultTmpl = template.Must(template.New("result").Parse(
	`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="utf-8">
			<title>Result of {{.Name}}</title>
		</head>
		<body>
			{{range .Prizes}}
				{{$winner := index $.Winners .ID}}
				{{if $winner }}
					<h2>{{.Name}} (${{.Amount}})</h2>
					<ul>
						{{range $winner.Numbers }}
						<li>{{.}}</li>
						{{end}}
					</ul>
				{{else}}
					<h2>{{.Name}} (${{.Amount}})</h2>
					No winners
				{{end}}
			{{end}}
		</body>
	</html>
	`),
)
