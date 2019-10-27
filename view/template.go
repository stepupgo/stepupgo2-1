package view

import "html/template"

var ListTmpl = template.Must(template.New("list").Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>宝くじAPP</title>
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
</html>`))

var PurchasePageTmpl = template.Must(template.New("purchase_page").Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Purchase {{.Name}}</title>
	</head>
	<body>
		<h2>{{.Name}}</h2>
		<h3>Prizes</h3>
		<ul>
		{{range .Prizes}}
			<li>{{.Name}} ${{.Amount}}</li>
		{{end}}
		</ul>

		{{if .Remain}}
		<form action="/purchase" method="post">
			<input name="id" type="hidden" value="{{.ID}}">
			<label for="amount">Amount({{.Remain}} / {{.Num}}): </label>
			<input name="num" type="number" min="1" max="{{.Remain}}">
			<input type="submit" value="Purchase">
		</form>
		{{else}}
			SOLD OUT
		{{end}}
	</body>
</html>`))

var ResultTmpl = template.Must(template.New("result").Parse(`<!DOCTYPE html>
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
</html>`))
