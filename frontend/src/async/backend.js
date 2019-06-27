import _ from "lodash"
import axios from "axios"

import config from "../../config/backend"
import utils from "../utils"

export default {
    async export(option) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/export`, option)
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    }
}