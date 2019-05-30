import axios from "axios"

import config from "../../config/backend"
import utils from "../utils"

export default {
    async add(relation) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/relations`, {
                method: "insert",
                params: [relation]
            })
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qry() {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/relations`, {
                method: "select",
                params: ["", []]
            })
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    }
}