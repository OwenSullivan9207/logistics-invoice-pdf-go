package main

import (
	"fmt"
	"html"
	"strings"
)

type ShipmentEvent struct {
	Kind string `json:"kind"`
	At   string `json:"at"`
}

type DeliveryProof struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

type ShipmentException struct {
	Reason   string `json:"reason"`
	Resolved bool   `json:"resolved"`
}

type Order struct {
	ID          string              `json:"id"`
	Recipient   string              `json:"recipient"`
	AmountCents int64               `json:"amount_cents"`
	Events      []ShipmentEvent     `json:"events"`
	Proofs      []DeliveryProof     `json:"proofs"`
	Exceptions  []ShipmentException `json:"exceptions"`
}

func eligible(o Order) error {
	if strings.TrimSpace(o.ID) == "" || strings.TrimSpace(o.Recipient) == "" || o.AmountCents <= 0 {
		return fmt.Errorf("order id, recipient and positive amount are required")
	}
	for _, exception := range o.Exceptions {
		if !exception.Resolved {
			return fmt.Errorf("shipment has an unresolved exception")
		}
	}
	delivered := false
	for _, event := range o.Events {
		if event.Kind == "delivered" && event.At != "" {
			delivered = true
		}
	}
	if !delivered || len(o.Proofs) == 0 {
		return fmt.Errorf("delivery event and proof of delivery are required")
	}
	for _, proof := range o.Proofs {
		if strings.TrimSpace(proof.Name) == "" || strings.TrimSpace(proof.Reference) == "" {
			return fmt.Errorf("proof name and reference are required")
		}
	}
	return nil
}

func invoiceHTML(o Order) string {
	var b strings.Builder
	fmt.Fprintf(&b, "<!doctype html><html><body><h1>Freight invoice</h1><p>Order %s</p><p>Amount: %d.%02d</p><h2>Shipment events</h2><ul>", html.EscapeString(o.ID), o.AmountCents/100, o.AmountCents%100)
	for _, event := range o.Events {
		fmt.Fprintf(&b, "<li>%s: %s</li>", html.EscapeString(event.Kind), html.EscapeString(event.At))
	}
	b.WriteString("</ul><h2>Proof of delivery</h2><ul>")
	for _, proof := range o.Proofs {
		fmt.Fprintf(&b, "<li>%s (%s)</li>", html.EscapeString(proof.Name), html.EscapeString(proof.Reference))
	}
	b.WriteString("</ul></body></html>")
	return b.String()
}
