package jobkit

var headerTemplate = `
{{ define "header" }}
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.2.0/css/uikit.min.css" />
	<script src="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.2.0/js/uikit.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/uikit/3.2.0/js/uikit-icons.min.js"></script>
	<style>
	#top-head {
		z-index: 9;
		top: 0;
		left:0;
		right:0;
	}
	/* Smaller Header */
	.uk-navbar-nav > li > a,
	.uk-navbar-item,
	.uk-navbar-toggle {
		/* navbar height */
	min-height: 52px;
	padding: 0 8px;
	font-size: 0.85rem;
	}
	#content {
		margin-top: 52px;
		padding: 30px 0 0 0;
		background-color: #f7f7f7;
		margin-left: 0;
		transition: margin 0.2s cubic-bezier(.4,0,.2,1);
	}
	.uk-form {
		display:inline-block;
	}
	</style>
</head>
<body>
<header id="top-head" class="uk-position-fixed">
	<div class="uk-container uk-container-expand uk-background-primary">
		<nav class="uk-navbar uk-light" data-uk-navbar="mode:click; duration: 250" uk-navbar>
			<div class="uk-navbar-left">
				<ul class="uk-navbar-nav">
					<li class="uk-active"><a href="/">Jobkit</a></li>
				</ul>
			</div>
			<div class="uk-navbar-center">
				<div class="uk-navbar-item uk-visible@s">
					<form action="/search" class="uk-search uk-search-default">
						<span data-uk-search-icon="" class="uk-icon uk-search-icon"><svg width="20" height="20" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg" data-svg="search-icon"><circle fill="none" stroke="#000" stroke-width="1.1" cx="9" cy="9" r="7"></circle><path fill="none" stroke="#000" stroke-width="1.1" d="M14,14 L18,18 L14,14 Z"></path></svg></span>
						<input class="uk-search-input search-field" type="search" placeholder="Search">
					</form>
				</div>
			</div>
			<div class="uk-navbar-right">
				<div class="uk-navbar-nav">
					Hello
				</div>
			</div>
		</nav>
	</div>
</header>
{{ end }}
`
