package httpapi_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"task109-loanamort/internal/httpapi"
	"task109-loanamort/internal/loan"
	"task109-loanamort/internal/store"
)

func TestAccruedInterestRejectsInvalidAsOf(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/loan.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	svc := loan.New(db)
	b, err := svc.CreateBorrower(ctx, loan.CreateBorrowerRequest{Name: "borrower"})
	if err != nil {
		t.Fatal(err)
	}
	l, err := svc.CreateLoan(ctx, loan.CreateLoanRequest{BorrowerID: b.BorrowerID, PrincipalCents: 100000, AnnualPercent: 12, Periods: 4, Type: loan.EqualInstallment})
	if err != nil {
		t.Fatal(err)
	}

	// A non-integer, a negative integer, a sign-prefixed value, and a
	// value with trailing non-digits must all be rejected with 400 rather
	// than silently falling back to a default period.
	for _, asOf := range []string{"abc", "-1", "+3", "3x", "1.5"} {
		r := httptest.NewRequest("GET", "/loans/"+l.LoanID+"/accrued-interest?as_of="+asOf, nil)
		w := httptest.NewRecorder()
		httpapi.NewMux(svc, "").ServeHTTP(w, r)
		if w.Code != 400 {
			t.Errorf("as_of=%q: status=%d, body=%s, want 400", asOf, w.Code, w.Body.String())
		}
	}

	// A valid non-negative integer must still succeed (regression guard so
	// the validation does not over-reject well-formed input).
	r := httptest.NewRequest("GET", "/loans/"+l.LoanID+"/accrued-interest?as_of=2", nil)
	w := httptest.NewRecorder()
	httpapi.NewMux(svc, "").ServeHTTP(w, r)
	if w.Code != 200 {
		t.Fatalf("as_of=2: status=%d, body=%s, want 200", w.Code, w.Body.String())
	}
}
