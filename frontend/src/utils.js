export default {
    getErrorMsg(res) {
        if (typeof res === "string") {
            return res
        } else if (res.message) {
            return res.message
        } else {
            return null
        }
    }
}