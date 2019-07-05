import _ from "lodash"
import axios from "axios"

import config from "../../config/backend"
import utils from "../utils"

export default {
    async add(link) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/links.insert`, _.omit(link, "id"))
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qry() {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/links.selectAll`)
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    }
}