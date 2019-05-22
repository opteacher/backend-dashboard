import config from "../../config/backend"
import axios from "axios"

export default {
    async post(model) {
        try {
            return await axios.post(`${config.url}/backend-dashboard/backend/modelsr`, model)
        } catch(e) {
            return e
        }
    }
}