package lint

import (
	"testing"
)

func TestJobDatasource(t *testing.T) {
	linter := NewTemplateJobRule()

	for _, tc := range []struct {
		result    Result
		dashboard Dashboard
	}{
		// Non-prometheus dashboards shouldn't fail.
		{
			result: Result{
				Severity: Success,
				Message:  "OK",
			},
			dashboard: Dashboard{
				Title: "test",
			},
		},
		// Missing job template.
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test' is missing the job template",
			},
			dashboard: Dashboard{
				Title: "test",
				Templating: struct {
					List []Template `json:"list"`
				}{
					List: []Template{
						{
							Type:  "datasource",
							Query: "prometheus",
						},
					},
				},
			},
		},
		// Wrong datasource.
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test' job template should use datasource $datasource, is currently 'foo'",
			},
			dashboard: Dashboard{
				Title: "test",
				Templating: struct {
					List []Template `json:"list"`
				}{
					List: []Template{
						{
							Type:  "datasource",
							Query: "prometheus",
						},
						{
							Name:       "job",
							Datasource: "foo",
						},
					},
				},
			},
		},
		// Wrong datasource (multiple datasources).
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test' job template should use datasource $prometheus_datasource or $loki_datasource, is currently 'foo'",
			},
			dashboard: Dashboard{
				Title: "test",
				Templating: struct {
					List []Template `json:"list"`
				}{
					List: []Template{
						{
							Type:  "datasource",
							Query: "prometheus",
						},
						{
							Type:  "datasource",
							Query: "loki",
						},
						{
							Name:       "job",
							Datasource: "foo",
						},
					},
				},
			},
		},
		// Wrong type.
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test' job template should be a Prometheus query, is currently 'bar'",
			},
			dashboard: Dashboard{
				Title: "test",
				Templating: struct {
					List []Template `json:"list"`
				}{
					List: []Template{
						{
							Type:  "datasource",
							Query: "prometheus",
						},
						{
							Name:       "job",
							Datasource: "$datasource",
							Type:       "bar",
						},
					},
				},
			},
		},
		// Wrong job label.
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test' job template should be a labelled 'job', is currently 'bar'",
			},
			dashboard: Dashboard{
				Title: "test",
				Templating: struct {
					List []Template `json:"list"`
				}{
					List: []Template{
						{
							Type:  "datasource",
							Query: "prometheus",
						},
						{
							Name:       "job",
							Datasource: "$datasource",
							Type:       "query",
							Label:      "bar",
						},
					},
				},
			},
		},
		// Missing instance templates.
		{
			result: Result{
				Severity: Error,
				Message:  "Dashboard 'test' is missing the instance template",
			},
			dashboard: Dashboard{
				Title: "test",
				Templating: struct {
					List []Template `json:"list"`
				}{
					List: []Template{
						{
							Type:  "datasource",
							Query: "prometheus",
						},
						{
							Name:       "job",
							Datasource: "$datasource",
							Type:       "query",
							Label:      "job",
							Multi:      true,
							AllValue:   ".+",
						},
					},
				},
			},
		},
		// What success looks like.
		{
			result: Result{
				Severity: Success,
				Message:  "OK",
			},
			dashboard: Dashboard{
				Title: "test",
				Templating: struct {
					List []Template `json:"list"`
				}{
					List: []Template{
						{
							Type:  "datasource",
							Query: "prometheus",
						},
						{
							Name:       "job",
							Datasource: "$datasource",
							Type:       "query",
							Label:      "job",
							Multi:      true,
							AllValue:   ".+",
						},
						{
							Name:       "instance",
							Datasource: "${datasource}",
							Type:       "query",
							Label:      "instance",
							Multi:      true,
							AllValue:   ".+",
						},
					},
				},
			},
		},
	} {
		testRule(t, linter, tc.dashboard, tc.result)
	}
}
