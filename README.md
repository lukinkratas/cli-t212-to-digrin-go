# CLI T212 to Digrin
Golang CLI tool for fetching T212 reports and transforming them to be used in Digrin portfolio tracker. Stores the reports in S3.

```
echo "T212_API_KEY=$T212_API_KEY" >> .env
echo "AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID" >> .env
echo "AWS_SECRET_ACCESS_KEY=$AWS_SECRET_ACCESS_KEY" >> .env
```

```
go mod tidy
```

```
go run main.py
```

# TODO

- [ ] investigate option of go routines

- [ ] add logging

- [ ] add decorators or decorators like functionality
