<template>
<dashboard>
    <div class="table-container">
        <el-button type="primary" @click="showAddDaoGroup = true">添加DAO组</el-button>
        <el-button type="primary" @click="showLoadDaoGroup = true">导入DAO组</el-button>
        <el-table class="mt-10" :data="daoGroups" style="width: 100%">
            <el-table-column type="expand">
                <template slot-scope="scope">
                    <el-table class="demo-table-expand" :data="scope.row.interfaces">
                        <el-table-column label="接口名" prop="name"/>
                        <el-table-column label="参数" prop="params">
                            <template slot-scope="subScope">
                                <div class="interval-container">
                                    <el-popover
                                        class="interval-item"
                                        v-for="param in subScope.row.params" :key="param.name"
                                        placement="top-start"
                                        trigger="hover"
                                        :content="param.type"
                                    >
                                        <el-tag size="small" slot="reference">{{param.name}}</el-tag>
                                    </el-popover>
                                </div>
                            </template>
                        </el-table-column>
                        <el-table-column label="返回值" prop="returns"/>
                        <el-table-column label="依赖模块" prop="requires"/>
                        <el-table-column label="描述" prop="desc"/>
                        <el-table-column label="配置" prop="setting">
                            <template slot-scope="subScope">
                                <el-button-group>
                                    <el-button size="mini" @click="showInterfaceDetail(scope.row, subScope.row)">详情</el-button>
                                    <el-button size="mini" type="danger" @click="delInterface(scope.row.name, subScope.row.name)">删除</el-button>
                                </el-button-group>
                            </template>
                        </el-table-column>
                    </el-table>
                </template>
            </el-table-column>
            <el-table-column label="组名" prop="name"/>
            <el-table-column label="类别" prop="category">
                <template slot-scope="scope">
                    <el-tag type="info">{{scope.row.category}}</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="语言" prop="language">
                <template slot-scope="scope">
                    <el-tag>{{scope.row.language}}</el-tag>
                </template>
            </el-table-column>
            <el-table-column label="实现" prop="implement">
                <template slot-scope="scope">
                    <el-button v-if="!scope.row.implement" size="mini" @click="implingGpName = scope.row.name">实例化DAO组</el-button>
                    <el-link v-else type="primary" size="mini" style="color:#409EFF" @click="showDaoImpl = {
                        grpName: scope.row.name, implId: scope.row.implement
                    }">
                        {{scope.row.implement}}
                    </el-link>
                </template>
            </el-table-column>
            <el-table-column label="配置" prop="setting" width="300">
                <template slot-scope="scope">
                    <el-button-group>
                        <el-button size="mini" @click="editingGroup = scope.row">添加接口</el-button>
                        <el-button v-show="scope.row.implement" size="mini" @click="confImpl(scope.row.implement)">实例配置</el-button>
                        <el-button size="mini" type="danger" @click="delGroup(scope.row.name)">删除</el-button>
                    </el-button-group>
                </template>
            </el-table-column>
        </el-table>
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="添加DAO组" :visible.sync="showAddDaoGroup" :modal-append-to-body="false" width="40vw">
            <edit-dao-group ref="add-dao-group-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="showAddDaoGroup = false">取 消</el-button>
                <el-button type="primary" @click="addGroup">确 定</el-button>
            </div>
        </el-dialog>
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="添加DAO接口" :visible="editingGroup !== null" :modal-append-to-body="false" width="40vw" @close="editingGroup = null">
            <edit-dao-interface ref="add-dao-interface-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="editingGroup = null">取 消</el-button>
                <el-button type="primary" @click="addInterface">确 定</el-button>
            </div>
        </el-dialog>
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="实例化DAO组" :visible="implingGpName.length != 0" :modal-append-to-body="false" width="50vw" @close="implingGpName = ''">
            <impl-dao-group ref="impl-dao-group-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="implingGpName = ''">取 消</el-button>
                <el-button type="primary" @click="implGroup">确 定</el-button>
            </div>
        </el-dialog>
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="实例化信息" :visible="showDaoImpl != null" :modal-append-to-body="false" width="30vw" @close="showDaoImpl = null">
            <dao-impl-info ref="dao-impl-info-form" :implId="showDaoImpl ? showDaoImpl.implId : ''"/>
            <div slot="footer" class="dialog-footer">
                <el-popover placement="top" width="160" v-model="showConfirmUistlImpl">
                    <p style="font-size: 0.5em">确定卸载DAO实例吗？</p>
                    <div style="text-align: right; margin: 0">
                        <el-button type="primary" size="mini" @click="uninstallImplement">确 定</el-button>
                    </div>
                    <el-button type="danger" slot="reference">卸 载</el-button>
                </el-popover>
            </div>
        </el-dialog>
        <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
        <el-dialog title="实例化DAO组" :visible.sync="showLoadDaoGroup" :modal-append-to-body="false" width="60vw">
            <load-dao-group ref="load-dao-group-form"/>
            <div slot="footer" class="dialog-footer">
                <el-button @click="showLoadDaoGroup = false">取 消</el-button>
                <el-button type="primary" @click="loadGroup">确 定</el-button>
            </div>
        </el-dialog>
    </div>
</dashboard>
</template>

<script>
import $ from "jquery"
import axios from "axios"
import backend from "../backend"
import dashboard from "../layouts/dashboard"
import editDaoGroup from "../forms/editDaoGroup"
import editDaoInterface from "../forms/editDaoInterface"
import implDaoGroup from "../forms/implDaoGroup"
import daoImplInfo from "../forms/daoImplInfo"
import chkLoadRes from "../forms/chkLoadRes"
import loadDaoGroup from "../forms/loadDaoGroup"

export default {
    components: {
        "dashboard": dashboard,
        "edit-dao-group": editDaoGroup,
        "edit-dao-interface": editDaoInterface,
        "impl-dao-group": implDaoGroup,
        "dao-impl-info": daoImplInfo,
        "load-dao-group": loadDaoGroup
    },
    data() {
        return {
            showAddDaoGroup: false,
            showLoadDaoGroup: false,
            showConfirmUistlImpl: false,
            editingGroup: null,
            implingGpName: "",
            showDaoImpl: null,
            daoGroups: [],
            categories: {
                databases: [{
                    name: "cctv"
                }]
            }
        }
    },
    async created() {
        await this.refresh()
    },
    methods: {
        async refresh() {
            let res = await backend.qryAllDaoGroups()
            if (typeof res === "string") {
                this.$message.error(`查询所有DAO组时发生错误：${res}`)
            } else {
                this.daoGroups = res.groups
            }
        },
        addGroup() {
            let addForm = this.$refs["add-dao-group-form"]
            addForm.$refs["edit-dao-group-form"].validate(async valid => {
                if (!valid) {
                    return false
                }
                let form = addForm.$refs["edit-dao-group-form"]
                let res = await backend.addDaoGroup(form.model)
                if (typeof res === "string") {
                    this.$message.error(`添加DAO组发生错误：${res}`)
                } else {
                    this.showAddDaoGroup = false
                    await this.refresh()
                }
            })
        },
        addInterface() {
            let addForm = this.$refs["add-dao-interface-form"]
            addForm.$refs["edit-dao-interface-form"].validate(async valid => {
                if (!valid) {
                    return false
                }
                let form = addForm.$refs["edit-dao-interface-form"]
                let res = await backend.addDaoInterface({
                    gpname: this.editingGroup.name,
                    interface: form.model
                })
                if (typeof res === "string") {
                    this.$message.error(`添加DAO接口发生错误：${res}`)
                } else {
                    this.editingGroup = null
                    await this.refresh()
                }
            })
        },
        delInterface(gpname, ifname) {
            this.$alert("确定删除接口？", "提示", {
                confirmButtonText: "确定",
                callback: async action => {
                    if (action !== "confirm") {
                        return
                    }
                    let res = await backend.delDaoInterface({gpname, ifname})
                    if (typeof res === "string") {
                        this.$message.error(`删除接口时发生错误：${res}`)
                    } else {
                        this.$message({
                            type: "info",
                            message: `接口（${ifname}）删除成功！`
                        })
                        await this.refresh()
                    }
                }
            })
        },
        delGroup(gpname) {
            this.$alert("确定删除组？", "提示", {
                confirmButtonText: "确定",
                callback: async action => {
                    if (action !== "confirm") {
                        return
                    }
                    let res = await backend.delDaoGroup(gpname)
                    if (typeof res === "string") {
                        this.$message.error(`删除DAO组时发生错误：${res}`)
                    } else {
                        this.$message({
                            type: "info",
                            message: `DAO组（${gpname}）删除成功！`
                        })
                        await this.refresh()
                    }
                }
            })
        },
        async implGroup() {
            const implDaoGroup = this.$refs["impl-dao-group-form"]
            const selImpl = implDaoGroup.$refs["impl-dao-group-form"].model
            let res = await backend.updDaoGroupImpl({
                gpname: this.implingGpName,
                implId: selImpl.id
            })
            if (typeof res === "string") {
                this.$message.error(`实例化DAO组时发生错误：${res}`)
                this.implingGpName = ""
                await this.refresh()
                return
            }
            const action = await this.$msgbox({
                title: "提示",
                message: this.$createElement(chkLoadRes, {
                    props: {
                        daoImpl: selImpl,
                    },
                    ref: "tip-load-temp-res-form"
                })
            })
            if (action !== "confirm") {
                this.implingGpName = ""
                await this.refresh()
                return
            }
            const result = this.$refs["tip-load-temp-res-form"]
            if (result.loadTmpStep) {
                const steps = (await axios.get(selImpl.tmpStepHref)).data.steps
                res = await backend.addTempSteps(steps)
                if (typeof res === "string") {
                    this.$message.error(`导入模板步骤时发生错误：${res}`)
                }
            }
            if (result.loadApiTemp) {
                const apis = (await axios.get(selImpl.apiTempHref)).data
                res = await backend.addTempApis(apis)
                if (typeof res === "string") {
                    this.$message.error(`导入模板接口时发生错误：${res}`)
                }
            }
            this.$message({
                type: "success",
                message: `DAO组（${this.implingGpName}）实例化成功`
            })
            this.implingGpName = ""
            await this.refresh()
        },
        showInterfaceDetail(group, itfc) {
            
        },
        async uninstallImplement() {
            let res = await backend.updDaoGroupImpl({
                gpname: this.showDaoImpl.grpName, implId: ""
            })
            if (typeof res === "string") {
                this.$message.error(`卸载DAO实例时发生错误：${res}`)
            } else {
                this.showConfirmUistlImpl = false
                this.showDaoImpl = null
                await this.refresh()
            }
        },
        async confImpl(implId) {
            let res = await backend.qryModSignById(implId)
            if (typeof res === "string") {
                this.$message.error(`查询模块标牌时发生错误：${res}`)
                return
            }
            const compTemps = (await axios.get(res.daoConfHref)).data
            let comps = []
            const h = this.$createElement;
            for (let compId in compTemps) {
                const cmpTmp = compTemps[compId]
                let content = null
                switch (cmpTmp.type) {
                    case "number":
                    case "text":
                    case "password":
                        content = [
                            h("input", {
                                "attrs": {
                                    "name": compId,
                                    "type": cmpTmp.type
                                },
                                "class": "form-control"
                            })
                        ]
                        break
                }
                comps.push(h("div", {"class": "form-group row"}, [
                    h("label", {"class": "col-sm-2 col-form-label"}, cmpTmp.label),
                    h("div", {"class": "col-sm-10"}, content)
                ]))
            }
            const action = await this.$msgbox({
                title: "接口信息",
                message: h("form", null, [comps]),
                showConfirmButton: true,
                customClass: "w-50"
            })
            if (action !== "confirm") {
                return
            }
            let configs = {}
            for (let compId in compTemps) {
                configs[compId] = $(`[name='${compId}']`)[0].value
            }
            res = await backend.addDaoConfig(implId, configs)
            if (typeof res === "string") {
                this.$message.error(`DAO实例配置时发生错误：${res}`)
            } else {
                this.$message({
                    type: "success",
                    message: `实例配置成功`
                })
                await this.refresh()
            }
        },
        async loadGroup() {
            const form = this.$refs["load-dao-group-form"].$refs["load-dao-group"]
            let res = await backend.addDaoGroup(form.model)
            if (typeof res === "string") {
                this.$message.error(`添加DAO组时发生错误：${res}`)
            } else {
                this.showLoadDaoGroup = false
                this.$message({
                    type: "success",
                    message: `模板组${form.model.name}导入成功`
                })
                await this.refresh()
            }
        }
    }
}
</script>

