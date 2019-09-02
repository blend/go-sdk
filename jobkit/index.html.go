package jobkit

var indexTemplate = `
{{ define "index" }}
{{ template "header" . }}
<div class="container">
		{{ range $index, $job := .ViewModel.Jobs }}
		<table class="job u-full-width">
			<thead>
				<tr>
					<th>Job Name</th>
					<th>Schedule</th>
					<th>Current</th>
					<th>Next Run</th>
					<th>Last Ran</th>
					<th>Last Result</th>
					<th>Last Elapsed</th>
					<th>Actions</th>
				</tr>
			</thead>
			<tbody>
				<tr>
					<td> <!-- job name -->
						{{ $job.Name }}
					</td>
					<td> <!-- schedule -->
						<pre>{{ $job.Schedule }}</pre>
						</td>
					<td> <!-- current -->
					{{ if $job.Current }}
						{{ $job.Current.Started | since_utc }}
					{{else}}
						<span>-</span>
					{{end}}
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
						<form method="POST" action="/job.enable/{{ $job.Name }}">
							<input type="submit" class="button" value="Enable" />
						</form>
					{{else}}
						<form method="POST" action="/job.disable/{{ $job.Name }}">
							<input type="submit" class="button" value="Disable" />
						</form>
					{{end}}
					{{ if $job.Current }}
					<form method="POST" action="/job.cancel/{{ $job.Name }}">
						<input type="submit" class="button button-danger" value="Cancel" />
					</form>
					{{else}}
					<form method="POST" action="/job.run/{{ $job.Name }}">
						<input type="submit" class="button button-primary" value="Run" />
					</form>
					{{end}}
					</td>
				</tr>
				{{ if $job.Description }}
				<tr>
					<td colspan=8>
						<h4>Description</h4>
						<pre>{{ $job.Description }}</pre>
					</td>
				</tr>
				{{ end }}
				<tr>
					<td colspan=8>
						<h4>History</h4>
						<table class="u-full-width">
							<thead>
								<tr>
									<th>Invocation</th>
									<th>Started</th>
									<th>Finished</th>
									<th>Timeout</th>
									<th>Cancelled</th>
									<th>Elapsed</th>
									<th>Error</th>
								</tr>
							</thead>
							<tbody>

							{{ if $job.Current }}
							<tr class="running">
								<td><a href="/job.invocation/{{$job.Name}}/{{ $job.Current.ID }}">{{ $job.Current.ID }}</a></td>
								<td>{{ $job.Current.Started | rfc3339 }}</td>
								<td>{{ if $job.Current.Finished.IsZero }}-{{ else }}{{ $job.Current.Finished | rfc3339 }}{{ end }}</td>
								<td>{{ if $job.Current.Timeout.IsZero }}-{{ else }}{{ $job.Current.Timeout | rfc3339 }}{{ end }}</td>
								<td>-</td>
								<td>{{ $job.Current.Elapsed }}</td>
								<td>-</td>
							<tr>
							{{ end }}

							{{ range $index, $ji := $job.History | reverse }}
							<tr class="{{ if $ji.Status | eq "failed" }}failed{{ else if $ji.Status | eq "cancelled"}}cancelled{{else}}ok{{end}}">
								<td><a href="/job.invocation/{{$ji.JobName}}/{{ $ji.ID }}">{{ $ji.ID }}</a></td>
								<td>{{ $ji.Started | rfc3339 }}</td>
								<td>{{ if $ji.Finished.IsZero }}-{{ else }}{{ $ji.Finished | rfc3339 }}{{ end }}</td>
								<td>{{ if $ji.Timeout.IsZero }}-{{ else }}{{ $ji.Timeout | rfc3339 }}{{ end }}</td>
								<td>{{ if $ji.Cancelled.IsZero }}-{{ else }}{{ $ji.Cancelled | rfc3339 }}{{ end }}</td>
								<td>{{ $ji.Elapsed }}</td>
								<td>{{ if $ji.Err }}<code>{{ $ji.Err }}</code>{{ else }}-{{end}}</td>
							</tr>
							{{ end }}
							</tbody>
						</table>
					</td>
				</tr>
			</tbody>
		</table>
		{{ else }}
		<h4>No Jobs Loaded</h4>
		{{ end }}
</div>
{{ template "footer" . }}
{{ end }}
`
