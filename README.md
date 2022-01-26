# 🐙 GithubGoAPI - API статистики пользоватя GitHub 📈

Семпл реквеста:
``` Elixir
GET https://hud0shnikgitapi.herokuapp.com/user/hud0shnik
```
``` Elixir
GET https://hud0shnikgitapi.herokuapp.com/user/hud0shnik/2022-01-20
```
Семпл респонса:
``` Json
{
"date":     "2022-01-21",
"username": "hud0shnik",
"commits":   9,
"color":     4
}
```
Параметр color - цвет ячейки. Всего есть 5 цветов: от серого (0) до ярко-зеленого(4)
