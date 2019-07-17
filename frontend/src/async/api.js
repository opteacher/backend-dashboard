import axios from "axios"

import config from "../../config/backend"
import utils from "../utils"

export default {
    async qryAll() {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/apis.selectAll`)
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qryByName(name) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/apis.selectByName`, {name})
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async add(api) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/apis.insert`, api)
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    }
}