drop table if exists guest_gift;

create table guest_gift
(
    guest_id   varchar(255) null comment '贵客图鉴ID',
    guest_name varchar(255) null comment '贵客名称',
    antique    varchar(255) null comment '礼物名',
    recipe     varchar(255) null comment '菜谱名'
)
    charset = utf8;

INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('001', '耶稣', '油火虫', '香酥鸭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('001', '耶稣', '暖石', '毛蟹年糕');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('001', '耶稣', '五星炒果', '肉蟹煲');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('001', '耶稣', '一昧真火', '瑞典肉丸');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('002', '如来', '一昧真火', '金缕虾');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('002', '如来', '千年煮鳖', '佛跳墙');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('002', '如来', '恐怖利刃', '孔雀开屏');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('002', '如来', '耐煮的水草', '菊花豆腐');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('003', '木良', '烤焦的菊花', '鸡蛋饼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('003', '木良', '蒸汽耳环', '泡椒凤爪');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('003', '木良', '暖石', '干炒牛河');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('003', '木良', '五香果', '虾酱空心菜');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('004', '木优', '刀嘴鹦鹉', '口水鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('004', '木优', '焦虫', '葱香烤鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('004', '木优', '油火虫', '土豆饼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('004', '木优', '一昧真火', '咖喱饺');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('005', '迷思吹法师', '蒸汽耳环', '葱油鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('005', '迷思吹法师', '油火虫', '牛肉煎包');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('005', '迷思吹法师', '防水的柠檬', '葱油拌面');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('005', '迷思吹法师', '刀嘴鹦鹉', '特色切糕');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('006', '大胃王', '油火虫', '锅贴');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('006', '大胃王', '焦虫', '烤肉拼盘');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('006', '大胃王', '一昧真火', '猪排饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('006', '大胃王', '蒸汽宝石', '梅菜扣肉');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('007', '屠夫老王', '防水的柠檬', '番茄肥牛');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('007', '屠夫老王', '暖石', '回锅肉');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('007', '屠夫老王', '油火虫', '椒盐排条');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('007', '屠夫老王', '香烤鱼排', '德式拼盘');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('008', '丐帮帮主', '焦虫', '葱香烤鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('008', '丐帮帮主', '蒸汽耳环', '白粽');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('008', '丐帮帮主', '刀嘴鹦鹉', '片皮鸭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('008', '丐帮帮主', '暖石', '湘西外婆菜');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('009', '潇洒李白', '一昧真火', '红烧排骨');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('009', '潇洒李白', '刀嘴鹦鹉', '醉蟹');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('009', '潇洒李白', '焦虫', '情人巧克力');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('009', '潇洒李白', '蒸馏杯', '清蒸武昌鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('010', '迷弟杜甫', '一昧真火', '红烧排骨');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('010', '迷弟杜甫', '防水的柠檬', '豚骨拉面');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('010', '迷弟杜甫', '刀嘴鹦鹉', '扣三丝');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('010', '迷弟杜甫', '千年煮鳖', '得莫利炖鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('011', '诸葛亮', '一昧真火', '猪排饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('011', '诸葛亮', '五香果', '扬州炒饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('011', '诸葛亮', '暖石', '炒鳝糊');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('011', '诸葛亮', '剪刀蟹', '番茄虾仁盅');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('012', '黄月英', '蒸汽耳环', '小笼包');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('012', '黄月英', '烤焦的菊花', '肉末茄子');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('012', '黄月英', '耐煮的水草', '萝卜小排汤');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('012', '黄月英', '耐煮的水草', '酥锅');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('013', '苏轼', '防水的柠檬', '蛤蜊豆腐汤');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('013', '苏轼', '耐煮的水草', '东坡肉');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('013', '苏轼', '一昧真火', '炸鱼排');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('013', '苏轼', '一昧真火', '麻辣小鱼干');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('014', '唐伯虎', '剪刀蟹', '糖番茄');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('014', '唐伯虎', '防水的柠檬', '酒酿小圆子');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('014', '唐伯虎', '鼓风机', '夏日风情堡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('014', '唐伯虎', '刀嘴鹦鹉', '三文鱼拼盘');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('015', '张飞', '油火虫', '牛肉煎包');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('015', '张飞', '耐煮的水草', '红烧牛肉面');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('015', '张飞', '暖石', '红咖喱鸭胸');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('016', '关羽', '防水的柠檬', '番茄肥牛');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('016', '关羽', '防水的柠檬', '红烧萝卜');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('016', '关羽', '耐煮的水草', '水煮鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('017', '刘备', '蒸汽耳环', '猪耳碟头');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('017', '刘备', '焦虫', '干锅香辣虾');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('017', '刘备', '香烤鱼排', '煎牛排');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('017', '刘备', '耐煮的水草', '三文鱼奶油汤');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('018', '吕布', '蒸汽耳环', '酿黄瓜');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('018', '吕布', '刀嘴鹦鹉', '罗宋汤');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('018', '吕布', '刀嘴鹦鹉', '豆豉蒸排骨');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('018', '吕布', '鼓风机', '风沙牛排');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('019', '貂蝉', '防水的柠檬', '皮蛋瘦肉粥');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('019', '貂蝉', '焦虫', '章鱼小丸子');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('019', '貂蝉', '耐煮的水草', '麻辣小龙虾');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('019', '貂蝉', '焦虫', '星空面包');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('020', '大内侍卫', '刀嘴鹦鹉', '罗宋汤');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('020', '大内侍卫', '防水的柠檬', '红烧萝卜');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('020', '大内侍卫', '五香果', '辣炒萝卜干');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('020', '大内侍卫', '千年煮鳖', '蟹黄鱼籽丸');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('021', '王子', '暖石', '咖喱土豆鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('021', '王子', '一昧真火', '菠萝咕咾肉');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('021', '王子', '焦虫', '肉酱千层面');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('021', '王子', '耐煮的水草', '豆腐脑');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('022', '公主', '暖石', '黄焖鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('022', '公主', '暖石', '菠萝饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('022', '公主', '蒸汽耳环', '水晶桃花');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('022', '公主', '蒸汽宝石', '柠檬蒸鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('023', '大明湖畔的美女', '油火虫', '臭豆腐');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('023', '大明湖畔的美女', '暖石', '爆炒猪肝');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('023', '大明湖畔的美女', '耐煮的水草', '海鲜面');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('023', '大明湖畔的美女', '焦虫', '白雪金狗');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('024', '皇帝', '蒸汽宝石', '虾饺皇');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('024', '皇帝', '五星炒果', '龙井虾仁');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('024', '皇帝', '香烤鱼排', '龙虾盛宴');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('024', '皇帝', '恐怖利刃', '战斧猪排');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('025', '皇后', '防水的柠檬', '老鸭粉丝汤');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('025', '皇后', '香烤鱼排', '蒜蓉芝士焗龙虾');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('025', '皇后', '耐煮的水草', '三文鱼茶泡饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('025', '皇后', '暖石', '甜筒披萨');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('026', '官爷', '油火虫', '家常豆腐');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('026', '官爷', '暖石', '洋葱小炒肉');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('026', '官爷', '恐怖利刃', '松鼠桂鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('026', '官爷', '鼓风机', '银河蛋糕');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('027', '富人', '蒸汽耳环', '蛤蜊蒸蛋');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('027', '富人', '焦虫', '香烤猪颈肉');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('027', '富人', '耐煮的水草', '腌笃鲜');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('027', '富人', '刀嘴鹦鹉', '海鲜甜甜圈');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('028', '富太太', '五香果', '麻婆豆腐');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('028', '富太太', '耐煮的水草', '鱼头豆腐汤');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('028', '富太太', '防水的柠檬', '海鲜粥');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('028', '富太太', '鼓风机', '腐衣黄鱼卷');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('029', '富二代', '蒸汽宝石', '蟹粉豆腐');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('029', '富二代', '耐煮的水草', '油焖笋');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('029', '富二代', '暖石', '宫保鸡丁');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('029', '富二代', '蒸馏杯', '炒面面包');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('030', '白富美', '防水的柠檬', '阳春面');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('030', '白富美', '暖石', '菠萝饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('030', '白富美', '刀嘴鹦鹉', '水晶肴蹄');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('030', '白富美', '剪刀蟹', '萝卜花花');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('048', '耶稣如来', '鼓风机', '九转大肠');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('048', '耶稣如来', '焦虫', '干锅鱼头');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('048', '耶稣如来', '千年煮鳖', '佛跳墙');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('048', '耶稣如来', '千年煮鳖', '金汤鱼翅娃娃菜');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('049', '木良木优', '五香果', '鸡肉蔬菜沙拉');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('049', '木良木优', '刀嘴鹦鹉', '鸡汤煮干丝');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('049', '木良木优', '蒸馏杯', '腊味合蒸');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('049', '木良木优', '暖石', '蘸酱菜');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('050', '潇洒李白迷弟杜甫', '刀嘴鹦鹉', '撒尿牛丸');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('050', '潇洒李白迷弟杜甫', '耐煮的水草', '大骨高汤');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('050', '潇洒李白迷弟杜甫', '恐怖利刃', '蟹黄鱼翅');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('050', '潇洒李白迷弟杜甫', '蒸汽宝石', '蘑菇蒸鳕鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('051', '诸葛亮黄月英', '刀嘴鹦鹉', '夫妻肺片');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('051', '诸葛亮黄月英', '香烤鱼排', '蒜蓉芝士焗龙虾');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('051', '诸葛亮黄月英', '一昧真火', '天妇罗拼盘');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('051', '诸葛亮黄月英', '一昧真火', '风味塔可');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('052', '张飞关羽刘备', '暖石', '三杯鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('052', '张飞关羽刘备', '香烤鱼排', '咖喱香辣虾');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('052', '张飞关羽刘备', '耐煮的水草', '鸡肉火锅');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('053', '吕布貂蝉', '蒸汽宝石', '蟹粉狮子头');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('053', '吕布貂蝉', '刀嘴鹦鹉', '金针菇牛肉卷');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('053', '吕布貂蝉', '恐怖利刃', '松鼠桂鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('053', '吕布貂蝉', '恐怖利刃', '天女散花');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('031', 'TapTap', '神秘魔方', '鸡清汤青拉面');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('032', '圣诞老人', '香烤鱼排', '仰望星空');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('032', '圣诞老人', '焦虫', '牛肉披萨');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('032', '圣诞老人', '香烤鱼排', '圣诞烤鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('033', '屈原', '蒸汽耳环', '蛋黄粽');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('033', '屈原', '耐煮的水草', '雄黄酒');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('033', '屈原', '恐怖利刃', '刺身拼盘');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('033', '屈原', '耐煮的水草', '海南椰子鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('034', '王母娘娘', '香烤鱼排', '绿咖喱鸡肉卷');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('034', '王母娘娘', '耐煮的水草', '血腥玛丽');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('034', '王母娘娘', '刀嘴鹦鹉', '日出鸡尾酒');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('035', '武松', '五星炒果', '香辣蟹');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('035', '武松', '一昧真火', '芝士猪肉卷');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('035', '武松', '刀嘴鹦鹉', '杏鲍菇牛肉粒');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('035', '武松', '暖石', '大盘鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('036', '伯爵大人', '鼓风机', '邪神烤鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('036', '伯爵大人', '暖石', '幽灵咖喱');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('036', '伯爵大人', '一昧真火', '木乃伊香肠');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('036', '伯爵大人', '蒸汽宝石', '恶作剧便当');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('037', '西瓜小哥', '鼓风机', '冷锅鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('037', '西瓜小哥', '一昧真火', '爱心宽油蛋');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('038', '福神', '耐煮的水草', '柠檬鸡丝');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('038', '福神', '蒸汽宝石', '雪笋蒸干贝');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('038', '福神', '五星炒果', '蛋白蟹肉');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('038', '福神', '香烤鱼排', '豪华海鲜桶');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('039', '禄神', '一昧真火', '富贵双方');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('039', '禄神', '耐煮的水草', '鸡肉火锅');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('039', '禄神', '耐煮的水草', '酸菜鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('039', '禄神', '蒸汽宝石', '黄陂三合');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('040', '寿神', '一昧真火', '灯笼茄子');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('040', '寿神', '蒸汽宝石', '竹筒饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('040', '寿神', '蒸馏杯', '雪花鱼糕');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('054', '福禄寿三神', '焦虫', '鱿鱼筒烤饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('054', '福禄寿三神', '千年煮鳖', '鱼跃龙门');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('054', '福禄寿三神', '蒸馏杯', '南瓜鱼翅盅');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('041', '罗爸爸', '耐煮的水草', '爷爷泡的茶');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('041', '罗爸爸', '千年煮鳖', '灵芝空间');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('042', '罗妈妈', '刀嘴鹦鹉', '冰激凌吐司');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('042', '罗妈妈', '刀嘴鹦鹉', '抹茶舒芙蕾');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('043', '黑道小三', '耐煮的水草', '海胆盖饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('043', '黑道小三', '蒸汽宝石', '土瓶蒸');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('044', '黑道小伊', '蒸汽耳环', '蛋黄酱盖饭');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('044', '黑道小伊', '焦虫', '烤火鸡腿');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('055', '不良青年', '千年煮鳖', '痛风锅');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('055', '不良青年', '一昧真火', '赏樱便当');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('045', '小宝', '焦虫', '京八件');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('045', '小宝', '蒸馏杯', '汽锅鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('045', '小宝', '耐煮的水草', '金玉田鸡');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('046', '容嬷嬷', '油火虫', '炸佛手');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('046', '容嬷嬷', '蒸汽宝石', '蒸汽石锅鱼');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('047', '太上皇', '千年煮鳖', '御凤还巢');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('047', '太上皇', '焦虫', '烤乳猪');
INSERT INTO guest_gift (guest_id, guest_name, antique, recipe) VALUES ('047', '太上皇', '蒸馏杯', '人参果');