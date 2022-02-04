# 🐙 API статистики пользоватя GitHub 📈

<h3>Семпл реквеста </h3>

``` Elixir
GET https://hud0shnikgitapi.herokuapp.com/user/hud0shnik
```
``` Elixir
GET https://hud0shnikgitapi.herokuapp.com/user/hud0shnik/2022-01-20
```
<h3>Семпл респонса </h3>

``` Json
{
"date":     "2022-01-21",
"username": "hud0shnik",
"commits":   9,
"color":     4
}
```
> Параметр color - цвет ячейки. Всего есть 5 цветов: от серого (0) до ярко-зеленого(4)
