package jobkit

var invocationTemplate = `
{{ define "invocation" }}
{{ template "header" . }}
<div class="container">
	<ul class="breadcrumbs">
		<li><a href="/">Jobs</a></li>
		<li>{{ .ViewModel.JobName }}</li>
		<li>{{ .ViewModel.ID }}</li>
	</ul>
	<table class="u-full-width">
		<thead>
			<tr>
				<th>Invocation</th>
				<th>Started</th>
				<th>Finished</th>
				<th>Timeout</th>
				<th>Cancelled</th>
				<th>Elapsed</th>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>{{ .ViewModel.ID }}</td>
				<td>{{ .ViewModel.Started | rfc3339 }}</td>
				<td>{{ if .ViewModel.Finished.IsZero }}-{{ else }}{{ .ViewModel.Finished | rfc3339 }}{{ end }}</td>
				<td>{{ if .ViewModel.Timeout.IsZero }}-{{ else }}{{ .ViewModel.Timeout | rfc3339 }}{{ end }}</td>
				<td>{{ if .ViewModel.Cancelled.IsZero }}-{{ else }}{{ .ViewModel.Cancelled | rfc3339 }}{{ end }}</td>
				<td>{{ if .ViewModel.Finished.IsZero }}{{ .ViewModel.Started | since_utc }}{{ else }}{{ .ViewModel.Elapsed }}{{ end }}</td>
			</tr>
		</tbody>
	</table>
	{{ if .ViewModel.Err }}
	<table class="u-full-width">
		<thead>
			<tr>
				<th>Error</th>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>
					<pre>{{ .ViewModel.Err }}</pre>
				</td>
			</tr>
		</tbody>
	</table>
	{{ end }}
	<table class="u-full-width">
		<tbody>
			<tr>
				<td class="align-right">
					<a href="/api/job.invocation.output/{{ .ViewModel.JobName }}/{{ .ViewModel.ID }}">JSON</a>
					<a href="/job.invocation.output/{{ .ViewModel.JobName }}/{{ .ViewModel.ID }}">Raw</a>
				</td>
			</tr>
		</tbody>
	</table>

	<!-- deps -->
	<link rel="stylesheet" href="/static/css/xterm.css" />
	<script src="/static/js/xterm.js"></script>
	<script src="/static/js/addons/fit/fit.js"></script>
	<script src="/static/js/addons/webLinks/webLinks.js"></script>
	<script src="/static/js/jquery.min.js"></script>

	{{ if .ViewModel.Output }}
	<table class="u-full-width">
		<thead>
			<tr>
				<th>Output (Combined)</th>
			</tr>
		</thead>
		<tbody>
			<tr>
				<td>
					<div id="terminal"></div>
					<script>
						Terminal.applyAddon(fit);
						Terminal.applyAddon(webLinks);
						var term = new Terminal();
						term.open(document.getElementById('terminal'));
						term.fit();

						const interval = 2000;
						var lastUpdateNanos = '';

						var getLogsURL = () => {
							var logsURL = "/api/job.invocation.output/{{ .ViewModel.JobName }}/{{ .ViewModel.ID }}";

							if (lastUpdateNanos != '') {
								logsURL = logsURL + '?afterNanos=' + lastUpdateNanos;
							}
							return logsURL;
						};

						var getLogs = (cb) => {
							$.get(getLogsURL()).then((res) => {
								lastUpdateNanos = res.serverTimeNanos.toString();
								if (res.lines != null) {
									for(var i = 0; i < res.lines.length; i++) {
										term.write(res.lines[i].line + '\r\n');
									}
								}
								if (cb) {
									cb();
								}
								return
							}).fail((res) => {
								cb();
							});
						};
						var poll = () => {
							getLogs(() => {
								setTimeout(poll, interval);
							});
						};
						poll();
					</script>
				</td>
			</tr>
		</tbody>
	</table>
	{{ end }}
</div>
{{ template "footer" . }}
{{ end }}
`
