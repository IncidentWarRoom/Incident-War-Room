# Error Handling in incident-service (`internal/errs`)

Every error carries a **category** (`Kind`), so upper layers (bot handlers)
can decide how to respond without parsing error messages. Fully compatible
with `errors.Is` / `errors.As` — the cause chain is preserved via `Unwrap()`.

## Kinds

| Kind | When to use |
|---|---|
| `KindInternal` | unexpected system failure (default for unclassified errors) |
| `KindUnavailable` | external dependency unreachable: Postgres, Telegram API, report-service |
| `KindNotFound` | entity does not exist (e.g. no active incident) |
| `KindConflict` | operation contradicts current state (e.g. closing a closed incident) |
| `KindValidation` | malformed input / business-rule violation |

## Raising errors

Wrap the cause at the lowest level (repository, HTTP client):

```go
return errs.Wrap(errs.KindUnavailable, "storage.Ping", err)
return errs.Wrapf(errs.KindUnavailable, "reportclient.Generate", err, "call report service")
return errs.New(errs.KindValidation, "incident.Create", "title is empty")
```

`Wrap` / `Wrapf` return `nil` for a `nil` cause. Message format is
`op: message: cause` (empty parts are omitted).

## Sentinel errors (`errs/domain.go`)

Return these from the service layer for known business-rule failures:

- `ErrIncidentNotFound`, `ErrNoActiveIncident` — NOT_FOUND
- `ErrIncidentAlreadyActive`, `ErrIncidentAlreadyClosed` — CONFLICT

## Handling errors

Match a specific sentinel or a whole category:

```go
err := svc.CloseIncident(chatID)
switch {
case errors.Is(err, errs.ErrNoActiveIncident):
    return c.Send("No active incident in this chat.")
case errs.Is(err, errs.KindUnavailable):
    return c.Send("Service temporarily unavailable, try again later.")
case err != nil:
    log.Printf("close incident: %v", err)
    return c.Send("Internal error.")
}
```

`errs.KindOf(err)` returns the Kind of the first `*errs.Error` in the chain
(`KindInternal` for unclassified errors); `errs.Is(err, kind)` is a shortcut.

## Adding a new domain error

One line in `errs/domain.go` — no registration needed:

```go
var ErrSomethingWrong = New(KindConflict, "incident", "human-readable description")
```
