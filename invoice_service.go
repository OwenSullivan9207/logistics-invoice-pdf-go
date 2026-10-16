package main

import (
	"context"
	"fmt"
	"html"
	"net/http"
	"strings"
)

type invoiceResult struct {
	URL       string `json:"url"`
	MessageID string `json:"message_id"`
}

func sendInvoice(ctx context.Context, c infraiClient, o Order) (invoiceResult, error) {
	var result invoiceResult
	if err := eligible(o); err != nil {
		return result, err
	}
	var pdf struct {
		URL string `json:"url"`
	}
	if err := c.post(ctx, "/v1/pdf/generate", map[string]any{"html": invoiceHTML(o), "page_size": "A4", "store": true}, "invoice-pdf-"+o.ID, &pdf); err != nil {
		return result, err
	}
	if pdf.URL == "" {
		return result, fmt.Errorf("PDF response has no URL")
	}
	result.URL = pdf.URL
	var mail struct {
		MessageID string `json:"message_id"`
	}
	body := "<p>Freight invoice for order " + html.EscapeString(o.ID) + ": <a href=\"" + html.EscapeString(pdf.URL) + "\">Download PDF</a></p>"
	if err := c.post(ctx, "/v1/email/send", map[string]any{"to": o.Recipient, "subject": "Freight invoice " + o.ID, "html": body}, "invoice-mail-"+o.ID, &mail); err != nil {
		return result, err
	}
	result.MessageID = mail.MessageID
	return result, nil
}

func clientStatus(err error) int {
	if _, ok := err.(APIError); ok {
		return http.StatusBadRequest
	}
	if strings.Contains(err.Error(), "required") || strings.Contains(err.Error(), "unresolved") {
		return http.StatusUnprocessableEntity
	}
	return http.StatusBadGateway
}
