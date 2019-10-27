package view

import "html/template"

var PurchasePageTmpl = template.Must(template.New("purchase_page").Parse(
	`
	<!DOCTYPE html>
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
	</html>
	`),
)
