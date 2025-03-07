module github.com/danielkhtse/supreme-adventure/transaction-service

go 1.24.1

require (
	github.com/danielkhtse/supreme-adventure/account-service v0.0.0
	github.com/danielkhtse/supreme-adventure/common v0.0.0
	gorm.io/gorm v1.25.10
)

replace (
	github.com/danielkhtse/supreme-adventure/account-service => ../account-service
	github.com/danielkhtse/supreme-adventure/common => ../common
)
