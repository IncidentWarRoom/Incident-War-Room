# Демо: инцидент из алерта Prometheus

Показывает цепочку: **кнопка «Сломать» → метрика в Prometheus → алерт → Alertmanager →
вебхук → инцидент** в Telegram-форуме и на дашборде.

```
demo-app (/metrics) ──scrape──▶ Prometheus ──alert firing──▶ Alertmanager
                                                                   │ webhook
                                                                   ▼
                                              incident-service /webhooks/alertmanager
                                                                   ▼
                                       инцидент: тема в Telegram (ALERT_CHAT_ID) + фронт
```

## Что где крутится

| Сервис | URL | Роль |
|---|---|---|
| demo-app | http://localhost:9000 | страница с кнопками 🔥 Сломать / ✅ Починить + метрика `demo_app_error_rate` |
| Prometheus | http://localhost:9090/alerts | скрейпит demo-app, считает алерт `HighErrorRate` |
| Alertmanager | http://localhost:9093 | шлёт вебхук в incident-service |
| incident-service | http://localhost:8080 | принимает вебхук, создаёт инцидент |
| frontend | http://localhost:3000 | список инцидентов |

## 0. Разовая настройка (уже сделана)

В `.env` должен быть id forum-супергруппы, где бот — админ с правом **Manage Topics**:
```dotenv
ALERT_CHAT_ID=-1004337045780
```
Если меняешь группу — как достать новый id: останови бота (`docker compose stop incident-service`),
напиши в группе `/start@имя_бота`, выполни
`curl -s "https://api.telegram.org/bot<BOT_TOKEN>/getUpdates"` и возьми `chat.id` (начинается с `-100`,
`is_forum: true`), затем верни бота (`docker compose up -d incident-service`).
> Бота нельзя опрашивать через `getUpdates`, пока `incident-service` запущен — будет `409 Conflict`.

## 1. Поднять стек

```bash
docker compose up --build -d
docker compose ps        # все 7 сервисов Up; postgres/report-service — healthy
```

## 2. Запустить демо (на камеру)

1. Открой вкладки: фронт `http://localhost:3000`, demo `http://localhost:9000`,
   Prometheus `http://localhost:9090/alerts`.
2. Покажи спокойный список инцидентов на фронте.
3. На demo-странице нажми **🔥 Сломать** (или `curl -XPOST localhost:9000/break`).
4. В Prometheus `/alerts`: `HighErrorRate` → `Pending` → через ~15с `Firing`.
5. Через ~5–10с после `Firing` в Telegram-форуме появляется новая тема инцидента
   **HIGH «High error rate on demo-app»**, он же — в списке на фронте (обнови страницу).

Итого от нажатия до инцидента ~20–30с.

## 3. Сбросить / повторить

```bash
curl -s -XPOST localhost:9000/fix      # вернуть demo-app в «здоров» (или кнопка ✅ Починить)
```
- «Починить» гасит алерт в Prometheus; **созданный инцидент остаётся** — закрытие инцидента ручное,
  командой бота в Telegram (`resolved`-вебхуки сервис игнорирует by design).
- Каждый цикл Сломать→Починить→Сломать создаёт **новый** инцидент, так что демо можно повторять.

## 4. Проверки / отладка

```bash
# метрика: 0 в норме, 1 после «Сломать»
curl -s localhost:9000/metrics | grep '^demo_app_error_rate'

# таргет demo-app должен быть UP
curl -s localhost:9090/api/v1/targets | grep -o '"health":"[^"]*"'

# состояние алерта
curl -s localhost:9090/api/v1/alerts | grep -o '"state":"[^"]*"'

# список инцидентов
curl -s localhost:8080/api/v1/incidents

# что прилетело в incident-service
docker compose logs incident-service --since 2m | grep -i alert
```

Быстрый прогон без ожидания Prometheus (дёрнуть вебхук напрямую):
```bash
curl -i -XPOST localhost:8080/webhooks/alertmanager -H 'Content-Type: application/json' \
  -d '{"alerts":[{"status":"firing","labels":{"alertname":"HighErrorRate","severity":"critical"},"annotations":{"summary":"High error rate on demo-app"}}]}'
```

Частые проблемы:
- В логах `alert chat is not configured` → не задан/не тот `ALERT_CHAT_ID` (см. шаг 0), перезапусти `incident-service`.
- Инцидент не пришёл в Telegram, но в API есть → бот не админ в группе или нет права Manage Topics.
- Алерт не доходит до `firing` → проверь, что demo-app `UP` в `/targets` и метрика `= 1`.

## 5. Остановить

```bash
docker compose stop              # пауза, данные сохраняются
# или
docker compose down              # снести контейнеры (том postgres сохраняется)
```
