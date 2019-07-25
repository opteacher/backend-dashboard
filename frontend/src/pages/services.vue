<template>
<dashboard>
    <info-bar @select-api="selectApi" @add-api="addApi"/>
    <div id="pnlFlows" class="w-100 h-100" style="position: absolute">
        <div style="position:absolute;width:0;height:0" v-for="flow in selApi.steps" :key="flow.index">
            <step-block :step="flow" @show-detail="showOperDetail"/>
        </div>
        <button v-for="btn in istStepBtns" :key="btn.nsuffix" :name="`istFlowBtn${btn.nsuffix}`" class="btn btn-success rounded-circle" type="button" style="position:absolute" @click="insertStep(btn.prev.index)">
            <i class="el-icon-plus"/>
        </button>
    </div>
    <svg id="pnlGraphs" class="w-100" style="position: absolute; z-index: -100; height: 100%">
        <flow-link v-for="btn in istStepBtns" :key="btn.nsuffix" :istFlowBtn="btn"/>        
    </svg>
    <el-button id="btnApiInfo" type="primary" size="small" icon="el-icon-info" v-if="selApi" @click="showApiInfo">
        {{selApi.name}}
    </el-button>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog :title="`步骤 #${selStep.index}`" :visible.sync="showStepDtlDlg" :modal-append-to-body="false" width="50vw">
        <step-detail ref="step-detail-form" :selStep="selStep"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="chgOperStep">编 辑</el-button>
            <el-button type="primary">确 定</el-button>
        </div>
    </el-dialog>
    <!-- +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++ -->
    <el-dialog :title="`添加步骤 #${flowInfo.flowIdx}`" :visible.sync="showAddStepDlg" :modal-append-to-body="false" width="50vw">
        <edit-flow ref="add-flow-form" :flowInfo="flowInfo"/>
        <div slot="footer" class="dialog-footer">
            <el-button @click="showAddStepDlg = false">取 消</el-button>
            <el-button type="primary" @click="addStep">确 定</el-button>
        </div>
    </el-dialog>
</dashboard>
</template>

<script>
import _ from "lodash"

import dashboard from "../layouts/dashboard"
import infoBar from "../components/infoBar"
import stepDetail from "../forms/stepDetail"
import backend from "../backend"
import apiInfo from "../forms/apiInfo"
import editFlow from "../forms/editFlow"
import stepBlock from "../components/stepBlock"
import flowLink from "../components/flowLink"

export default {
    components: {
        "dashboard": dashboard,
        "info-bar": infoBar,
        "step-detail": stepDetail,
        "edit-flow": editFlow,
        "step-block": stepBlock,
        "flow-link": flowLink
    },
    data() {
        return {
            selApi: {steps: []},
            selStep: {index: -1},
            index: 1,
            showStepDtlDlg: false,
            showAddStepDlg: false,
            fors: [],
            istStepBtns: [],
            flowInfo: {flowIdx: -1}
        }
    },
    methods: {
        async refresh() {
            let res = await backend.qryApiByName(this.selApi.name)
            if (typeof res === "string") {
                this.$message.error(`查询名为${this.selApi.name}的接口时发生错误：${res}`)
            } else {
                this.selectApi(res)
            }
        },
        selectApi(selApi) {
            this.selApi = selApi
            this.istStepBtns = []
            let locVars = this.selApi.params ? _.keys(this.selApi.params) : []
            this.selApi.steps = this.selApi.steps ? this.selApi.steps.map((flow, idx) => {
                // 做一些处理：只包含一个元素的输出数组全部设为空
                if (!flow.outputs || (flow.outputs.length === 1 && flow.outputs[0] === "")) {
                    flow.outputs = []
                }
                if (!flow.requires || (flow.requires.length === 1 && flow.requires[0] === "")) {
                    flow.requires = []
                }
                flow.index = idx

                if (idx === this.selApi.steps.length - 1) {
                    flow.isLast = true
                    this.istStepBtns.push({
                        apiName: selApi.name,
                        nsuffix: `_${idx}`,
                        prev: flow,
                        next: null,
                        locVars: locVars
                    })
                } else {
                    this.istStepBtns.push({
                        apiName: selApi.name,
                        nsuffix: `_${idx}_${idx + 1}`,
                        prev: flow,
                        next: this.selApi.steps[idx + 1],
                        locVars: locVars
                    })
                }

                if (idx !== 0) {
                    locVars.concat(flow.outputs)
                }
                return flow
            }) : []
            if (this.istStepBtns.length === 0) {
                // 没有按钮，说明流程中没有步骤，添加一个按钮用于初始化
                this.istStepBtns.push({
                    apiName: selApi.name,
                    nsuffix: "_0",
                    locVars: locVars
                })
            }
        },
        chgOperStep() {
            this.$refs["step-detail-form"].mode = (
                this.$refs["step-detail-form"].mode === "display" ? "editing" : "display"
            )
        },
        showApiInfo() {
            this.index++
            this.$msgbox({
                title: "接口信息",
                message: this.$createElement(apiInfo, {
                    props: {
                        api: this.selApi,
                    },
                    key: this.index
                }),
                showConfirmButton: false
            }).catch(err => {})
        },
        addApi(newApi) {
            console.log(newApi)
        },
        async addStep() {
            let form = this.$refs["add-flow-form"]
            let res = await backend.addStep({
                index: form.flowInfo.flowIdx,
                operStep: Object.assign(form.flow, {
                    apiName: form.flowInfo.apiName
                })
            })
            if (typeof res === "string") {
                this.$message.error(`插入流程发生错误：${res}`)
            } else {
                this.$message({
                    type: "success",
                    message: "插入成功！"
                })
                this.showAddStepDlg = false
                await this.refresh()
            }
        },
        showOperDetail(step) {
            this.selStep = step
            this.showStepDtlDlg = true
        },
        insertStep(prevIdx) {
            this.flowInfo.flowIdx = prevIdx + 1
            this.showAddStepDlg = true
        }
    }
}
</script>

<style lang="scss">
.api-params, .local-vars {
    font-size: 0.2rem;
    padding: .5vh .5vw;
}
.desc-panel:hover {
    cursor: pointer;
    background-color: #f8f9fa;
}
#btnApiInfo {
    position: fixed;
    right: 20px;
    bottom: 0;
    border-radius: 4px 4px 0 0;
}
</style>
