import axios from "axios"
import config from "../config/backend"
import utils from "./utils"

export default {
    async export(option) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/export`, option)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async addModel(model) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/models.insert`, _.omit(model, "id"))).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async delModel(name) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/models.delete`, {name})).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async updModel(model) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/models.update`, model)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qryAllModels() {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/models.selectAll`)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qryAllStructs() {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/structs.selectAll`)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async addLink(link) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/links.insert`, _.omit(link, "id"))).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qryAllLinks() {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/links.selectAll`)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qryAllApis() {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/apis.selectAll`)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qryApiByName(name) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/apis.selectByName`, {name})).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async addApi(api) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/apis.insert`, api)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async delApiByName(name) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/apis.deleteByName`, {name})).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qryStepTmp() {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/steps.selectTemp`)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async delStep(delInfo) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/steps.delete`, delInfo)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async qryStepSymbols() {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/specialSymbols`)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    },
    async addStep(flow) {
        try {
            return (await axios.post(`${config.url}/backend-dashboard/backend/steps.insert`, flow)).data.data
        } catch(e) {
            return utils.getErrorMsg(e)
        }
    }
}