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
				<div id="terminal">
					<pre>{{ .ViewModel.Output }}</pre>
				</div>
				<script>
					var es = new EventSource("/api/job.invocation.output.stream/{{ .ViewModel.JobName }}/{{ .ViewModel.ID }}");
					var terminal = document.getElementById('terminal')
					es.onmessage = (e) => {
						var line = document.createElement('pre');
						line.textContent = e.data;
						terminal.appendChild(line);
					};
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
