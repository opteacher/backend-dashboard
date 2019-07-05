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
    }
}