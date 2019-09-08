package jobkit

var indexTemplate = `
{{ define "index" }}
{{ template "header" . }}
	<div id="content" class="uk-container uk-container-expand">
		<table class="uk-table uk-table-striped">
			<thead>
				<tr>
					<th>Job Name</th>
					<th>Success Rate</th>
					<th>Schedule</th>
					<th>Next Run</th>
					<th>Last Ran</th>
					<th>Last Result</th>
					<th>Last Elapsed</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
		{{ range $index, $job := .ViewModel.Jobs }}
				<tr class="job-details">
					<td> <!-- job name -->
						{{ $job.Name }}
					</td>
					<td> <!-- success rate -->
						{{ $job.Stats.SuccessRate | format_pct }}
					</td>
					<td> <!-- schedule -->
					{{ if $job.Schedule }}
						{{ $job.Schedule }}
					{{ else }}
						<span>-</span>
					{{ end }}
					</td>
					<td> <!-- next run-->
					{{ if $job.Disabled }}
						<span>-</span>
					{{ else }}
						{{ $job.NextRuntime | rfc3339 }}
					{{ end }}
					</td>
					<td> <!-- last run -->
					{{ if $job.Last }}
						{{ $job.Last.Started | rfc3339 }}
					{{ else }}
						<span class="none">-</span>
					{{ end }}
					</td>
					<td> <!-- last status -->
					{{ if $job.Last }}
						{{ if $job.Last.Err }}
							{{ $job.Last.Err }}
						{{ else }}
						<span class="none">Success</span>
						{{ end }}
					{{ else }}
						<span class="none">-</span>
					{{ end }}
					</td>
					<td> <!-- last elapsed -->
					{{ if $job.Last }}
						{{ $job.Last.Elapsed }}
					{{ else }}
						<span class="none">-</span>
					{{ end }}
					</td>
					<td> <!-- actions -->
					{{ if $job.Disabled }}
						<form method="POST" action="/job.enable/{{ $job.Name }}" class="uk-form">
							<input type="submit" class="uk-button" value="Enable" />
						</form>
					{{else}}
						<form method="POST" action="/job.disable/{{ $job.Name }}" class="uk-form">
							<input type="submit" class="uk-button" value="Disable" />
						</form>
					{{end}}
					{{ if $job.CanRun }}
					<form method="POST" action="/job.run/{{ $job.Name }}" class="uk-form">
						<input type="submit" class="uk-button uk-button-primary" value="Run" />
					</form>
					{{ else }}
					<button class="uk-button uk-button-disabled" disabled>Running</button>
					{{ end }}
					</td>
				</tr>
		{{ else }}
		<tr colspan=8><td>No Jobs Loaded</td></tr>
		{{ end }}
			</tbody>
		</table>
</div>
{{ template "footer" . }}
{{ end }}
`
