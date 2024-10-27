@echo off

for /f "delims=\" %%a in ('dir /b /a-d /on ".\sql\*.sql"') do (
	echo %%a
	docker exec mysql /usr/bin/mysql -uroot -p111111 -e "source /sql/%%a"
)

pause