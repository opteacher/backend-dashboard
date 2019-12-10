export default {
    getErrorMsg(res) {
        if (typeof res === "string") {
            return res
        } else if (res.message) {
            return res.message
        } else {
            return null
        }
    },
    between(n, a, b) {
        return n <= Math.max(a, b) && n >= Math.min(a, b)
    },
    supportTypes: [{
        title: "文本",
        value: "string"
    }, {
        title: "数字",
        value: "int32"
    }, {
        title: "日期",
        value: "uint64"
    }, {
        title: "布尔",
        value: "bool"
    }]
}