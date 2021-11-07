## Project Layout

The starter kit uses the following project layout:

```
.
├── cmd                  main applications of the project
│   └── server           the API server application
├── config               configuration files for different environments
├── internal             private application and library code
│   ├── config           configuration library
│   ├── deposit          deposit-related features
│   ├── entity           database models
│   ├── errors           error types and handling
│   ├── test             helpers for testing purpose
│   └── transaction      transaction-related features
├── pkg                  public library code
│   ├── accesslog        access log middleware
│   ├── log              structured and context-aware logger
└── testdata             test data scripts
```

Этот проект основан на [Go RESTful API Starter Kit](https://github.com/qiangxue/go-rest-api).
Огромная благодарность [@qiangxue](https://github.com/qiangxue) за его работу.