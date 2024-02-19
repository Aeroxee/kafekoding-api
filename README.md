# KafeKoding API

KafeKoding API is project to make RestFul API to handling account, class, blog of articles in KafeKoding Community.

## Requirement

1. [GIN](https://github.com/gin-gonic/gin)
2. [UUID](https://github.com/google/uuid)
3. [GoSlug](https://github.com/gosimple/slug)
4. [GORM](https://gorm.io)

## How To Run

1. Clone this repository
```
git clone https://github.com/Aeroxee/kafekoding-api.git
```
2. Create `.env`
```
SMTP_SERVER=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your email
SMTP_PASSWORD=your app password
```
3. Run & execution
```
go build -o kafekodingapi cmd/kafekoding-api/main.go
./kafekodingpi
```
4. Server run on `:8000`

## Contribution

If you encounter errors, have questions, or would like to contribute to the development of this class, please open an issue or submit a pull request. Your contribution is greatly appreciated!

## License

This project is licensed under the [MIT license](). Please refer to the [LICENSE](LICENSE) file for details.