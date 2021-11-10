# Перевод денег между пользователями

Перевести деньги со счета одного пользователя на счет другого пользователя. Если счета получателя не существует, то он 
будет создан.

**URL** : `/v1/deposits/transfer`

**Метод** : `POST`

**Формат запроса**

Деньги будут списаны со счета пользователя с ID равным `sender_id` и зачислены на счет пользователя с ID равным 
`recipient_id`.

```json
{
  "sender_id"   : "[строка, UUID]",
  "recipient_id": "[строка, UUID]",
  "amount"      : "[число, положительное]",
  "description" : "[строка, опционально, до 100 символов]"
}
```

**Пример запроса**

```json
{
  "sender_id": "8c5593a0-37d3-11ec-8d3d-0242ac130001",
  "recipient_id": "6e726185-586e-49a7-89a4-6cfc2b03b0a2",
  "amount" : 300,
  "description": "happy birthday!"
}
```

## Ответ - успех

**Код** : `200 OK`

**Пример ответа**: транзакция, отражающая соответсвующий перевод.

```json
{
  "id": 8,
  "sender_id": "8c5593a0-37d3-11ec-8d3d-0242ac130001",
  "recipient_id": "6e726185-586e-49a7-89a4-6cfc2b03b0a2",
  "amount": 300,
  "description": "happy birthday!",
  "transaction_date": "2021-11-10T14:24:17.4145906Z"
}
```

## Ответ - ошибка

**Причина** : Параметры запроса некорректны

**Код** : `400 BAD REQUEST`

**Пример ответа** :

```json
{
  "status": 400,
  "message": "There is some problem with the data you submitted.",
  "details": [
    {
      "field": "amount",
      "error": "must be greater than 0"
    },
    {
      "field": "sender_id",
      "error": "must be a valid UUID"
    }
  ]
}
```

### ИЛИ

**Причина** : Отправитель не имеет достаточно средств для совершения перевода

**Код** : `403 FORBIDDEN`

**Пример ответа**

```json
{
  "status": 403,
  "message": "Insufficient funds to perform operation."
}
```