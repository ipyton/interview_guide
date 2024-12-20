/* global use, db */
// MongoDB Playground
// To disable this template go to Settings | MongoDB | Use Default Template For Playground.
// Make sure you are connected to enable completions and to be able to run a playground.
// Use Ctrl+Space inside a snippet or a string literal to trigger completions.
// The result of the last command run in a playground is shown on the results panel.
// By default the first 20 documents will be returned with a cursor.
// Use 'console.log()' to print to the debug output.
// For more documentation on playgrounds please refer to
// https://www.mongodb.com/docs/mongodb-vscode/playgrounds/

// Select the database to use.
// rs.initiate()


categories = [{
    "code": "1",
    "name": "农、林、牧、渔业",
    "child": [{
        "code": "01",
        "name": "农业"
    },
    {
        "code": "02",
        "name": "林业"
    },
    {
        "code": "03",
        "name": "畜牧业"
    },
    {
        "code": "04",
        "name": "渔业"
    },
    {
        "code": "05",
        "name": "农、林、牧、渔服务业"
    }
    ]

},
{
    "code": "2",
    "name": "采矿业",
    "child": [{
        "code": "06",
        "name": "煤炭开采和洗选业"
    },
    {
        "code": "07",
        "name": "石油和天然气开采业"
    },
    {
        "code": "08",
        "name": "黑色金属矿采选业"
    },
    {
        "code": "09",
        "name": "有色金属矿采选业"
    },
    {
        "code": "10",
        "name": "非金属矿采选业"
    },
    {
        "code": "11",
        "name": "开采辅助活动"
    },
    {
        "code": "12",
        "name": "其他采矿业"
    }
    ]
},
{
    "code": "3",
    "name": "制造业",
    "child": [{
        "code": "13",
        "name": "农副食品加工业"
    },
    {
        "code": "14",
        "name": "食品制造业"
    },
    {
        "code": "15",
        "name": "酒、饮料和精制茶制造业"
    },
    {
        "code": "16",
        "name": "烟草制品业"
    },
    {
        "code": "17",
        "name": "纺织业"
    },
    {
        "code": "18",
        "name": "纺织服装、服饰业"
    },
    {
        "code": "19",
        "name": "皮革、毛皮、羽毛及其制品和制鞋业"
    },
    {
        "code": "20",
        "name": "木材加工和木、竹、藤、棕、草制品业"
    },
    {
        "code": "21",
        "name": "家具制造业"
    },
    {
        "code": "22",
        "name": "造纸和纸制品业"
    },
    {
        "code": "23",
        "name": "印刷和记录媒介复制业"
    },
    {
        "code": "24",
        "name": "文教、工美、体育和娱乐用品制造业"
    },
    {
        "code": "25",
        "name": "石油加工、炼焦和核燃料加工业"
    },
    {
        "code": "26",
        "name": "化学原料和化学制品制造业"
    },
    {
        "code": "27",
        "name": "医药制造业"
    },
    {
        "code": "28",
        "name": "化学纤维制造业"
    },
    {
        "code": "29",
        "name": "橡胶和塑料制品业"
    },
    {
        "code": "30",
        "name": "非金属矿物制品业"
    },
    {
        "code": "31",
        "name": "黑色金属冶炼和压延加工业"
    },
    {
        "code": "32",
        "name": "有色金属冶炼和压延加工业"
    },
    {
        "code": "33",
        "name": "金属制品业"
    },
    {
        "code": "34",
        "name": "通用设备制造业"
    },
    {
        "code": "35",
        "name": "专用设备制造业"
    },
    {
        "code": "36",
        "name": "汽车制造业"
    },
    {
        "code": "37",
        "name": "铁路、船舶、航空航天和其他运输设备制造业"
    },
    {
        "code": "38",
        "name": "电气机械和器材制造业"
    },
    {
        "code": "39",
        "name": "计算机、通信和其他电子设备制造业"
    },
    {
        "code": "40",
        "name": "仪器仪表制造业"
    },
    {
        "code": "41",
        "name": "其他制造业"
    },
    {
        "code": "42",
        "name": "废弃资源综合利用业"
    },
    {
        "code": "43",
        "name": "金属制品、机械和设备修理业"
    }
    ]
},
{
    "code": "4",
    "name": "电力、热力、燃气及水生产和供应业",
    "child": [{
        "code": "44",
        "name": "电力、热力生产和供应业"
    },
    {
        "code": "45",
        "name": "燃气生产和供应业"
    },
    {
        "code": "46",
        "name": "水的生产和供应业"
    }
    ]
},
{
    "code": "5",
    "name": "建筑业",
    "child": [{
        "code": "47",
        "name": "房屋建筑业"
    },
    {
        "code": "48",
        "name": "土木工程建筑业"
    },
    {
        "code": "49",
        "name": "建筑安装业"
    },
    {
        "code": "50",
        "name": "建筑装饰和其他建筑业"
    }
    ]
},
{
    "code": "6",
    "name": "批发和零售业",
    "child": [{
        "code": "51",
        "name": "批发业"
    },
    {
        "code": "52",
        "name": "零售业"
    }
    ]
},
{
    "code": "7",
    "name": "交通运输、仓储和邮政业",
    "child": [{
        "code": "53",
        "name": "铁路运输业"
    },
    {
        "code": "54",
        "name": "道路运输业"
    },
    {
        "code": "55",
        "name": "水上运输业"
    },
    {
        "code": "56",
        "name": "航空运输业"
    },
    {
        "code": "57",
        "name": "管道运输业"
    },
    {
        "code": "58",
        "name": "装卸搬运和运输代理业"
    },
    {
        "code": "59",
        "name": "仓储业"
    },
    {
        "code": "60",
        "name": "邮政业"
    }
    ]
},
{
    "code": "8",
    "name": "住宿和餐饮业",
    "child": [{
        "code": "61",
        "name": "住宿业"
    },
    {
        "code": "62",
        "name": "餐饮业"
    }
    ]
},
{
    "code": "9",
    "name": "信息传输、软件和信息技术服务业",
    "child": [{
        "code": "63",
        "name": "电信、广播电视和卫星传输服务"
    },
    {
        "code": "64",
        "name": "互联网和相关服务"
    },
    {
        "code": "65",
        "name": "软件和信息技术服务业"
    }
    ]
},
{
    "code": "10",
    "name": "金融业",
    "child": [{
        "code": "66",
        "name": "货币金融服务"
    },
    {
        "code": "67",
        "name": "资本市场服务"
    },
    {
        "code": "68",
        "name": "保险业"
    },
    {
        "code": "69",
        "name": "其他金融业"
    }
    ]
},
{
    "code": "11",
    "name": "房地产业",
    "child": [{
        "code": "70",
        "name": "房地产业"
    }]
},
{
    "code": "12",
    "name": "租赁和商务服务业",
    "child": [{
        "code": "71",
        "name": "租赁业"
    },
    {
        "code": "72",
        "name": "商务服务业"
    }
    ]
},
{
    "code": "13",
    "name": "科学研究和技术服务业",
    "child": [{
        "code": "73",
        "name": "研究和试验发展"
    },
    {
        "code": "74",
        "name": "专业技术服务业"
    },
    {
        "code": "75",
        "name": "科技推广和应用服务业"
    }
    ]
},
{
    "code": "14",
    "name": "水利、环境和公共设施管理业",
    "child": [{
        "code": "76",
        "name": "水利管理业"
    },
    {
        "code": "77",
        "name": "生态保护和环境治理业"
    },
    {
        "code": "78",
        "name": "公共设施管理业"
    }
    ]
},
{
    "code": "15",
    "name": "居民服务、修理和其他服务业",
    "child": [{
        "code": "79",
        "name": "居民服务业"
    },
    {
        "code": "80",
        "name": "机动车、电子产品和日用产品修理业"
    },
    {
        "code": "81",
        "name": "其他服务业"
    }
    ]
},
{
    "code": "16",
    "name": "教育",
    "child": [{
        "code": "82",
        "name": "教育"
    }]
},
{
    "code": "17",
    "name": "卫生和社会工作",
    "child": [{
        "code": "83",
        "name": "卫生"
    },
    {
        "code": "84",
        "name": "社会工作"
    }
    ]
},
{
    "code": "18",
    "name": "文化、体育和娱乐业",
    "child": [{
        "code": "85",
        "name": "新闻和出版业"
    },
    {
        "code": "86",
        "name": "广播、电视、电影和影视录音制作业"
    },
    {
        "code": "87",
        "name": "文化艺术业"
    },
    {
        "code": "88",
        "name": "体育"
    },
    {
        "code": "89",
        "name": "娱乐业"
    }
    ]
},
{
    "code": "19",
    "name": "公共管理、社会保障和社会组织",
    "child": [{
        "code": "90",
        "name": "中国共产党机关"
    },
    {
        "code": "91",
        "name": "国家机构"
    },
    {
        "code": "92",
        "name": "人民政协、民主党派"
    },
    {
        "code": "93",
        "name": "社会保障"
    },
    {
        "code": "94",
        "name": "群众团体、社会团体和其他成员组织"
    },
    {
        "code": "95",
        "name": "基层群众自治组织"
    }
    ]
},
{
    "code": "20",
    "name": "国际组织",
    "child": [{
        "code": "96",
        "name": "国际组织"
    }]
}
]
let id = 0
function traverse(root) {
    let father_node
    let queue = []
    let result = []
    result.push({ parent_class_id: -1, class_name: root.name,  class_id: id, isLeaf: !root.child})
    root.class_id = id
    id ++
    queue.push(root)
    while (queue.length !== 0) {
        let father_node = queue.shift()
        if (father_node.child) {
            father_node.child.forEach(element => {
                result.push({ parent_class_id: father_node.class_id, class_id: id, class_name: element.name, isLeaf: !element.child})
                element.class_id = id
                id ++
                queue.push(element)
            })
        }
    }
    return result

}
let result = []
categories.forEach(element => {
    result = [...result, ...traverse(element)]
});
print(result)

use('interview_guide');

// Insert a few documents into the sales collection.
db.getCollection('classes').insertMany(result);


// Create a new index in the collection.
db.getCollection('classes')
    .createIndex(
        {
            parent_class_id: 1
        }, {
        unique: false
    }
    );


