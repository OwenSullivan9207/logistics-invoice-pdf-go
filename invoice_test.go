package main

import "testing"

func TestInvoiceEligibility(t *testing.T) {
	base := Order{ID: "L-42", Recipient: "billing@example.com", AmountCents: 7250, Events: []ShipmentEvent{{Kind: "delivered", At: "2026-09-20T10:00:00Z"}}, Proofs: []DeliveryProof{{Name: "signed receipt", Reference: "POD-42"}}}
	for _, tc := range []struct {
		name    string
		change  func(*Order)
		allowed bool
	}{
		{"delivered with proof", func(*Order) {}, true},
		{"missing proof", func(o *Order) { o.Proofs = nil }, false},
		{"not delivered", func(o *Order) { o.Events = nil }, false},
		{"open exception", func(o *Order) { o.Exceptions = []ShipmentException{{Reason: "damage"}} }, false},
		{"resolved exception", func(o *Order) { o.Exceptions = []ShipmentException{{Reason: "address correction", Resolved: true}} }, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o := base
			tc.change(&o)
			if (eligible(o) == nil) != tc.allowed {
				t.Fatalf("eligible=%v, want %v", eligible(o) == nil, tc.allowed)
			}
		})
	}
}
