# test_project_go
goose postgres "host=localhost port=5432 user=myuser password=mypassword dbname=test_db sslmode=disable" -dir migrations up
goose -dir migrations create courier sql