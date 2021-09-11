create table laboratory
(
    target_name    varchar(255) null,
    target_type    varchar(255) null,
    rarity         int          null,
    skill          varchar(255) null,
    antique        varchar(255) null,
    antique_number int          null,
    equips         text         null,
    recipes        text         null
)
    charset = utf8;

INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('猪肉脯', 'recipe', 1, '炒', '五香果', 4, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('豆腐青菜', 'recipe', 1, '炒', '五香果', 4, '[]', '["炒青菜","香菇青菜","青菜豆腐汤"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('青椒香干', 'recipe', 1, '炒', '五香果', 4, '[]', '["剁椒跑蛋","罗宋汤"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('锅巴小龙虾', 'recipe', 3, '炒', '五香果', 2, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('九珍海鲜', 'recipe', 3, '炒', '五香果', 2, '[]', '["酱爆鱿鱼","金针菇牛肉卷"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('菠萝饭', 'recipe', 3, '炒', '五香果', 2, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('毛蟹年糕', 'recipe', 4, '炒', '暖石', 3, '["银锅铲"]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('蛋白蟹肉', 'recipe', 5, '炒', '五星炒果', 1, '["银锅铲","银制剪刀"]', '["咖喱饺","天妇罗拼盘","撒尿牛丸","蟹粉狮子头"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('烤蔬菜沙拉', 'recipe', 2, '烤', '烤焦的菊花', 5, '[]', '["浇汁菜心"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('鱿鱼筒烤饭', 'recipe', 4, '烤', '焦虫', 3, '["银制酒提"]', '["蟹粉豆腐","金缕虾","竹筒饭"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('红酒鸭胸', 'recipe', 4, '烤', '焦虫', 3, '["精致面包机"]', '["香酥鸭","片皮鸭","红油抄手"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('辣条', 'recipe', 1, '切', '剪刀蟹', 4, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('柠檬鸡丝', 'recipe', 3, '切', '剪刀蟹', 2, '[]', '["酸辣鸡杂","鱼香肉丝"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('特色切糕', 'recipe', 4, '切', '刀嘴鹦鹉', 3, '["红木蒸饭桶"]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('战斧猪排', 'recipe', 5, '切', '恐怖利刃', 1, '["银制剪刀","精致面包机"]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('素丸子', 'recipe', 2, '炸', '油火虫', 5, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('脆皮烧肉', 'recipe', 2, '炸', '油火虫', 5, '[]', '["蒜泥白肉"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('黑椒牛仔骨', 'recipe', 4, '炸', '一昧真火', 3, '["白猫筷子架"]', '["京酱肉丝","生拌牛肉","煎牛排"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('黄瓜鸡蛋卷', 'recipe', 2, '蒸', '蒸汽耳环', 5, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('菜饭', 'recipe', 2, '蒸', '蒸汽耳环', 5, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('八宝饭', 'recipe', 3, '蒸', '蒸汽耳环', 2, '[]', '["菠萝饭","猪排饭","鳗鱼炒饭"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('猪血糕', 'recipe', 3, '蒸', '蒸汽耳环', 2, '[]', '["酒酿小圆子","咖喱饺","虾饺皇"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('三色凉糕', 'recipe', 3, '蒸', '蒸汽耳环', 2, '[]', '["牛肉卷饼","南瓜饼","土豆饼"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('荷花豆腐', 'recipe', 3, '蒸', '蒸汽耳环', 2, '[]', '["豆豉蒸鱼","海鲜面","鱼头豆腐汤"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('三鲜砂锅', 'recipe', 3, '蒸', '蒸汽耳环', 2, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('儿童餐', 'recipe', 4, '蒸', '蒸汽宝石', 3, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('雪花鱼糕', 'recipe', 5, '蒸', '蒸馏杯', 1, '["红木蒸饭桶","银制剪刀"]', '["松鼠桂鱼","蒜蓉芝士焗龙虾"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('海鲜粥', 'recipe', 2, '煮', '防水的柠檬', 5, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('酥锅', 'recipe', 3, '煮', '防水的柠檬', 2, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('老鸭笋尖汤', 'recipe', 3, '煮', '防水的柠檬', 2, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('汤圆', 'recipe', 3, '煮', '防水的柠檬', 2, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('南瓜粥', 'recipe', 3, '煮', '防水的柠檬', 2, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('鸡肉火锅', 'recipe', 4, '煮', '耐煮的水草', 3, '["银水壶"]', '["咖喱土豆牛肉","辣子鸡","黄焖鸡","三杯鸡"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('猪肚鸡', 'recipe', 4, '煮', '耐煮的水草', 3, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('蟹黄鱼翅', 'recipe', 5, '煮', '千年煮鳖', 1, '["银水壶","铜水壶"]', '["咖喱香辣虾","雪笋蒸干贝","鳗鱼炒饭","生滚鱼片粥"]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('金汤鱼翅娃娃菜', 'recipe', 5, '煮', '千年煮鳖', 1, '["精致煎锅铲","银水壶"]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('花果茶', 'recipe', 2, '煮', '防水的柠檬', 5, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('豆皮寿司', 'recipe', 2, '切', '剪刀蟹', 5, '[]', '[]');
INSERT INTO laboratory (target_name, target_type, rarity, skill, antique, antique_number, equips, recipes) VALUES ('爱心宽油蛋', 'recipe', 3, '炸', '油火虫', 2, '[]', '[]');