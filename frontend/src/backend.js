import axios from "axios"
import config from "../config/backend"
import utils from "./utils"

export default {
    async export(option) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/export`, option))
    },
    async addModel(model) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/models.insert`, _.omit(model, "id")))
    },
    async delModel(name) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/models.delete`, {name}))
    },
    async updModel(model) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/models.update`, model))
    },
    async qryAllModels() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/models.selectAll`, {type: "model"}))
    },
    async qryAllStructs() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/models.selectAll`, {type: "struct"}))
    },
    async qryAllAvaTypes() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/models.selectAll`))
    },
    async qryAllBaseStructsName() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/structs.selectAllBases`))
    },
    async addLink(link) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/links.insert`, _.omit(link, "id")))
    },
    async qryAllLinks() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/links.selectAll`))
    },
    async qryAllApis() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/apis.selectAll`))
    },
    async qryApiByName(name) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/apis.selectByName`, {name}))
    },
    async addApi(api) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/apis.insert`, api))
    },
    async delApiByName(name) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/apis.deleteByName`, {name}))
    },
    async qryAllTempStep() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/temp.steps.selectAll`))
    },
    async addTempSteps(steps) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/temp.steps.insertMany`, {steps}))
    },
    async addTempStep(step) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/temp.steps.insert`, step))
    },
    async delTempStepByKey(key) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/temp.steps.deleteByKey`, {key}))
    },
    async delStep(delInfo) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/steps.delete`, delInfo))
    },
    async qryStepSymbols() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/specialSymbols`))
    },
    async addStep(step) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/steps.insert`, step))
    },
    async addDaoGroup(group) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/dao.groups.insert`, group))
    },
    async delDaoGroup(gpname) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/dao.groups.deleteByName`, {name: gpname}))
    },
    async addDaoInterface(istDaoItfcInf) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/dao.interface.insert`, istDaoItfcInf))
    },
    async delDaoInterface(delDaoItfcInf) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/dao.interface.delete`, delDaoItfcInf))
    },
    async addDaoConfig(implId, configs) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/dao.config.insert`, {"implement": implId, configs}))
    },
    async qryAllDaoGroups() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/dao.groups.selectAll`))
    },
    async qryAllImplMods() {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/mod.sign.selectAll`, {type: "dao_implement"}))
    },
    async qryModSignById(id) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/mod.info.selectBySignId`, {id}))
    },
    async updDaoGroupImpl(setDaoGpImplInf) {
        return await utils.reqBackend(axios.post(`${config.url}/backend-dashboard/backend/dao.groups.updateImplement`, setDaoGpImplInf))
    }
}