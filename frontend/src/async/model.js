import _ from "lodash"
import axios from "axios"

import config from "../../config/backend"
import utils from "../utils"

export default {
    async post(model) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/models`, {
                method: "insert",
                params: [_.omit(model, "id")]
            })
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async delete(id) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/models`, {
                method: "delete",
                params: [id]
            })
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async put(model) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/models`, {
                method: "update",
                params: [model]
            })
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    }
}