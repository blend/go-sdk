package jobkit

var headerTemplate = `
{{ define "header" }}
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="referrer" content="no-referrer"/>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1"/>
	<meta name="robots" content="noimageindex">
	<title>Jobkit</title>
	<link rel="stylesheet" href="/static/css/normalize.css"/>
	<link rel="stylesheet" href="/static/css/skeleton.css"/>
	<style type="text/css">
		table.u-full-width {
			font-size: 12px;
		}

		form {
			margin: 0;
			padding: 0;
		}

		.align-left {
			text-align: left;
		}
		.align-center {
			text-align: center;
		}
		.align-right {
			text-align: right;
		}
		.container {
			max-width: 1200px;
			border: 1px solid #efefef;
			padding: 10px;
		}
		h4 {
			font-size: inherit;
			font-weight: bold;
		}
		.small-text {
			font-size: 10px;
		}

		tr.running {
			color: white;
			background-color: #0F9960;
		}
		tr.running a {
			color: white;
		}
		tr.failed {
			background-color: #F55656
		}
		tr.cancelled {
			background-color: #FFB366;
		}

		ul.breadcrumbs {
			font-size: 12px;
			margin-top: 10px;
			list-style: none;
		}
		ul.breadcrumbs li {
			display: inline-block;
		}
		ul.breadcrumbs li:after {
			margin-left: 5px;
			content: ">";
		}
		ul.breadcrumbs li:last-child:after {
			margin-left: 0;
			content: none;
		}
	</style>
</head>
<body>
{{ end }}
`
