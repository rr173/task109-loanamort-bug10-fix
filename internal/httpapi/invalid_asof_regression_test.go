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
	r := httptest.NewRequest("GET", "/loans/"+l.LoanID+"/accrued-interest?as_of=abc", nil)
	w := httptest.NewRecorder()
	httpapi.NewMux(svc, "").ServeHTTP(w, r)
	if w.Code != 400 {
		t.Fatalf("status=%d, body=%s", w.Code, w.Body.String())
	}
}
