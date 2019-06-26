import _ from "lodash"
import axios from "axios"

import config from "../../config/backend"
import utils from "../utils"

export default {
    async add(model) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/models.insert`, _.omit(model, "id"))
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async del(id) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/models`, {
                method: "delete",
                params: ["`id`=?", [id]]
            })
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async upd(model) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/models`, {
                method: "update",
                params: ["`id`=?", [model.id], model]
            })
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qry() {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/models.selectAll`)
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    }
}