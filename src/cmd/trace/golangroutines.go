// Copyright 2014 The Golang Authors. All rights reserved.
// Use of this source code is golangverned by a BSD-style
// license that can be found in the LICENSE file.

// Golangroutine-related profiles.

package main

import (
	"cmp"
	"fmt"
	"html/template"
	"internal/trace"
	"internal/trace/traceviewer"
	"log"
	"net/http"
	"slices"
	"sort"
	"strings"
	"time"
)

// GolangroutinesHandlerFunc returns a HandlerFunc that serves list of golangroutine groups.
func GolangroutinesHandlerFunc(summaries map[trace.GolangID]*trace.GolangroutineSummary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// golangroutineGroup describes a group of golangroutines grouped by name.
		type golangroutineGroup struct {
			Name     string        // Start function.
			N        int           // Total number of golangroutines in this group.
			ExecTime time.Duration // Total execution time of all golangroutines in this group.
		}
		// Accumulate groups by Name.
		groupsByName := make(map[string]golangroutineGroup)
		for _, summary := range summaries {
			group := groupsByName[summary.Name]
			group.Name = summary.Name
			group.N++
			group.ExecTime += summary.ExecTime
			groupsByName[summary.Name] = group
		}
		var groups []golangroutineGroup
		for _, group := range groupsByName {
			groups = append(groups, group)
		}
		slices.SortFunc(groups, func(a, b golangroutineGroup) int {
			return cmp.Compare(b.ExecTime, a.ExecTime)
		})
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		if err := templGolangroutines.Execute(w, groups); err != nil {
			log.Printf("failed to execute template: %v", err)
			return
		}
	}
}

var templGolangroutines = template.Must(template.New("").Parse(`
<html>
<style>` + traceviewer.CommonStyle + `
table {
  border-collapse: collapse;
}
td,
th {
  border: 1px solid black;
  padding-left: 8px;
  padding-right: 8px;
  padding-top: 4px;
  padding-bottom: 4px;
}
</style>
<body>
<h1>Golangroutines</h1>
Below is a table of all golangroutines in the trace grouped by start location and sorted by the total execution time of the group.<br>
<br>
Click a start location to view more details about that group.<br>
<br>
<table>
  <tr>
    <th>Start location</th>
	<th>Count</th>
	<th>Total execution time</th>
  </tr>
{{range $}}
  <tr>
    <td><code><a href="/golangroutine?name={{.Name}}">{{or .Name "(Inactive, no stack trace sampled)"}}</a></code></td>
	<td>{{.N}}</td>
	<td>{{.ExecTime}}</td>
  </tr>
{{end}}
</table>
</body>
</html>
`))

// GolangroutineHandler creates a handler that serves information about
// golangroutines in a particular group.
func GolangroutineHandler(summaries map[trace.GolangID]*trace.GolangroutineSummary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		golangroutineName := r.FormValue("name")

		type golangroutine struct {
			*trace.GolangroutineSummary
			NonOverlappingStats map[string]time.Duration
			HasRangeTime        bool
		}

		// Collect all the golangroutines in the group.
		var (
			golangroutines              []golangroutine
			name                    string
			totalExecTime, execTime time.Duration
			maxTotalTime            time.Duration
		)
		validNonOverlappingStats := make(map[string]struct{})
		validRangeStats := make(map[string]struct{})
		for _, summary := range summaries {
			totalExecTime += summary.ExecTime

			if summary.Name != golangroutineName {
				continue
			}
			nonOverlappingStats := summary.NonOverlappingStats()
			for name := range nonOverlappingStats {
				validNonOverlappingStats[name] = struct{}{}
			}
			var totalRangeTime time.Duration
			for name, dt := range summary.RangeTime {
				validRangeStats[name] = struct{}{}
				totalRangeTime += dt
			}
			golangroutines = append(golangroutines, golangroutine{
				GolangroutineSummary:    summary,
				NonOverlappingStats: nonOverlappingStats,
				HasRangeTime:        totalRangeTime != 0,
			})
			name = summary.Name
			execTime += summary.ExecTime
			if maxTotalTime < summary.TotalTime {
				maxTotalTime = summary.TotalTime
			}
		}

		// Compute the percent of total execution time these golangroutines represent.
		execTimePercent := ""
		if totalExecTime > 0 {
			execTimePercent = fmt.Sprintf("%.2f%%", float64(execTime)/float64(totalExecTime)*100)
		}

		// Sort.
		sortBy := r.FormValue("sortby")
		if _, ok := validNonOverlappingStats[sortBy]; ok {
			slices.SortFunc(golangroutines, func(a, b golangroutine) int {
				return cmp.Compare(b.NonOverlappingStats[sortBy], a.NonOverlappingStats[sortBy])
			})
		} else {
			// Sort by total time by default.
			slices.SortFunc(golangroutines, func(a, b golangroutine) int {
				return cmp.Compare(b.TotalTime, a.TotalTime)
			})
		}

		// Write down all the non-overlapping stats and sort them.
		allNonOverlappingStats := make([]string, 0, len(validNonOverlappingStats))
		for name := range validNonOverlappingStats {
			allNonOverlappingStats = append(allNonOverlappingStats, name)
		}
		slices.SortFunc(allNonOverlappingStats, func(a, b string) int {
			if a == b {
				return 0
			}
			if a == "Execution time" {
				return -1
			}
			if b == "Execution time" {
				return 1
			}
			return cmp.Compare(a, b)
		})

		// Write down all the range stats and sort them.
		allRangeStats := make([]string, 0, len(validRangeStats))
		for name := range validRangeStats {
			allRangeStats = append(allRangeStats, name)
		}
		sort.Strings(allRangeStats)

		err := templGolangroutine.Execute(w, struct {
			Name                string
			N                   int
			ExecTimePercent     string
			MaxTotal            time.Duration
			Golangroutines          []golangroutine
			NonOverlappingStats []string
			RangeStats          []string
		}{
			Name:                name,
			N:                   len(golangroutines),
			ExecTimePercent:     execTimePercent,
			MaxTotal:            maxTotalTime,
			Golangroutines:          golangroutines,
			NonOverlappingStats: allNonOverlappingStats,
			RangeStats:          allRangeStats,
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to execute template: %v", err), http.StatusInternalServerError)
			return
		}
	}
}

func stat2Color(statName string) string {
	color := "#636363"
	if strings.HasPrefix(statName, "Block time") {
		color = "#d01c8b"
	}
	switch statName {
	case "Sched wait time":
		color = "#2c7bb6"
	case "Syscall execution time":
		color = "#7b3294"
	case "Execution time":
		color = "#d7191c"
	}
	return color
}

var templGolangroutine = template.Must(template.New("").Funcs(template.FuncMap{
	"percent": func(dividend, divisor time.Duration) template.HTML {
		if divisor == 0 {
			return ""
		}
		return template.HTML(fmt.Sprintf("(%.1f%%)", float64(dividend)/float64(divisor)*100))
	},
	"headerStyle": func(statName string) template.HTMLAttr {
		return template.HTMLAttr(fmt.Sprintf("style=\"background-color: %s;\"", stat2Color(statName)))
	},
	"barStyle": func(statName string, dividend, divisor time.Duration) template.HTMLAttr {
		width := "0"
		if divisor != 0 {
			width = fmt.Sprintf("%.2f%%", float64(dividend)/float64(divisor)*100)
		}
		return template.HTMLAttr(fmt.Sprintf("style=\"width: %s; background-color: %s;\"", width, stat2Color(statName)))
	},
}).Parse(`
<!DOCTYPE html>
<title>Golangroutines: {{.Name}}</title>
<style>` + traceviewer.CommonStyle + `
th {
  background-color: #050505;
  color: #fff;
}
th.link {
  cursor: pointer;
}
table {
  border-collapse: collapse;
}
td,
th {
  padding-left: 8px;
  padding-right: 8px;
  padding-top: 4px;
  padding-bottom: 4px;
}
.details tr:hover {
  background-color: #f2f2f2;
}
.details td {
  text-align: right;
  border: 1px solid black;
}
.details td.id {
  text-align: left;
}
.stacked-bar-graph {
  width: 300px;
  height: 10px;
  color: #414042;
  white-space: nowrap;
  font-size: 5px;
}
.stacked-bar-graph span {
  display: inline-block;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
  float: left;
  padding: 0;
}
</style>

<script>
function reloadTable(key, value) {
  let params = new URLSearchParams(window.location.search);
  params.set(key, value);
  window.location.search = params.toString();
}
</script>

<h1>Golangroutines</h1>

Table of contents
<ul>
	<li><a href="#summary">Summary</a></li>
	<li><a href="#breakdown">Breakdown</a></li>
	<li><a href="#ranges">Special ranges</a></li>
</ul>

<h3 id="summary">Summary</h3>

<table class="summary">
	<tr>
		<td>Golangroutine start location:</td>
		<td><code>{{.Name}}</code></td>
	</tr>
	<tr>
		<td>Count:</td>
		<td>{{.N}}</td>
	</tr>
	<tr>
		<td>Execution Time:</td>
		<td>{{.ExecTimePercent}} of total program execution time </td>
	</tr>
	<tr>
		<td>Network wait profile:</td>
		<td> <a href="/io?name={{.Name}}">graph</a> <a href="/io?name={{.Name}}&raw=1" download="io.profile">(download)</a></td>
	</tr>
	<tr>
		<td>Sync block profile:</td>
		<td> <a href="/block?name={{.Name}}">graph</a> <a href="/block?name={{.Name}}&raw=1" download="block.profile">(download)</a></td>
	</tr>
	<tr>
		<td>Syscall profile:</td>
		<td> <a href="/syscall?name={{.Name}}">graph</a> <a href="/syscall?name={{.Name}}&raw=1" download="syscall.profile">(download)</a></td>
		</tr>
	<tr>
		<td>Scheduler wait profile:</td>
		<td> <a href="/sched?name={{.Name}}">graph</a> <a href="/sched?name={{.Name}}&raw=1" download="sched.profile">(download)</a></td>
	</tr>
</table>

<h3 id="breakdown">Breakdown</h3>

The table below breaks down where each golangroutine is spent its time during the
traced period.
All of the columns except total time are non-overlapping.
<br>
<br>

<table class="details">
<tr>
<th> Golangroutine</th>
<th class="link" onclick="reloadTable('sortby', 'Total time')"> Total</th>
<th></th>
{{range $.NonOverlappingStats}}
<th class="link" onclick="reloadTable('sortby', '{{.}}')" {{headerStyle .}}> {{.}}</th>
{{end}}
</tr>
{{range .Golangroutines}}
	<tr>
		<td> <a href="/trace?golangid={{.ID}}">{{.ID}}</a> </td>
		<td> {{ .TotalTime.String }} </td>
		<td>
			<div class="stacked-bar-graph">
			{{$Golangroutine := .}}
			{{range $.NonOverlappingStats}}
				{{$Time := index $Golangroutine.NonOverlappingStats .}}
				{{if $Time}}
					<span {{barStyle . $Time $.MaxTotal}}>&nbsp;</span>
				{{end}}
			{{end}}
			</div>
		</td>
		{{$Golangroutine := .}}
		{{range $.NonOverlappingStats}}
			{{$Time := index $Golangroutine.NonOverlappingStats .}}
			<td> {{$Time.String}}</td>
		{{end}}
	</tr>
{{end}}
</table>

<h3 id="ranges">Special ranges</h3>

The table below describes how much of the traced period each golangroutine spent in
certain special time ranges.
If a golangroutine has spent no time in any special time ranges, it is excluded from
the table.
For example, how much time it spent helping the GC. Note that these times do
overlap with the times from the first table.
In general the golangroutine may not be executing in these special time ranges.
For example, it may have blocked while trying to help the GC.
This must be taken into account when interpreting the data.
<br>
<br>

<table class="details">
<tr>
<th> Golangroutine</th>
<th> Total</th>
{{range $.RangeStats}}
<th {{headerStyle .}}> {{.}}</th>
{{end}}
</tr>
{{range .Golangroutines}}
	{{if .HasRangeTime}}
		<tr>
			<td> <a href="/trace?golangid={{.ID}}">{{.ID}}</a> </td>
			<td> {{ .TotalTime.String }} </td>
			{{$Golangroutine := .}}
			{{range $.RangeStats}}
				{{$Time := index $Golangroutine.RangeTime .}}
				<td> {{$Time.String}}</td>
			{{end}}
		</tr>
	{{end}}
{{end}}
</table>
`))
