package domain_test

import (
	"reflect"
	"testing"

	"github.com/pt9912/u-boot/internal/hexagon/domain"
)

func TestSeverity_String(t *testing.T) {
	t.Parallel()
	cases := []struct {
		sev  domain.Severity
		want string
	}{
		{domain.SeverityOK, "ok"},
		{domain.SeverityWarn, "warn"},
		{domain.SeverityError, "error"},
		{domain.Severity(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.sev.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tc.sev, got, tc.want)
		}
	}
}

func TestSeverity_Ordering(t *testing.T) {
	t.Parallel()
	// Why: the LH-FA-DIAG-003 exit-code dispatch relies on the
	// numerical ordering SeverityOK < SeverityWarn < SeverityError.
	// Pin it so a future renumbering breaks the test instead of the
	// CLI semantics.
	if domain.SeverityOK >= domain.SeverityWarn || domain.SeverityWarn >= domain.SeverityError {
		t.Errorf("severity ordering broken: OK=%d Warn=%d Error=%d",
			domain.SeverityOK, domain.SeverityWarn, domain.SeverityError)
	}
}

func TestDiagnosticReport_MaxSeverity(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   domain.DiagnosticReport
		want domain.Severity
	}{
		{
			name: "empty",
			in:   domain.DiagnosticReport{},
			want: domain.SeverityOK,
		},
		{
			name: "all ok",
			in: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "a", Severity: domain.SeverityOK},
				{ID: "b", Severity: domain.SeverityOK},
			}},
			want: domain.SeverityOK,
		},
		{
			name: "ok+warn",
			in: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "a", Severity: domain.SeverityOK},
				{ID: "b", Severity: domain.SeverityWarn},
			}},
			want: domain.SeverityWarn,
		},
		{
			name: "warn+error",
			in: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "a", Severity: domain.SeverityWarn},
				{ID: "b", Severity: domain.SeverityError},
			}},
			want: domain.SeverityError,
		},
		{
			name: "single error",
			in: domain.DiagnosticReport{Items: []domain.Diagnostic{
				{ID: "a", Severity: domain.SeverityError},
			}},
			want: domain.SeverityError,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.in.MaxSeverity(); got != tc.want {
				t.Errorf("MaxSeverity() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDiagnosticReport_HasErrors(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   domain.DiagnosticReport
		want bool
	}{
		{"empty", domain.DiagnosticReport{}, false},
		{"only-ok", domain.DiagnosticReport{Items: []domain.Diagnostic{{Severity: domain.SeverityOK}}}, false},
		{"only-warn", domain.DiagnosticReport{Items: []domain.Diagnostic{{Severity: domain.SeverityWarn}}}, false},
		{"with-error", domain.DiagnosticReport{Items: []domain.Diagnostic{{Severity: domain.SeverityError}}}, true},
		{"warn+error", domain.DiagnosticReport{Items: []domain.Diagnostic{
			{Severity: domain.SeverityWarn}, {Severity: domain.SeverityError},
		}}, true},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.in.HasErrors(); got != tc.want {
				t.Errorf("HasErrors() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDiagnosticReport_HasWarnings(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		in   domain.DiagnosticReport
		want bool
	}{
		{"empty", domain.DiagnosticReport{}, false},
		{"only-ok", domain.DiagnosticReport{Items: []domain.Diagnostic{{Severity: domain.SeverityOK}}}, false},
		{"only-warn", domain.DiagnosticReport{Items: []domain.Diagnostic{{Severity: domain.SeverityWarn}}}, true},
		{
			// Warn shadowed by Error → HasWarnings is FALSE (the max
			// is Error, not Warn). LH-FA-DIAG-003 distinguishes
			// strictly between the two.
			"warn+error-prefers-error",
			domain.DiagnosticReport{Items: []domain.Diagnostic{
				{Severity: domain.SeverityWarn}, {Severity: domain.SeverityError},
			}},
			false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.in.HasWarnings(); got != tc.want {
				t.Errorf("HasWarnings() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestDiagnosticReport_SortedByIssuesFirst(t *testing.T) {
	t.Parallel()
	in := domain.DiagnosticReport{Items: []domain.Diagnostic{
		{ID: "z.ok", Severity: domain.SeverityOK},
		{ID: "a.error", Severity: domain.SeverityError},
		{ID: "b.warn", Severity: domain.SeverityWarn},
		{ID: "a.ok", Severity: domain.SeverityOK},
		{ID: "a.warn", Severity: domain.SeverityWarn},
	}}
	got := in.SortedByIssuesFirst()
	want := []domain.Diagnostic{
		{ID: "a.error", Severity: domain.SeverityError},
		{ID: "a.warn", Severity: domain.SeverityWarn},
		{ID: "b.warn", Severity: domain.SeverityWarn},
		{ID: "a.ok", Severity: domain.SeverityOK},
		{ID: "z.ok", Severity: domain.SeverityOK},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("SortedByIssuesFirst()\n got  %v\n want %v", got, want)
	}
	// Original report must not be mutated.
	if in.Items[0].ID != "z.ok" {
		t.Errorf("SortedByIssuesFirst mutated receiver: in.Items[0].ID = %q", in.Items[0].ID)
	}
}
