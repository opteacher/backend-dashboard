import axios from "axios"

import config from "../../config/backend"
import utils from "../utils"

export default {
    async post(model) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/insert/models`, model)
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async delete(id) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/delete/models`, {id})
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    }
}