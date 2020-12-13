# 查菜谱

## 1. 基础查询

### 1.1 通过菜谱名查询

查询方式：`#菜谱 菜谱名`

可输入全称查询或输入关键字进行模糊查询，有唯一匹配的菜谱则会直接返回结果

示例：`#菜谱 鱼生` `#菜谱 顺德鱼生`

![]()

### 1.2 通过菜谱编号查询

查询方式：`#菜谱 菜谱编号`

有多个符合查询条件的结果时会返回一个列表，此时可以通过菜谱编号进行查询

示例：`#菜谱 015`

![]()

## 2. 复合查询

复合查询功能可以根据指定的条件得到一个菜谱列表，当只有一个菜谱符合要求时会返回该菜谱的信息。

**复合查询格式为 `#菜谱 参数列表`，参数列表的每个参数之间使用空格分隔，前后有关联的一组参数使用减号 `-` 连接。**

如: `#菜谱 食材-茄子 金币效率 p2`

![]()

> 复合查询一般结合任务需求进行查询，胡乱组合的查询结果往往没有实际意义哦

目前有以下三类查询参数：

### 2.1 普通筛选参数

筛选技法：`炒技法`

筛选调料类型：`甜味`

设置单价下限：`$100`

设置稀有度下限：`3火` / `3星`

有较多数据时翻页：`p2`

### 2.2 需要说明参数类型的参数

> 此类参数需要使用 `-` 连接所有条件

筛选食材：`食材-茄子`，`食材-螃蟹-贝壳`，`食材-肉类`， `食材-蔬菜`

筛选技法：`技法-炒`，`技法-煮-切`

筛选调料类型：`调料-甜`，`调料-咸`

筛选贵客：`贵客-丐帮帮主`，`稀有客人-大胃王`

筛选符文：`符文-剪刀蟹`

筛选来源：`来源-不夜城池`，`来源-实验室`

### 2.3 排序参数

目前提供以下排序方式，默认按照稀有度排序：

`图鉴序`、`单时间`、`总时间`、`单价`、`金币效率`、`耗材效率`、`稀有度`

### 复合查询示例

> 上面的各种参数类型也许看的有点晕，别急，以下会结合一些具体任务来示例如何使用复合查询，理解之后就很容易啦~

任务：`开业中使用任意食材4260个`。想要快速消耗食材，根据耗材效率排序菜谱。

`#菜谱 耗材效率`

![]()

任务：`累计开业获得200000金币`。想根据金币效率排序菜谱。第一页似乎没有想找的菜，那就翻到第二页.

`#菜谱 金币效率 p2`

![]()

任务：`制作使用鸡胸肉的料理60份`。想查询使用了鸡胸肉的菜谱，并根据单时间排序。

`#菜谱 食材-鸡胸肉 单时间`

![]()

任务：`制作200份使用炸技的料理且菜谱稀有度至少为4星`。想查询4星以上使用了炸技法的菜谱，并根据单时间排序。

`#菜谱 技法-炸 4星 单时间` 或者 `#菜谱 炸技法 四火 单时间`

![]()

任务：`制作单价不低于300金币的炒技料理150份`。想查询单价在300以上并且使用了炒技法的菜谱，并根据单时间排序。

`#菜谱 炒技法 $300 单时间` 或者 `#菜谱 技法-炒 $300 单时间`

![]()

任务：`制作3星以上使用鸡蛋和面粉的料理各80份`。想查询有没有哪个三星以上的菜同时使用了鸡蛋和面粉，并根据单时间排序。

`#菜谱 食材-鸡蛋-面粉 3星 单时间`

![]()

任务1：`制作使用水果的料理80份` 任务2：`制作3星切技法的料理80份`。想要同时做这两个任务，就想查询3星以上用了水果的切技法料理，并根据单时间排序。

`#菜谱 食材-水果 切技法 3星 单时间` 或者 `#菜谱 食材-水果 技法-切 3火 单时间 `

![]()

> 以下再举例一些其它场景的用法

想要查看不夜城池有哪些菜谱。

`#菜谱 来源-不夜城池`

![]()

想要查看哪些菜谱有概率出恐怖利刃，并根据每组总时间排序。

`#菜谱 符文-恐怖利刃 总时间`

![]()

