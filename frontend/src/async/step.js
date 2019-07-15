import axios from "axios"

import config from "../../config/backend"
import utils from "../utils"

export default {
    async qryTmp() {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/steps.selectTemp`)
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    }
}