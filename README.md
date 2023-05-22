ExcelMate 应用用于处理 Excel 文件，包含 csv 格式，可使用多种简单的条件搜索数据，并将不同条件的 Excel 数据定制输出为自己想要的数据格式。


# ExcelMate 教程

## 界面
打开后右上角两个数据库，上面写条件，下面写模板。
![条件模板](./docs/1.png)


## 案例数据
假设 Excel 数据如下：

A | B | C | D
--- | --- | --- |----
时间 | 交易对象 | 金额 | 描述
2022-02-08 | 建设银行 | 20 | 我还款了
2022-03-08 | 招商银行 | -20 | 我借钱了


## 条件语法

### 列名

`A - Z` 的大写字母名称表示对应哪一列的数据。

### 数字列

如果某一列你认为是数字，则需要 `number(C)` 这样写，表示 C 列是数字。

### 大于 

`number(C) > 0` 查找 C 列大于 0 的数据。


### 小于

`number(C) < 0` 查找 C 列小于 0 的数据。

### 大于等于

`number(C) >= 0` 查找 C 列大于等于 0 的数据。


### 小于等于

`number(C) <= 0` 查找 C 列小于等于 0 的数据。



### 等于

` B == "建设银行" ` 查找 B 列为 “建设银行” 的数据。

### 不等于

`B != "建设银行"` 查找 B 列不是 “建设银行” 的数据。


### 模糊匹配

` D matches "还款"` 模糊匹配 D 列数据中包含 “还款” 字样的数据。

### 多条件满足

` B == "建设银行" and D matches "还款"` 使用 `and` 表示这两个条件都满足。

### 多条件有一个满足

` B == "建设银行" or D matches "还款"` 使用 `or` 表示这两个条件有一个满足就行。

## 模板语法
模板的作用是拿条件匹配的数据组装自己想要的最终内容。

`{{.A}}` 使用双括号、英文句号、列名组装一块，表示取出 A 列的数据。 
### Beancount
拿 Beancount 格式举个例子。
```
{{.A}} * "{{.B}}" "{{.D}}" 
    Expenses:Others              {{.G}} CNY
    Assets:Card:CCB
```

### 格式化时间
如果 Excel 中的时间不是自己想要的格式，有两种办法解决。

1. 使用 Excel 软件处理好。
2. 使用 模板语法。

展开说说模板语法:

* 数据 `2023-01-27 19:51:41`，目标 `2023-01-27`，模板语法 `{{toDate "2006-01-02 15:04:05" .A | date "2006-01-02"}}`

 

