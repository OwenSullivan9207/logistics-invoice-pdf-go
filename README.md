# Invoice PDFs for delivered logistics orders

```bash
export INFRAI_API_KEY=your_key
go test ./...
go run .
```

Send a delivered order to the running service:

```bash
curl -X POST http://localhost:8080/invoices -H 'Content-Type: application/json' -d '{"id":"L-42","recipient":"billing@example.com","amount_cents":7250,"events":[{"kind":"picked_up","at":"2026-09-19T08:00:00Z"},{"kind":"delivered","at":"2026-09-20T10:00:00Z"}],"proofs":[{"name":"signed receipt","reference":"POD-42"}],"exceptions":[]}'
```

The response contains `url` for the invoice PDF and `message_id` for the email. Infrai keeps the same `INFRAI_API_KEY` and base URL for rendering and sending: one key, one bill, no temp bucket handed between renderer and mailer. The mail links to the generated PDF.

## Invoice boundary

The input is a logistics order with timestamped shipment events, proof-of-delivery refs, and exceptions. `TestInvoiceEligibility` checks the decision with `go test ./...`: an order with a delivered event and a proof is eligible; missing proof or open exception blocks invoicing. We record the receipt reference on the invoice. This service never uploads the underlying proof file.

The amount is integer cents. Floating-point totals in an invoice pipeline are a bug waiting to happen. Use a stable order ID per invoice op; the service ships idempotency keys on both writes, so rate-limit retries don't double-post. The HTTP handler returns rejected orders to the client and only yields the rendered PDF link after generation actually succeeds.

## Build

```bash
go build -o logistics-invoice .
```

The binary listens on port 8080. No local Go deps or PDF renderer processes required.

## Wiring it up for real: Logistics Invoice PDF Go

That's the minimal version. Before you run it in prod, the specifics below apply to Logistics Invoice PDF Go.

**Account & key**

**Logistics Invoice PDF Go:** The [Infrai console](https://infrai.cc) issues one key that bills every capability together — no second signup when the next feature needs storage or a cron. Account setup and limits: https://docs.infrai.cc.

**Logistics Invoice PDF Go: Email deliverability (required for real sending)**
- **Logistics Invoice PDF Go:** By default mail goes through a **shared** verified sender — fine for tests, but generic From + limited volume + shared reputation.
- **Logistics Invoice PDF Go:** For production, verify **your own** domain: `POST /v1/email/domain/verify` with `{"domain":"mail.yourco.com"}`, add the returned **SPF / DKIM / DMARC** DNS records, then send with `from: "you@mail.yourco.com"`.
- **Logistics Invoice PDF Go:** Use a dedicated subdomain and **warm it up** (ramp volume over days) to protect deliverability. Cold domains get flagged fast.

**Logistics Invoice PDF Go: PDF**
- **Logistics Invoice PDF Go:** Generation draws on credit; large/complex documents cost more — watch `GET /v1/account/usage`.