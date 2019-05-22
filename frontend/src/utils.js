export default {
    getErrorMsg(res) {
        if (res instanceof String) {
            return res
        } else if (res.message) {
            return res.message
        } else {
            return null
        }
    }
}