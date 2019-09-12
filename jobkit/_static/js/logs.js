import { Terminal } from '/static/js/xterm.js'
const term = new Terminal();
term.open(document.getElementById('terminal'));
var es = new EventSource("/api/job.invocation.output.stream/{{ .ViewModel.JobName }}/{{ .ViewModel.ID }}");
es.onmessage = (e) => {
	term.write(e.data);
};